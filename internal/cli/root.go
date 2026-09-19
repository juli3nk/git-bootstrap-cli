package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juli3nk/git-bootstrap-cli/internal/cli/commit"
	"github.com/juli3nk/git-bootstrap-cli/internal/cli/templates"

	"github.com/spf13/cobra"
)

var templateDir string

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "git-bootstrap",
		Short: "Bootstrap new Git repositories with standardized templates",
	}

	cmd.AddCommand(commit.NewCommand(TemplateDir))
	cmd.AddCommand(templates.NewCommand(TemplateDir))
	cmd.AddCommand(NewDoctorCommand())
	cmd.AddCommand(NewNewCommand())
	cmd.AddCommand(NewUpgradeCommand())
	cmd.AddCommand(NewVersionCommand())

	return cmd
}

func Execute() error {
	return NewCommand().Execute()
}

func init() {
	cobra.OnInitialize(resolveTemplateDir)
}

func resolveTemplateDir() {
	templateDir = resolveBootstrapDir("templates")
}

func resolveBootstrapDir(sub string) string {
	if d := os.Getenv("PROJECT_BOOTSTRAP_HOME"); d != "" {
		return filepath.Join(d, sub)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot determine home directory: %v\n", err)
		os.Exit(1)
	}

	return filepath.Join(home, ".local", "share", "git-bootstrap", sub)
}

func TemplateDir() string {
	return templateDir
}
