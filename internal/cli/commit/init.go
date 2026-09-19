package commit

import (
	"context"
	"fmt"
	"os"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap/commit"
	"github.com/juli3nk/git-bootstrap-cli/internal/config"
	"github.com/spf13/cobra"
)

func newInitCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Create bootstrap commits",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInit(cmd.Context(), templateDir())
		},
	}

	return cmd
}

func runInit(ctx context.Context, templateDir string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting current directory: %w", err)
	}

	manifest, err := config.Load(templateDir)
	if err != nil {
		return err
	}

	svc := commit.NewService()
	return svc.InitCommit(ctx, dir, manifest)
}
