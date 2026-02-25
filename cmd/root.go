package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tacokumo/tacoctl/cmd/config"
	"github.com/tacokumo/tacoctl/internal/version"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:     "tacoctl",
		Version: version.Get().String(),
		Short:   "A command-line interface for the Portal API",
		Long:    "tacoctl is a CLI tool for interacting with the Portal API, providing kubectl-like experience for application management.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	// Add subcommands
	c.AddCommand(config.NewConfigCmd())

	return c
}
