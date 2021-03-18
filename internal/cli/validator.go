package cli

import "fmt"

// Validator is an interface for something that can be validated.
type Validator interface {
	Validate() error
}

// ValidateInputs validates a map of named inputs using the Validator interface.
func ValidateInputs(inputs map[string]Validator) error {
	for name, validator := range inputs {
		if err := validator.Validate(); err != nil {
			return fmt.Errorf("error with input \"%s\": %w", name, err)
		}
	}
	return nil
}
