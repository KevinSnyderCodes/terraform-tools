package cli

import (
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"
)

type (
	// FuncRun is a function that a cobra.Command runs.
	FuncRun func(*cobra.Command, []string)
	// FuncRunE is a FuncRun that returns an error.
	FuncRunE func(*cobra.Command, []string) error
)

// Run wraps a FuncRunE with fatal logging via logrus. It returns a FuncRun that a cobra.Command can run.
func Run(fn FuncRunE) FuncRun {
	return func(cmd *cobra.Command, args []string) {
		if err := fn(cmd, args); err != nil {
			log.Fatal(err)
		}
	}
}
