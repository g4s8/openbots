package spec

import (
	"fmt"
)

// Check is a type for validator checks.
type Check string

const (
	CheckNotEmpty Check = "not_empty"
	CheckIsInt    Check = "is_int"
	CheckIsFloat  Check = "is_float"
	CheckIsBool   Check = "is_bool"
)

// Validators is a struct for validator configuration.
type Validators struct {
	// ErrorMessage is a message to send if validation fails.
	ErrorMessage string `yaml:"error_message"`
	// Checks is a list of checks to perform.
	Checks []Check `yaml:"checks"`
}

func (v *Validators) validate() (errs []error) {
	for _, check := range v.Checks {
		switch check {
		case CheckNotEmpty, CheckIsInt, CheckIsFloat, CheckIsBool:
			return
		default:
			errs = append(errs, fmt.Errorf("unknown validator: %q", check))
		}
	}
	return
}
