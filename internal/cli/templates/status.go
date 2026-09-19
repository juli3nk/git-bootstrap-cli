package templates

import (
	"context"
	"fmt"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap/templates"
	"github.com/spf13/cobra"
)

func newStatusCommand(templateDir func() string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show bootstrap templates status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runTemplatesStatus(cmd.Context(), templateDir())
		},
	}

	return cmd
}

func runTemplatesStatus(ctx context.Context, target string) error {
	svc := templates.NewService()

	fmt.Printf("Templates directory: %s\n", target)

	info, err := svc.Status(ctx, target)
	if err != nil {
		return err
	}

	if !info.Installed {
		fmt.Println("Status: not installed")
		return nil
	}

	fmt.Println("Status: directory exists")

	if !info.IsGitRepo {
		fmt.Println("Git repository: no")
		return nil
	}

	fmt.Println("Git repository: yes")

	if info.VersionError != nil {
		fmt.Printf("Version: unknown (%v)\n", info.VersionError)
	} else {
		fmt.Printf("Version: %s\n", info.Version)
	}

	if info.RemoteError != nil {
		fmt.Printf("Remote origin: unknown (%v)\n", info.RemoteError)
	} else {
		fmt.Printf("Remote origin: %s\n", info.RemoteOrigin)
	}

	return nil
}
