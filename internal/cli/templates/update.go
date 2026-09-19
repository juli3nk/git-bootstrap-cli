package templates

import (
	"context"
	"fmt"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap/templates"
	"github.com/spf13/cobra"
)

func newUpdateCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update installed bootstrap templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesUpdate(cmd.Context(), templateDir())
		},
	}

	return cmd
}

func runTemplatesUpdate(ctx context.Context, target string) error {
	svc := templates.NewService()

	fmt.Println("Updating bootstrap templates...")

	version, err := svc.Update(ctx, target)
	if err != nil {
		return err
	}

	fmt.Println("\nCurrent version:")
	fmt.Println(version)
	return nil
}
