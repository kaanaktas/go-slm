package datafilter

import "regexp"

type owasp struct {
	pattern
}

func (o owasp) validate(data *string) bool {
	matched, _ := regexp.MatchString(o.Rule, *data)
	return matched
}

func (o owasp) toString() string {
	return o.Name + " " + o.Message
}

func (o owasp) disable() bool {
	return o.Disable
}
