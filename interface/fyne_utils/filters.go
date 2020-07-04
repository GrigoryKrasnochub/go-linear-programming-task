package fyne_utils

import (
	"regexp"
)

var regexpFloatNumber = regexp.MustCompile(`\d+\.*\d*`)
var regexpNegativeFloatNumber = regexp.MustCompile(`(?:(?:-\d*)|(?:\d+))\.?\d*`)
var regexpIntegerNumber = regexp.MustCompile(`\d+`)

func isNumeric(input string) {
	regexpFloatNumber.FindString(input)
}

func FilterIntegerNumber(input *string) {
	*input = regexpIntegerNumber.FindString(*input)
}

func FilterPositiveFloatNumber(input *string) {
	*input = regexpFloatNumber.FindString(*input)
}

func FilterFloatNumber(input *string) {
	*input = regexpNegativeFloatNumber.FindString(*input)
}
