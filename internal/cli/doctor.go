package cli

import (
	"fmt"

	"github.com/juli3nk/git-bootstrap-cli/internal/bootstrap"
	"github.com/spf13/cobra"
)

func NewDoctorCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "doctor [stack] [kind]",
		Short: "Check project files against templates",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runDoctor,
	}

	return cmd
}

func runDoctor(cmd *cobra.Command, args []string) error {
	svc := bootstrap.NewService()

	result, err := svc.Verify(TemplateDir(), args, ".")
	if err != nil {
		return err
	}

	for _, report := range result.Reports {
		switch report.Status {
		case bootstrap.StatusMissing:
			fmt.Printf("[MISSING] %s\n", report.Path)
		case bootstrap.StatusOK:
			fmt.Printf("[OK]      %s\n", report.Path)
		case bootstrap.StatusDiff:
			fmt.Printf("[DIFF]    %s\n", report.Path)
		}
	}

	if result.HasIssues() {
		return fmt.Errorf("doctor found %d missing and %d diff file(s)", result.Missing, result.Diff)
	}

	return nil
}
