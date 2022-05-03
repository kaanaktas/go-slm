package datafilter

import "regexp"

type pan struct {
	pattern
}

func (p pan) validate(data *string) bool {
	r := regexp.MustCompile(p.Rule)
	matchList := r.FindAllString(*data, -1)
	for _, v := range matchList {
		if isValidPanNumber(v) {
			return true
		}
	}

	return false
}

func (p pan) toString() string {
	return p.Name + " " + p.Message
}

func (p pan) disable() bool {
	return p.Disable
}
