package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
)

func Run(ctx context.Context, dir string, args ...string) (string, error) {
	//nolint:gosec // git arguments are controlled by the caller, not user input
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}

func IsRepo(ctx context.Context, dir string) bool {
	_, err := Run(ctx, dir, "rev-parse", "--git-dir")
	return err == nil
}

func HasCommits(ctx context.Context, dir string) bool {
	_, err := Run(ctx, dir, "rev-parse", "--verify", "HEAD")
	return err == nil
}

func Init(ctx context.Context, dir, branch string) error {
	_, err := Run(ctx, dir, "init", "--initial-branch="+branch)
	return err
}

func Clone(ctx context.Context, dir, url string) error {
	_, err := Run(ctx, "", "clone", url, dir)
	return err
}

func Add(ctx context.Context, dir string, paths ...string) error {
	args := append([]string{"add"}, paths...)
	_, err := Run(ctx, dir, args...)
	return err
}

func Commit(ctx context.Context, dir, message string) error {
	_, err := Run(ctx, dir, "commit", "-m", message)
	return err
}

func DiffCachedQuiet(ctx context.Context, dir string) bool {
	_, err := Run(ctx, dir, "diff", "--cached", "--quiet")
	return err == nil
}

func CommitIfStaged(ctx context.Context, dir, message string) error {
	if !DiffCachedQuiet(ctx, dir) {
		return Commit(ctx, dir, message)
	}
	return nil
}

func AddAndCommitIfStaged(ctx context.Context, dir, path, message string) error {
	if err := Add(ctx, dir, path); err != nil {
		return err
	}
	return CommitIfStaged(ctx, dir, message)
}

func TagExists(ctx context.Context, dir, tag string) bool {
	_, err := Run(ctx, dir, "rev-parse", "--verify", "refs/tags/"+tag)
	return err == nil
}

func Tag(ctx context.Context, dir, tag, message string) error {
	_, err := Run(ctx, dir, "tag", "-a", tag, "-m", message)
	return err
}

func RemoteURL(ctx context.Context, dir, name string) (string, error) {
	return Run(ctx, dir, "remote", "get-url", name)
}
