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
	Name    string         `yaml:"name"`
	Regex   *regexp.Regexp `yaml:"-"`
	Rule    string         `yaml:"rule"`
	Sample  string         `yaml:"sample"`
	Message string         `yaml:"message"`
	Disable bool           `yaml:"disable"`
}

// Validate checks whether the given data string satisfies the regular expression pattern.
func (pv *patternValidator) Validate(data *string) bool {
	return pv.Regex.MatchString(*data)
}

// ToString returns a string representation of the validation rule.
func (pv *patternValidator) ToString() string {
	return pv.Name + " " + pv.Message
}

// IsDisabled indicates whether the validation rule is disabled.
func (pv *patternValidator) IsDisabled() bool {
	return pv.Disable
}

func (pv *patternValidator) compileRule() {
	pv.Regex = regexp.MustCompile(pv.Rule)
}
