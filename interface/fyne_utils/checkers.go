package fyne_utils

import (
	"errors"
	"strconv"
)

func IsPositiveIntNumber(val string) error {
	result, err := strconv.Atoi(val)
	if err != nil {
		return errors.New("value should be numeric")
	}
	if result < 1 {
		return errors.New("value should be positive")
	}
	return nil
}
