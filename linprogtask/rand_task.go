package linprogtask

import (
	"fmt"
	"math/rand"
	"os"
)

func CalcRandom(resultsCount int, variableCount int, conditionCount int, destination string, showProgress bool) {
	results := make([]Result, resultsCount)
	tasks := make([]Task, resultsCount)
	progressPercent := 0
	onePercent := float64(resultsCount) / float64(100)
	for solvedTaskResultsCount := 0; solvedTaskResultsCount < resultsCount; {
		task := RandomLinProgTask(variableCount, conditionCount)
		tasks[solvedTaskResultsCount] = task
		result := task.DoCalc()

		if !result.CalculationWasFinished {
			continue
		}

		results[solvedTaskResultsCount] = result
		solvedTaskResultsCount += 1

		//Percent calculation
		if !showProgress {
			continue
		}

		tProgressPercent := int(float64(solvedTaskResultsCount) / onePercent)
		if progressPercent != tProgressPercent {
			progressPercent = tProgressPercent
			fmt.Printf("Progress %d%% \n", progressPercent)
		}

	}

	headerString := fmt.Sprintf("SolutionNumber\tIterationsCount\tTargetFunctionValue\tSolutionCoordinates\tVariablesCount\tAParams\tCParams\tBParams\n")

	if destination == "" {
		fmt.Println("Calculation was finished")
		fmt.Println(headerString)
		for index, result := range results {
			fmt.Printf("%d\t%d\t%f\t%+v\t%d\t%+v\t%+v\t%+v\n", index, result.IterationCalc, result.TargetFunctionResult, result.SolutionCoordinates, variableCount, tasks[index].AParams, tasks[index].CParams, tasks[index].BParams)
		}
	} else {

		fileContent := make([]string, resultsCount+1)

		fileContent[0] = headerString
		for i := 0; i < resultsCount; i++ {
			result := results[i]
			fileContent[i+1] = fmt.Sprintf("%d\t%d\t%f\t%+v\t%d\t%+v\t%+v\t%+v\n", i, result.IterationCalc, result.TargetFunctionResult, result.SolutionCoordinates, variableCount, tasks[i].AParams, tasks[i].CParams, tasks[i].BParams)
		}

		err := saveToFile(destination, fileContent)
		if err != nil {
			fmt.Println("Oops something went wrong\nError:", err)
		}
	}

}

func RandomLinProgTask(variableCount int, conditionCount int) Task {
	var task Task
	task.SetConditionsCount(conditionCount)
	task.SetVariablesCount(variableCount)
	task.randomCValue()
	task.randomBValue()
	task.randomAValue()
	return task
}

func (t *Task) randomCValue() {
	index := 0
	generatedCParam := randFloats(0, 20000, t.ConditionsCount)
	for index < t.ConditionsCount {
		validateCParamErr := t.ValidateCParam(generatedCParam[index])
		if validateCParamErr == nil {
			t.SetC(index, generatedCParam[index])
			index++
		} else {
			generatedCParam[index] = randFloats(0, 20000, 1)[0]
		}
	}
}

func (t *Task) randomBValue() {
	generatedBParam := randFloats(-20000, 20000, t.VariablesCount)
	t.SetBRow(generatedBParam)
}

func (t *Task) randomAValue() {
	index := 0
	for index < t.ConditionsCount {
		generatedAParam := randFloats(-20000, 20000, t.VariablesCount)
		t.SetARow(index, generatedAParam)
		index++
	}
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}

func saveToFile(destination string, rows []string) error {
	f, err := os.Create(destination)
	if err != nil {
		return err
	}

	defer func() {
		closeErr := f.Close()
		if err == nil {
			err = closeErr
		}
	}()

	for _, row := range rows {
		_, err = f.WriteString(row)
		if err != nil {
			return err
		}
		err = f.Sync()
		if err != nil {
			return err
		}
	}

	return nil
}
