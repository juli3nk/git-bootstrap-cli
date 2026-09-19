package commit

import (
	"context"
	"fmt"
	"os"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap/commit"
	"github.com/juli3nk/git-bootstrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func newAppCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "app",
		Short: "Create application commit",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(cmd.Context(), templateDir())
		},
	}

	return cmd
}

func runApp(ctx context.Context, templateDir string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	manifest, err := config.Load(templateDir)
	if err != nil {
		return err
	}

	svc := commit.NewService()
	return svc.AppCommit(ctx, dir, manifest)
}
