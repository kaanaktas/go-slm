package datafilter

type validate interface {
	validate(data *string) bool
	toString() string
	disable() bool
}

//filter types
const (
	PAN   = "pan"
	OWASP = "owasp"
)

type pattern struct {
	Name    string `json:"name"`
	Rule    string `json:"rule"`
	Sample  string `json:"sample"`
	Message string `json:"message"`
	Disable bool   `json:"disable"`
}
