package datafilter

import "regexp"

type pan struct {
	pattern
}

func (p pan) Validate(data *string) bool {
	r := regexp.MustCompile(p.Rule)
	matchList := r.FindAllString(*data, -1)
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
