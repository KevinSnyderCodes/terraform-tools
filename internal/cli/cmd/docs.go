package cmd

import (
	"fmt"

	"github.com/pkg/browser"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/kevinsnydercodes/terraform-tools/internal/cli"
	"github.com/kevinsnydercodes/terraform-tools/internal/cli/value"
)

var (
	docsCmdType = value.String{
		Default:  "resource",
		Required: true,
		ValidateFunc: value.ValidateStringOneOf([]string{
			"resource",
			"data-source",
		}),
	}
	docsCmdVersion string

	docsCmd = cobra.Command{
		Use:   "docs PROVIDER RESOURCE",
		Short: "View the documentation for a Terraform resource.",
		Run: cli.Run(func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("expected 2 arguments")
			}

			provider := args[0]
			resource := args[1]

			url := fmt.Sprintf(
				"https://registry.terraform.io/providers/hashicorp/%s/%s/docs/%s/%s",
				provider,
				docsCmdVersion,
				docsCmdType.Get()+"s",
				resource,
			)

			// TODO: Resource does not always map one-to-one to documentation URL
			// Find a better way to construct the documentation URL
			// (Example: `aws_vpc` is just `vpc` in the URL)

			// Visit this page and watch network tab:
			// https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/vpc_dhcp_options
			//
			// API calls:
			// 1. https://registry.terraform.io/v2/providers (get provider ID)
			// 2. https://registry.terraform.io/v2/providers/[provider_id] (get provider version)
			// 3. https://registry.terraform.io/v2/provider_versions/[provider_version_id] (get ???)
			// 4. https://registry.terraform.io/v2/provider_docs (get resource ID)
			// 5. https://registry.terraform.io/v2/provider_docs/[resource_id] (get docs)
			// 6. Last API call includes "path" field with path to markdown file -- can parse this to build URL
			//    ("path" uses /r/ or /d/ for resources and data sources, and contains appropriate resource name for URL)

			log.Info("Opening web browser...")
			if err := browser.OpenURL(url); err != nil {
				return fmt.Errorf("error opening browser: %w", err)
			}

			return nil
		}),
	}
)

func init() {
	docsCmd.Flags().Var(&docsCmdType, "type", "Type of documentation.")
	docsCmd.Flags().StringVar(&docsCmdVersion, "version", "latest", "Version of the provider.")
}
