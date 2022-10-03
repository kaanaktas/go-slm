package datafilter

import "regexp"

// Validator interface represents a validation rule.
type Validator interface {
	// Validate checks whether the given data string satisfies the validation rule.
	Validate(data *string) bool
	// ToString returns a string representation of the validation rule.
	ToString() string
	// IsDisabled indicates whether the validation rule is disabled.
	IsDisabled() bool
}

// PatternValidator represents a validation rule based on a regular expression patternValidator.
type patternValidator struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Sample  string `json:"sample"`
	Message string `json:"message"`
	Disable bool   `json:"disable"`
}

// Validate checks whether the given data string satisfies the regular expression pattern.
func (pv *patternValidator) Validate(data *string) bool {
	matched, _ := regexp.MatchString(pv.Rule, *data)
	return matched
}

// ToString returns a string representation of the validation rule.
func (pv *patternValidator) ToString() string {
	return pv.Name + " " + pv.Message
}

// IsDisabled indicates whether the validation rule is disabled.
func (pv *patternValidator) IsDisabled() bool {
	return pv.Disable
}
