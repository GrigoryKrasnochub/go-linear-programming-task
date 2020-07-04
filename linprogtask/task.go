package linprogtask

import (
	"errors"
	"fmt"
	"log"
)

type Task struct {
	VariablesCount  int
	ConditionsCount int
	CParams         []float64
	BParams         []float64
	AParams         [][]float64
}

func (t *Task) SetConditionsCount(conditionsCount int) {
	t.ConditionsCount = conditionsCount
	t.UpdateTask()
}

func (t *Task) ValidateConditionsCountValue(conditionsCount int) error {
	if conditionsCount < 1 {
		return errors.New("conditionCount(M) <= 0. It should be > 0")
	}

	if conditionsCount > t.VariablesCount {
		return errors.New("conditionCount(M) > VariablesCount(n). It should be <= VariablesCount(n)")
	}

	return nil
}

func (t *Task) SetVariablesCount(variablesCount int) {
	t.VariablesCount = variablesCount
	t.UpdateTask()
}

func (t *Task) ValidateVariablesCountValue(variablesCount int) error {
	if variablesCount < 1 {
		return errors.New("VariablesCount(n) <= 0. It should be > 0")
	}

	if variablesCount < t.ConditionsCount {
		return errors.New("conditionCount(M) > VariablesCount(n). It should be >= ConditionsCount(M)")
	}

	return nil
}

func (t *Task) SetC(numb int, c float64) {
	log.Printf("set CParam %d, value: %f \nCParams:\n%f", numb, c, t.CParams)
	t.CParams[numb] = c
}

func (t *Task) GetC(numb int) float64 {
	return t.CParams[numb]
}

func (t *Task) SetB(numb int, b float64) {
	log.Printf("set BParam %d, value: %f \nBParams:\n%f", numb, b, t.BParams)
	t.BParams[numb] = b
}

func (t *Task) SetBRow(b []float64) {
	log.Printf("set BParams value: %f \nBParams:\n%f", b, t.BParams)
	t.BParams = b
}

func (t *Task) GetB(numb int) float64 {
	return t.BParams[numb]
}

func (t *Task) ValidateCParam(param float64) error {
	if param < 0 {
		return errors.New(fmt.Sprintf("CParams should be bigger than 0"))
	}

	return nil
}

func (t *Task) SetA(i int, j int, a float64) {
	log.Printf("set AParam %d %d, value: %f \nAParams:\n%f", i, j, a, t.AParams)
	t.AParams[i][j] = a
}

func (t *Task) SetARow(i int, a []float64) {
	t.AParams[i] = a;
}

func (t *Task) GetA(i int, j int) float64 {
	return t.AParams[i][j]
}

/* recalculate size of CParams, AParams, BParams relative to current n, M*/
func (t *Task) UpdateTask() {
	if t.ConditionsCount > 0 {
		//update C slice
		t.CParams = t.updateParamSliceCap(t.CParams, t.ConditionsCount, "CParams")

		//update A slice
		if len(t.AParams) > 0 {
			if cap(t.AParams) >= t.ConditionsCount {
				t.AParams = t.AParams[0:t.ConditionsCount]
			} else {
				t.AParams = t.AParams[0:cap(t.AParams)]
				sliceSizeDif := t.ConditionsCount - cap(t.AParams)
				addSlice := make([][]float64, sliceSizeDif)
				t.AParams = append(t.AParams, addSlice...)
			}
		} else {
			t.AParams = make([][]float64, t.ConditionsCount)
		}
		log.Printf("AParams slice is changed cap:%d; len:%d; \n %f", cap(t.AParams), len(t.AParams), t.AParams)
	}

	if t.VariablesCount > 0 {
		//update B slice
		t.BParams = t.updateParamSliceCap(t.BParams, t.VariablesCount, "BParams")

		//update A slice
		if len(t.AParams) > 0 {
			for index, aParams := range t.AParams {
				t.AParams[index] = t.updateParamSliceCap(aParams, t.VariablesCount, fmt.Sprintf("AParams sub slice %d", index))
			}
		}
	}
}

func (t *Task) updateParamSliceCap(slice []float64, newcap int, name string) []float64 {
	if len(slice) > 0 {
		if cap(slice) >= newcap {
			slice = slice[0:newcap]
		} else {
			slice = slice[0:cap(slice)]
			sliceSizeDif := newcap - cap(slice)
			addSlice := make([]float64, sliceSizeDif)
			slice = append(slice, addSlice...)
		}
	} else {
		slice = make([]float64, newcap)
	}
	log.Printf("%s slice is changed cap:%d; len:%d; \n %f", name, cap(slice), len(slice), slice)

	return slice
}

/*Check all values to be correct*/
func (t *Task) IsSystemReadyToCalc() error {
	conditionsCount := t.ConditionsCount
	variablesCount := t.VariablesCount

	conditionCountError := t.ValidateConditionsCountValue(conditionsCount)
	if conditionCountError != nil {
		return conditionCountError
	}

	variablesCountError := t.ValidateVariablesCountValue(variablesCount)
	if variablesCountError != nil {
		return variablesCountError
	}

	var CParamError error
	for index, cParam := range t.CParams {
		CParamError = t.ValidateCParam(cParam)
		if CParamError != nil {
			CParamError = errors.New(fmt.Sprintf("%s in CParam %d", CParamError.Error(), index))
			break
		}
	}
	if CParamError != nil {
		return CParamError
	}

	return nil
}

func (t *Task) DoCalc() Result {
	calc := calc{task: *t}
	result := calc.doCalc()
	return result
}
