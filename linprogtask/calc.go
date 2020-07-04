package linprogtask

import (
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
)

type calc struct {
	task                  Task
	negativeBParamIndex   int
	bParamsAlter          []float64
	ficLimParams          []float64
	extendedAParams       [][]float64 // from task with xi <= 0 limitations
	extendedCParams       []float64   // from task with xi <= 0 limitations
	activeLimitationIndex *int
	minCoefAParamsCParams float64
	result                Result
}

type Result struct {
	CalculationStatus      string
	SolutionCoordinates    []float64
	TargetFunctionResult   float64
	IterationCalc          int
	CalculationWasFinished bool
}

func (c *calc) doCalc() Result {
	log.Printf("Calculation started with params: \n%+v\n", c.task)
	prepareDataForEvolutionErr := c.prepareDataForEvolution()
	if prepareDataForEvolutionErr != nil {
		c.result.CalculationStatus = prepareDataForEvolutionErr.Error()
		return c.result
	}

	log.Printf("Start evolution: \n%+v\n", c)
	var doEvolutionErr error
	doEvolutionErr = nil
	for doEvolutionErr == nil {
		doEvolutionErr = c.doEvolution()
		if doEvolutionErr != nil {
			c.result.CalculationStatus = doEvolutionErr.Error()
			return c.result
		}
	}

	return c.result
}

func (c *calc) prepareDataForEvolution() error {
	//add to AParams and CParams limitation xi <= 0
	log.Println("expand AParams and CParams with xi <= 0 limitations")
	c.expandACParams()

	c.result.SolutionCoordinates = make([]float64, c.task.VariablesCount)
	c.result.CalculationWasFinished = false

	log.Println("Search for negative BParam index (i*)")
	searchNegativeBParamIndexErr := c.searchNegativeBParamIndex()
	if searchNegativeBParamIndexErr != nil {
		return searchNegativeBParamIndexErr
	}

	log.Println("Generate alter BParams (b~)")
	c.generateAlterBParams()

	log.Println("Generate fictitious limitation params (hi)")
	c.calculateFictitiousLimitationParams()
	log.Println("Replace fictional limitation. Calculate limitation params using extendedAParams.")
	c.replaceFictionalLimitationInExtendedAParams()
	return nil
}

func (c *calc) expandACParams() {

	c.extendedAParams = make([][]float64, c.task.VariablesCount)
	c.extendedCParams = make([]float64, c.task.VariablesCount)
	for index, _ := range c.extendedAParams {
		c.extendedAParams[index] = make([]float64, c.task.VariablesCount)
		c.extendedAParams[index][index] = 1
	}

	for _, extendedAParam := range c.task.AParams {
		tmpExtendedAParam := make([]float64, c.task.ConditionsCount)
		_ = copy(tmpExtendedAParam, extendedAParam)
		c.extendedAParams = append(c.extendedAParams, tmpExtendedAParam)
	}

	extendedCParams := make([]float64, c.task.VariablesCount)
	_ = copy(extendedCParams, c.task.CParams)
	c.extendedCParams = append(c.extendedCParams, extendedCParams...)

	log.Printf("extended AParams: \n%+v\n extended CParams \n%+v\n", c.extendedAParams, c.extendedCParams)
}

func (c *calc) searchNegativeBParamIndex() error {
	bParamCalculationPossible := false
	for index, bParam := range c.task.BParams {
		if bParam < 0 {
			bParamCalculationPossible = true
			c.negativeBParamIndex = index
			break
		}
	}
	log.Printf("Negative BParam index %d\n", c.negativeBParamIndex)

	if !bParamCalculationPossible {
		return errors.New("Calculation couldn't be started, all of BParams >= 0 ")
	}
	return nil
}

func (c *calc) generateAlterBParams() {
	c.bParamsAlter = make([]float64, c.task.VariablesCount)
	for index, _ := range c.bParamsAlter {
		c.bParamsAlter[index] = rand.Float64() + 1
	}
	log.Printf("Alter BParams generated (b~) \n%f\n", c.bParamsAlter)
}

func (c *calc) calculateFictitiousLimitationParams() {
	c.ficLimParams = make([]float64, c.task.VariablesCount)
	for index, bParam := range c.task.BParams {
		if index == c.negativeBParamIndex {
			c.ficLimParams[index] = bParam / c.bParamsAlter[index]
		} else {
			c.ficLimParams[index] = (bParam - c.bParamsAlter[index]) / c.bParamsAlter[c.negativeBParamIndex]
		}
	}
	log.Printf("Additional params (h)\n%f\n", c.ficLimParams)

}

func (c *calc) replaceFictionalLimitationInExtendedAParams() {
	extendedAParamsCopy := make([][]float64, len(c.extendedAParams))
	for index, _ := range extendedAParamsCopy {
		extendedAParamsCopy[index] = make([]float64, len(c.extendedAParams[index]))
		_ = copy(extendedAParamsCopy[index], c.extendedAParams[index])
	}
	for index, _ := range c.extendedAParams {
		if c.extendedAParams[index][c.negativeBParamIndex] != 0 {
			for innerIndex, _ := range c.extendedAParams[index] {
				if innerIndex == c.negativeBParamIndex {
					c.extendedAParams[index][innerIndex] = extendedAParamsCopy[index][innerIndex] / c.ficLimParams[c.negativeBParamIndex]
				} else {
					c.extendedAParams[index][innerIndex] = extendedAParamsCopy[index][innerIndex] - (extendedAParamsCopy[index][c.negativeBParamIndex] * c.ficLimParams[innerIndex] / c.ficLimParams[c.negativeBParamIndex])
				}
			}
		}
	}

	log.Printf("Limitations params after add fictitional condition \n%+v\n", c.extendedAParams)
}

func (c *calc) doEvolution() error {
	c.result.IterationCalc++
	log.Printf("Evolution iteration %d started", c.result.IterationCalc)

	log.Println("Search for active limitation")
	searchActiveLimitationErr := c.searchActiveLimitationIndex()
	if searchActiveLimitationErr != nil {
		return searchActiveLimitationErr
	}

	c.result.SolutionCoordinates[c.negativeBParamIndex] += c.minCoefAParamsCParams
	c.result.TargetFunctionResult += c.minCoefAParamsCParams * c.bParamsAlter[c.negativeBParamIndex]

	log.Println("Move point of start coordinate")
	c.moveCoordinateStartPoint()

	log.Println("Searching and replacing variable")
	searchAndRepalceVariableErr := c.searchAndReplaceVariable()
	if searchAndRepalceVariableErr != nil {
		return searchAndRepalceVariableErr
	}

	return nil
}

func (c *calc) searchActiveLimitationIndex() error {
	c.activeLimitationIndex = nil

	AParamsBParamNegativeIndex := make([]float64, len(c.extendedCParams))
	coefsCParamsAParamsBNegative := make([]float64, len(c.extendedCParams))

	for index, _ := range coefsCParamsAParamsBNegative {
		AParamsBParamNegativeIndex[index] = c.extendedAParams[index][c.negativeBParamIndex]
		coefsCParamsAParamsBNegative[index] = c.extendedCParams[index] / c.extendedAParams[index][c.negativeBParamIndex]
	}

	c.minCoefAParamsCParams = math.MaxFloat64
	for index, coefCParamsAParamsBNegative := range coefsCParamsAParamsBNegative {
		if AParamsBParamNegativeIndex[index] > 0 {
			if c.minCoefAParamsCParams > coefCParamsAParamsBNegative {
				newIndex := index
				c.activeLimitationIndex = &newIndex
				c.minCoefAParamsCParams = coefCParamsAParamsBNegative
			}
		}
	}

	if c.activeLimitationIndex == nil {
		return errors.New("solution is infinitely positive number")
	}
	return nil
}

func (c *calc) moveCoordinateStartPoint() {
	for index, extendedAParam := range c.extendedAParams {
		if extendedAParam[c.negativeBParamIndex] != 0 {
			c.extendedCParams[index] = c.extendedCParams[index] - c.minCoefAParamsCParams*extendedAParam[c.negativeBParamIndex]
		}
	}
}

func (c *calc) searchAndReplaceVariable() error {
	coefsAlterBParamsAparms := make([]float64, c.task.VariablesCount)
	AParamsOfActiveLim := make([]float64, c.task.VariablesCount)

	for index, _ := range coefsAlterBParamsAparms {
		AParamsOfActiveLim[index] = c.extendedAParams[*c.activeLimitationIndex][index]
		coefsAlterBParamsAparms[index] = c.bParamsAlter[index] / c.extendedAParams[*c.activeLimitationIndex][index]
	}

	minAlterBParamsAparms := math.MaxFloat64
	var replacedVariableIndex *int
	replacedVariableIndex = nil
	for index, coefAlterBParamsAparms := range coefsAlterBParamsAparms {
		if AParamsOfActiveLim[index] > 0 && coefAlterBParamsAparms < minAlterBParamsAparms {
			newIndex := index
			minAlterBParamsAparms = coefAlterBParamsAparms
			replacedVariableIndex = &newIndex
		}
	}

	if replacedVariableIndex == nil {
		return errors.New("cant find variable for replacing")
	}

	//replacing
	bParamsAlterCopy := make([]float64, len(c.bParamsAlter))
	_ = copy(bParamsAlterCopy, c.bParamsAlter)

	for index, _ := range bParamsAlterCopy {
		switch index {
		case *replacedVariableIndex:
			bParamsAlterCopy[index] = c.bParamsAlter[index] / c.extendedAParams[*c.activeLimitationIndex][*replacedVariableIndex]
		case c.negativeBParamIndex:
			bParamsAlterCopy[index] = c.bParamsAlter[index] - c.bParamsAlter[*replacedVariableIndex]*c.extendedAParams[*c.activeLimitationIndex][c.negativeBParamIndex]/c.extendedAParams[*c.activeLimitationIndex][*replacedVariableIndex]
		default:
			bParamsAlterCopy[index] = c.bParamsAlter[index] - c.bParamsAlter[*replacedVariableIndex]*c.extendedAParams[*c.activeLimitationIndex][index]/c.extendedAParams[*c.activeLimitationIndex][*replacedVariableIndex]
		}
	}

	_ = copy(c.bParamsAlter, bParamsAlterCopy)

	for index, extendedAParam := range c.extendedAParams {
		if c.extendedAParams[index][*replacedVariableIndex] != 0 && *c.activeLimitationIndex != index {
			for innerIndex, _ := range extendedAParam {
				if innerIndex != *replacedVariableIndex {
					c.extendedAParams[index][innerIndex] = c.extendedAParams[index][innerIndex] - c.extendedAParams[index][*replacedVariableIndex]*c.extendedAParams[*c.activeLimitationIndex][innerIndex]/c.extendedAParams[*c.activeLimitationIndex][*replacedVariableIndex]
				}
			}
			for innerIndex, _ := range extendedAParam {
				if innerIndex == *replacedVariableIndex {
					c.extendedAParams[index][innerIndex] = c.extendedAParams[index][innerIndex] / c.extendedAParams[*c.activeLimitationIndex][*replacedVariableIndex]
				}
			}
		}
	}

	for index, extendedAParam := range c.extendedAParams {
		if c.extendedAParams[index][*replacedVariableIndex] != 0 && *c.activeLimitationIndex == index {
			c.extendedAParams[index] = make([]float64, len(extendedAParam))
			c.extendedAParams[index][*replacedVariableIndex] = 1

			c.extendedCParams[index] = 0
		}
	}

	if *replacedVariableIndex == c.negativeBParamIndex {
		c.result.CalculationWasFinished = true
		return errors.New(fmt.Sprintf("Calculation complete!\nSolution was found\nTarget function value = %f\nItteration counter = %d", c.result.TargetFunctionResult, c.result.IterationCalc))
	}
	return nil
}
