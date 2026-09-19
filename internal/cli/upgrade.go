package cli

import (
	"fmt"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap"
	"github.com/spf13/cobra"
)

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade [stack] [kind]",
		Short: "Upgrade project files from templates",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runUpgrade,
	}

	return cmd
}

func runUpgrade(cmd *cobra.Command, args []string) error {
	svc := bootstrap.NewService()

	sources, err := bootstrap.ResolveSources(TemplateDir(), args)
	if err != nil {
		return err
	}

	fmt.Printf("Applying %d template source(s)...\n", len(sources))

	for _, src := range sources {
		switch src.Mode {
		case bootstrap.ModeCopy:
			fmt.Printf("Copying %s\n", src.Path)
		case bootstrap.ModeAppend:
			fmt.Printf("Appending %s\n", src.Path)
		}
	}

	if _, err := svc.ApplySources(TemplateDir(), sources, "", "."); err != nil {
		return err
	}

	fmt.Println("\nUpgrade complete.")

	return nil
}
