package commit

import (
	"github.com/spf13/cobra"
)

// NewCommand creates the commit command group.
func NewCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit",
		Short: "Create project commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}

	cmd.AddCommand(
		newAppCommand(templateDir),
		newInitCommand(templateDir),
	)

	return cmd
}
