package commit

import (
	"context"
	"fmt"
	"os"

	"github.com/juli3nk/git-bootstrap-cli/internal/config"
	"github.com/juli3nk/git-bootstrap-cli/internal/git"
)

// Service provides commit creation operations for bootstrapped projects.
type Service struct{}

// NewService creates a new commit service.
func NewService() *Service {
	return &Service{}
}

// InitCommit initializes a Git repository and creates the configured bootstrap
// commits and initial tag.
func (s *Service) InitCommit(ctx context.Context, dir string, manifest *config.Manifest) error {
	if !git.IsRepo(ctx, dir) {
		if err := git.Init(ctx, dir, "develop"); err != nil {
			return fmt.Errorf("initializing repository: %w", err)
		}
	}

	if git.HasCommits(ctx, dir) {
		return fmt.Errorf("repository already contains commits")
	}

	for _, commit := range manifest.Init.Commits {
		if err := addAndCommitIfAnyExist(ctx, dir, commit.Files, commit.Message); err != nil {
			return err
		}
	}

	if git.HasCommits(ctx, dir) && !git.TagExists(ctx, dir, manifest.Init.Tag.Name) {
		if err := git.Tag(ctx, dir, manifest.Init.Tag.Name, manifest.Init.Tag.Message); err != nil {
			return fmt.Errorf("creating initial tag: %w", err)
		}
	}

	return nil
}

// AppCommit creates an application commit if there are staged changes.
func (s *Service) AppCommit(ctx context.Context, dir string, manifest *config.Manifest) error {
	return git.CommitIfStaged(ctx, dir, manifest.AppCommit.Message)
}

func addAndCommitIfAnyExist(ctx context.Context, dir string, files []string, message string) error {
	staged := false
	for _, f := range files {
		if _, err := os.Stat(f); err == nil {
			if err := git.Add(ctx, dir, f); err != nil {
				return fmt.Errorf("adding %s: %w", f, err)
			}
			staged = true
		}
	}

	if staged {
		if err := git.CommitIfStaged(ctx, dir, message); err != nil {
			return fmt.Errorf("creating commit %q: %w", message, err)
		}
	}

	return nil
}
