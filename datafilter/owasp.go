package datafilter

import "regexp"

type owasp struct {
	pattern
}

func (o owasp) Validate(data *string) bool {
	matched, _ := regexp.MatchString(o.Rule, *data)
	return matched
}

func (o owasp) ToString() string {
	return o.Name + " " + o.Message
}

func (o owasp) Disable() bool {
	return o.IsDisabled
}
