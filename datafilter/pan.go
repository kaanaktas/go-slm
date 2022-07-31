package datafilter

import (
	"regexp"
	"strings"
)

type pan struct {
	pattern
}

func (p pan) Validate(data *string) bool {
	dataWithoutSpace := strings.ReplaceAll(*data, " ", "")
	r := regexp.MustCompile(p.Rule)
	matchList := r.FindAllString(dataWithoutSpace, -1)
	for _, v := range matchList {
		if isValidPanNumber(v) {
			return true
		}
	}

	return false
}

func (p pan) ToString() string {
	return p.Name + " " + p.Message
}

func (p pan) Disable() bool {
	return p.IsDisabled
}
