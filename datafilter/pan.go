package datafilter

import (
	"math"
	"strconv"
	"strings"
)

type pan struct {
	patternValidator
}

func (p *pan) Validate(data *string) bool {
	dataWithoutSpace := strings.ReplaceAll(*data, " ", "")
	matchList := p.Regex.FindAllString(dataWithoutSpace, -1)
	for _, v := range matchList {
		if isValidPan(v) {
			return true
		}
	}

	return false
}

func isValidPan(number string) bool {
	clearText := strings.ReplaceAll(number, " ", "")
	if _, err := strconv.ParseInt(clearText, 10, 64); err != nil {
		return false
	}
	digits := strings.Split(clearText, "")
	lengthOfString := len(digits)

	if lengthOfString < 2 {
		return false
	}

	sum := 0
	doubleChecker := false

	for i := lengthOfString - 1; i > -1; i-- {
		digit, _ := strconv.Atoi(digits[i])
		if doubleChecker {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		doubleChecker = !doubleChecker
	}

	return math.Mod(float64(sum), 10) == 0
}
