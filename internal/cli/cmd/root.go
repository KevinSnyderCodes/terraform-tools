package cmd

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kevinsnydercodes/terraform-tools/internal/cli"
	"github.com/kevinsnydercodes/terraform-tools/internal/cli/value"
)

var (
	logLevel = value.String{
		Default:  "info",
		Env:      "LOG_LEVEL",
		Required: true,
		ValidateFunc: value.ValidateStringOneOf([]string{
			"trace",
			"debug",
			"info",
			"warning",
			"error",
			"fatal",
			"panic",
		}),
	}

	rootCmd = cobra.Command{
		Use:   "terraform-tools",
		Short: "Collection of commands for working with Terraform.",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Validate inputs
			inputs := map[string]cli.Validator{
				"log-level": &logLevel,
			}
			if err := cli.ValidateInputs(inputs); err != nil {
				return fmt.Errorf("error validating inputs: %w", err)
			}

			// Set log level
			level, err := log.ParseLevel(logLevel.Get())
			if err != nil {
				return fmt.Errorf("error parsing log level from input: %w", err)
			}
			log.SetLevel(level)

			return nil
		},
	}
)

func init() {
	rootCmd.PersistentFlags().VarP(&logLevel, "log-level", "", "Level to use for logging.")

	rootCmd.AddCommand(&docsCmd)
	rootCmd.AddCommand(&stateCmd)
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
