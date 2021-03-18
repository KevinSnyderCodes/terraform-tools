package value

import (
	"fmt"
	"os"
)

type ValidateStringFunc func(string) error

// String is a string value with various options.
type String struct {
	value string

	Default      string
	Env          string
	Required     bool
	ValidateFunc ValidateStringFunc
}

// Get gets the value.
func (o *String) Get() string {
	if o.value != "" {
		return o.value
	}

	if o.Env != "" {
		if value := os.Getenv(o.Env); value != "" {
			return value
		}
	}

	if o.Default != "" {
		return o.Default
	}

	return ""
}

// Set sets the value.
func (o *String) Set(value string) error {
	o.value = value
	return nil
}

// Validate validates a value.
func (o *String) Validate() error {
	value := o.Get()

	if o.Required && value == "" {
		return fmt.Errorf("must provide value")
	}

	if o.ValidateFunc != nil {
		if err := o.ValidateFunc(value); err != nil {
			return err
		}
	}

	return nil
}

// String returns the value.
func (o *String) String() string {
	return o.Get()
}

// Type returns the type.
func (o *String) Type() string {
	return "string"
}
