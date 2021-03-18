package cmd

import (
	"fmt"
	"os"
	"path"
	"regexp"
	"time"

	terraformcli "github.com/kevinsnydercodes/go-terraform-cli"
	"github.com/kevinsnydercodes/terraform-tools/internal/cli"
	"github.com/kevinsnydercodes/terraform-tools/internal/cli/value"
	"github.com/kevinsnydercodes/terraform-tools/internal/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	stateCmd = cobra.Command{
		Use:   "state",
		Short: "Commands for working with state.",
	}
)

var (
	statemoveCmdSrc = value.String{
		Required: true,
	}
	statemoveCmdDst = value.String{
		Required: true,
	}

	statemoveCmd = cobra.Command{
		Use:   "move [flags] REGULAR_EXPRESSION",
		Short: "Move state from one workspace to another workspace.",
		Long:  "", // TODO: Add long documentation, describe process in detail
		Run: cli.Run(func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("expected 1 argument")
			}

			src := statemoveCmdSrc.Get()
			dst := statemoveCmdDst.Get()
			expr := args[0]

			log.Info("Compiling regular expression")
			re, err := regexp.Compile(expr)
			if err != nil {
				return fmt.Errorf("error compiling regular expression: %w", err)
			}

			log.Info("Creating Terraform CLI for source workspace")
			srcCLI, err := terraformcli.NewCLI(src)
			if err != nil {
				return fmt.Errorf("error creating terraform cli for source workspace: %w", err)
			}

			log.Info("Creating Terraform CLI for destination workspace")
			dstCLI, err := terraformcli.NewCLI(dst)
			if err != nil {
				return fmt.Errorf("error creating terraform cli for destination workspace: %w", err)
			}

			log.Info("Listing resources in source workspace")
			resources, err := srcCLI.State.List(terraformcli.ListOptions{})
			if err != nil {
				return fmt.Errorf("error listing resources in source workspace: %w", err)
			}

			log.Info("Filtering resources by regular expression")
			resources = types.StringSlice(resources).FilterRegexp(re)

			log.Info("Pulling state for source workspace")
			srcState, err := srcCLI.State.Pull(terraformcli.PullOptions{})
			if err != nil {
				return fmt.Errorf("error pulling state for source workspace: %w", err)
			}

			log.Info("Pulling state for destination workspace")
			dstState, err := dstCLI.State.Pull(terraformcli.PullOptions{})
			if err != nil {
				return fmt.Errorf("error pulling state for destination workspace: %w", err)
			}

			now := time.Now()

			// TODO: Handle files already existing

			log.Info("Writing state of source workspace to file")
			srcStateName := path.Join(src, "state.tfstate")
			if err := os.WriteFile(srcStateName, srcState, 0644); err != nil {
				return fmt.Errorf("error writing state file for source workspace: %w", err)
			}

			log.Info("Writing state of destination workspace to file")
			dstStateName := path.Join(dst, "state.tfstate")
			if err := os.WriteFile(dstStateName, dstState, 0644); err != nil {
				return fmt.Errorf("error writing state file for destination workspace: %w", err)
			}

			log.Info("Writing state of source workspace to backup file")
			srcStateBackupName := path.Join(src, fmt.Sprintf("backup-%s.tfstate", now.Format(time.RFC3339)))
			if err := os.WriteFile(srcStateBackupName, srcState, 0644); err != nil {
				return fmt.Errorf("error writing state backup file for source workspace: %w", err)
			}

			log.Info("Writing state of destination workspace to backup file")
			dstStateBackupName := path.Join(dst, fmt.Sprintf("backup-%s.tfstate", now.Format(time.RFC3339)))
			if err := os.WriteFile(dstStateBackupName, dstState, 0644); err != nil {
				return fmt.Errorf("error writing state backup file for source workspace: %w", err)
			}

			log.Info("Moving resources")
			for _, resource := range resources {
				log.Infof("Moving resource %s", resource)
				options := terraformcli.MoveOptions{
					State:    srcStateName,
					StateOut: dstStateName,
				}
				if err := dstCLI.State.Move(resource, resource, options); err != nil {
					return fmt.Errorf("error moving resource %s: %w", resource, err)
				}
			}

			log.Info("Pushing state for source workspace")
			if err := srcCLI.State.Push(srcStateName, terraformcli.PushOptions{}); err != nil {
				return fmt.Errorf("error pushing state for source workspace: %w", err)
			}

			log.Info("Pushing state for destination workspace")
			if err := dstCLI.State.Push(dstStateName, terraformcli.PushOptions{}); err != nil {
				return fmt.Errorf("error pushing state for destination workspace: %w", err)
			}

			return nil
		}),
	}
)

func init() {
	statemoveCmd.Flags().Var(&statemoveCmdSrc, "src", "Source workspace to move state from.")
	statemoveCmd.Flags().Var(&statemoveCmdDst, "dst", "Destination workspace to move state to.")

	stateCmd.AddCommand(&statemoveCmd)
}
