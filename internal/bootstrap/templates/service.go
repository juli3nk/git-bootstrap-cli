package templates

import (
	"context"
	"fmt"
	"os"

	"github.com/juli3nk/git-bootstrap-cli/internal/git"
)

// Service manages the bootstrap templates repository.
type Service struct{}

// NewService creates a new templates service.
func NewService() *Service {
	return &Service{}
}

// Install clones a Git repository into the target directory.
func (s *Service) Install(ctx context.Context, repoURL, target string) error {
	if _, err := os.Stat(target); err == nil {
		return fmt.Errorf("templates directory already exists: %s", target)
	}

	if err := os.MkdirAll(target, 0o750); err != nil {
		return fmt.Errorf("creating templates directory: %w", err)
	}

	if err := git.Clone(ctx, target, repoURL); err != nil {
		return fmt.Errorf("cloning templates: %w", err)
	}

	return nil
}

// StatusInfo describes the status of an installed templates repository.
type StatusInfo struct {
	Directory    string
	Installed    bool
	IsGitRepo    bool
	Version      string
	VersionError error
	RemoteOrigin string
	RemoteError  error
}

// Status returns the current status of the installed templates.
func (s *Service) Status(ctx context.Context, target string) (*StatusInfo, error) {
	info := &StatusInfo{Directory: target}

	if _, err := os.Stat(target); err != nil {
		info.Installed = false
		return info, nil
	}

	info.Installed = true

	if !git.IsRepo(ctx, target) {
		info.IsGitRepo = false
		return info, nil
	}

	info.IsGitRepo = true

	output, err := git.Run(ctx, target, "describe", "--tags", "--always")
	if err != nil {
		info.VersionError = err
	} else {
		info.Version = output
	}

	remote, err := git.RemoteURL(ctx, target, "origin")
	if err != nil {
		info.RemoteError = err
	} else {
		info.RemoteOrigin = remote
	}

	return info, nil
}

// Update pulls the latest changes from origin/main and returns the new version.
func (s *Service) Update(ctx context.Context, target string) (string, error) {
	if _, err := os.Stat(target); err != nil {
		return "", fmt.Errorf("templates directory not found: %s", target)
	}

	if !git.IsRepo(ctx, target) {
		return "", fmt.Errorf("not a git repository: %s", target)
	}

	if _, err := git.Run(ctx, target, "pull", "--ff-only", "origin", "main"); err != nil {
		return "", fmt.Errorf("pulling updates: %w", err)
	}

	output, err := git.Run(ctx, target, "describe", "--tags", "--always")
	if err != nil {
		return "", fmt.Errorf("describing version: %w", err)
	}

	return output, nil
}
