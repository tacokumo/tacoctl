package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tacokumo/tacoctl/internal/version"
)

func New() *cobra.Command {
	c := &cobra.Command{
		Use:     "tacoctl",
		Version: version.Get().String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		SilenceUsage: true,
	}
	return c
}
