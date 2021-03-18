package shell

import (
	"fmt"
	"log"
	"os/exec"
)

// MustLookPath searches for an executable and exits fatally if not found.
func MustLookPath(file string) string {
	result, err := exec.LookPath(file)
	if err != nil {
		log.Fatalln(fmt.Errorf("error looking for %s executable in path: %w", file, err))
	}

	return result
}
