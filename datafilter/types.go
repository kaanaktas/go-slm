package datafilter

type Validate interface {
	Validate(data *string) bool
	ToString() string
	Disable() bool
}

//filter types
const (
	PAN   = "pan"
	OWASP = "owasp"
)

type pattern struct {
	Name       string `json:"name"`
	Rule       string `json:"rule"`
	Sample     string `json:"sample"`
	Message    string `json:"message"`
	IsDisabled bool   `json:"disable"`
}
