package templates

import (
	"github.com/spf13/cobra"
)

// NewCommand creates the templates command group.
func NewCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "templates",
		Short: "Manage bootstrap templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		newInstallCommand(templateDir),
		newStatusCommand(templateDir),
		newUpdateCommand(templateDir),
	)

	return cmd
}
