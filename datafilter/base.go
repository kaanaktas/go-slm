package datafilter

import "regexp"

type Validate interface {
	Validate(data *string) bool
	ToString() string
	Disable() bool
}

type pattern struct {
	Name       string `json:"name"`
	Rule       string `json:"rule"`
	Sample     string `json:"sample"`
	Message    string `json:"message"`
	IsDisabled bool   `json:"disable"`
}

func (p *pattern) Validate(data *string) bool {
	matched, _ := regexp.MatchString(p.Rule, *data)
	return matched
}

func (p *pattern) ToString() string {
	return p.Name + " " + p.Message
}

func (p *pattern) Disable() bool {
	return p.IsDisabled
}
