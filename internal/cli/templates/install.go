package templates

import (
	"context"
	"fmt"
	"os"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap/templates"
	"github.com/spf13/cobra"
)

func newInstallCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install <repo-url>",
		Short: "Install bootstrap templates from a Git repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesInstall(cmd.Context(), templateDir(), args[0])
		},
	}

	return cmd
}

func runTemplatesInstall(ctx context.Context, target, repoURL string) error {
	svc := templates.NewService()

	if _, err := os.Stat(target); err == nil {
		fmt.Printf("templates directory already exists: %s\n", target)
		fmt.Println("run `git-bootstrap templates update` to update it")
		return nil
	}

	fmt.Printf("Installing templates from %s into %s\n", repoURL, target)

	if err := svc.Install(ctx, repoURL, target); err != nil {
		return err
	}

	fmt.Printf("Templates installed at %s\n", target)
	return nil
}
