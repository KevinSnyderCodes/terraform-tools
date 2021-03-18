package value

import (
	"fmt"
	"strings"

	"github.com/kevinsnydercodes/terraform-tools/internal/types"
)

func ValidateStringOneOf(valid []string) ValidateStringFunc {
	return func(value string) error {
		if !types.StringSlice(valid).Contains(value) {
			return fmt.Errorf("Must provide one of: %s", strings.Join(valid, ", "))
		}
		return nil
	}
}
