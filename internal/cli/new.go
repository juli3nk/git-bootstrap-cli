package cli

import (
	"fmt"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap"
	"github.com/spf13/cobra"
)

var license string

func NewNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [stack] [kind]",
		Short: "Create a new project",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runNew,
	}

	cmd.Flags().StringVar(&license, "license", "MIT", "License to include (MIT, APACHE-2.0)")

	return cmd
}

func runNew(cmd *cobra.Command, args []string) error {
	svc := bootstrap.NewService()

	applied, err := svc.Apply(TemplateDir(), args, license, ".")
	if err != nil {
		return err
	}

	for _, src := range applied {
		fmt.Printf("Applied %s (%s)\n", src.Path, src.Mode)
	}

	if license != "" {
		fmt.Printf("Added %s license\n", license)
	}

	return nil
}
