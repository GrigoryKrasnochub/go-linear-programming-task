package _interface

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"github.com/GrigoryKrasnochub/go-linear-programming-task/interface/fyne_utils"
	"github.com/GrigoryKrasnochub/go-linear-programming-task/linprogtask"
)

type ProgramInterface struct {
	application fyne.App
	window      fyne.Window
	widgets     ProgramInterfaceWidgets
	linTask     linprogtask.Task
}

type ProgramInterfaceWidgets struct {
	setConditionsNumbersWidget   *widget.Form
	setConditionsVariablesWidget *widget.Box
}

func InitInterface() ProgramInterface {
	a := app.New()
	w := a.NewWindow("Linear-programming")
	a.Settings().SetTheme(theme.LightTheme())

	w.CenterOnScreen()

	pinter := ProgramInterface{}
	linTask := linprogtask.Task{}

	conditionsCount := widget.NewEntry()
	conditionsCount.SetPlaceHolder("enter number > 0")

	variablesCount := widget.NewEntry()
	variablesCount.SetPlaceHolder("enter number > 0")

	conditionsCount.OnChanged = func(val string) {
		if val != "" {
			err := fyne_utils.IsPositiveIntNumber(val)
			if err == nil && variablesCount.Text != "" {
				variablesCountVal, _ := strconv.Atoi(variablesCount.Text)
				conditionsCountVal, _ := strconv.Atoi(val)
				if conditionsCountVal > variablesCountVal {
					err = errors.New("number should be smaller or equal (n) value")
				}
			}

			if err != nil {
				pinter.showError(err)
				conditionsCount.SetText("1")
			} else {
				if variablesCount.Text != "" && val != "" {
					variablesCountVal, _ := strconv.Atoi(variablesCount.Text)
					conditionsCountVal, _ := strconv.Atoi(val)
					pinter.linTask.SetVariablesCount(variablesCountVal)
					pinter.linTask.SetConditionsCount(conditionsCountVal)
					pinter.buildTable(variablesCountVal, conditionsCountVal)
				}
			}
		}
	}

	variablesCount.OnChanged = func(val string) {
		if val != "" {
			err := fyne_utils.IsPositiveIntNumber(val)
			variablesCountVal, _ := strconv.Atoi(variablesCount.Text)
			conditionsCountVal, _ := strconv.Atoi(val)
			if conditionsCountVal > variablesCountVal {
				err = errors.New("number should be bigger or equal (M) value")
			}

			if err != nil {
				pinter.showError(err)
				variablesCount.SetText("1")
			} else {
				if conditionsCount.Text != "" && val != "" {
					variablesCountVal, _ := strconv.Atoi(variablesCount.Text)
					conditionsCountVal, _ := strconv.Atoi(val)
					pinter.linTask.SetVariablesCount(variablesCountVal)
					pinter.linTask.SetConditionsCount(conditionsCountVal)
					pinter.buildTable(variablesCountVal, conditionsCountVal)
				}
			}
		}
	}

	setConditionNumbersWidget := &widget.Form{
		Items: []*widget.FormItem{
			{"Conditions Count (M)", conditionsCount},
			{"Variables Count (n)", variablesCount},
		},
	}

	setConditionsVariablesWidget := widget.NewHBox()

	pinter = ProgramInterface{
		application: a,
		window:      w,
		widgets: ProgramInterfaceWidgets{
			setConditionsNumbersWidget:   setConditionNumbersWidget,
			setConditionsVariablesWidget: setConditionsVariablesWidget,
		},
		linTask: linTask,
	}

	pinter.redrawWindow()

	return pinter
}

func (pinter *ProgramInterface) redrawWindow() {
	pinter.window.SetContent(widget.NewVBox(
		pinter.widgets.setConditionsNumbersWidget,
		pinter.widgets.setConditionsVariablesWidget,
	))
}

func (pinter *ProgramInterface) ShowInterface() {
	pinter.window.ShowAndRun()
}

func (pinter *ProgramInterface) showError(err error) {
	dialog.ShowError(err, pinter.window)
}

func (pinter *ProgramInterface) showMessage(header string, message string) {
	dialog.ShowInformation(header, message, pinter.window)
}

/*redraw table*/
func (pinter *ProgramInterface) buildTable(variablesCount int, conditionsCount int) {
	log.Printf("Start table redrawing\n%+v\n", pinter.linTask)

	newTable := widget.NewHBox()

	for tableColumnsCount := -4; tableColumnsCount < variablesCount; tableColumnsCount++ {
		newColumn := widget.NewVBox()

		//generate info columns
		if tableColumnsCount < 0 {
			emptyLabel := widget.NewLabel("")
			newColumn.Append(emptyLabel)

			newItem := widget.NewEntry()
			newItem.Disable()
			if tableColumnsCount == -1 {
				newItem.Text = "B"
			}
			newColumn.Append(newItem)

			for index := 0; index < conditionsCount; index++ {
				switch tableColumnsCount {
				case -4:
					//Row number
					newItem := widget.NewEntry()
					newItem.SetText(fmt.Sprint(index + 1))
					newItem.Disable()
					newColumn.Append(newItem)
				case -3:
					//CParam title
					newItem := widget.NewEntry()
					newItem.SetText(fmt.Sprint("C"))
					newItem.Disable()
					newColumn.Append(newItem)
				case -2:
					//CParam Entry
					cParamEntryIndex := index
					cParamEntry := widget.NewEntry()
					cParamEntry.Text = fmt.Sprint(pinter.linTask.GetC(cParamEntryIndex))
					cParamEntry.OnChanged = func(val string) {
						fyne_utils.FilterPositiveFloatNumber(&val)
						cParamEntry.Text = val

						if val != "" {
							parsedVal, _ := strconv.ParseFloat(val, 64)
							err := pinter.linTask.ValidateCParam(parsedVal)
							if err != nil {
								pinter.showError(err)
							} else {
								pinter.linTask.SetC(cParamEntryIndex, parsedVal)
							}
						}
					}
					newColumn.Append(cParamEntry)
				case -1:
					//AParam tittle
					newItem := widget.NewEntry()
					newItem.SetText("A")
					newItem.Disable()
					newColumn.Append(newItem)
				}
			}
		}

		//generate columns
		if tableColumnsCount >= 0 {
			colNumbLabel := widget.NewLabelWithStyle(fmt.Sprint(tableColumnsCount+1), fyne.TextAlignCenter, fyne.TextStyle{
				Bold:      false,
				Italic:    false,
				Monospace: false,
			})
			newColumn.Append(colNumbLabel)

			bParamEntryIndex := tableColumnsCount
			bParamEntry := widget.NewEntry()
			bParamEntry.Text = fmt.Sprint(pinter.linTask.GetB(bParamEntryIndex))
			bParamEntry.OnChanged = func(val string) {
				fyne_utils.FilterFloatNumber(&val)
				bParamEntry.SetText(val)
				if val != "" {
					parsedVal, _ := strconv.ParseFloat(val, 64)
					pinter.linTask.SetB(bParamEntryIndex, parsedVal)
				}
				bParamEntry.Refresh()
			}
			newColumn.Append(bParamEntry)

			for index := 0; index < conditionsCount; index++ {
				cParamEntryIndexCol := index
				cParamEntryIndexRow := tableColumnsCount
				aParamEntry := widget.NewEntry()
				aParamEntry.Text = fmt.Sprint(pinter.linTask.GetA(cParamEntryIndexCol, cParamEntryIndexRow))
				aParamEntry.OnChanged = func(val string) {
					fyne_utils.FilterFloatNumber(&val)
					aParamEntry.SetText(val)
					if val != "" {
						parsedVal, _ := strconv.ParseFloat(val, 64)
						pinter.linTask.SetA(cParamEntryIndexCol, cParamEntryIndexRow, parsedVal)
					}
				}
				newColumn.Append(aParamEntry)
			}
		}

		newTable.Append(newColumn)
	}

	startCalcButton := widget.NewButton("Start calculation", func() {
		err := pinter.linTask.IsSystemReadyToCalc()
		if err != nil {
			pinter.showError(err)
		}

		calcResult := pinter.linTask.DoCalc()
		pinter.showMessage("Calculation status", fmt.Sprintf("\n%+v\n", calcResult))
	})

	pinter.widgets.setConditionsVariablesWidget = widget.NewVBox(newTable, startCalcButton)

	pinter.redrawWindow()
	pinter.widgets.setConditionsVariablesWidget.Refresh()
}
