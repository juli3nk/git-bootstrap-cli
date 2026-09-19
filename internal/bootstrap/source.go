package bootstrap

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Mode describes how a source directory should be applied to a project.
type Mode string

const (
	ModeCopy   Mode = "copy"
	ModeAppend Mode = "append"
)

// Source describes a single template source directory and how it should be
// applied to the current project.
type Source struct {
	Path string
	Mode Mode
}

// ResolveSources returns the ordered list of template sources for the given
// arguments.
//
//   - no args: common/
//   - [stack]: common/ + stack/base/ + stack/base/append/
//   - [stack, kind]: ... + stack/kind/ + stack/kind/append/
func ResolveSources(templateDir string, args []string) ([]Source, error) {
	common := filepath.Join(templateDir, "common")
	if err := requireDir(common); err != nil {
		return nil, fmt.Errorf("templates not found: %w", err)
	}

	sources := []Source{{Path: common, Mode: ModeCopy}}

	if len(args) == 0 {
		return sources, nil
	}

	stack := args[0]
	stackDir := resolveStackDir(templateDir, stack)
	if stackDir == "" {
		return nil, fmt.Errorf("stack template not found: %s", stack)
	}

	baseDir := filepath.Join(stackDir, "base")
	if err := requireDir(baseDir); err != nil {
		available := AvailableKinds(stackDir)
		return nil, fmt.Errorf("base template not found in stack %q; available kinds: %s", stack, strings.Join(available, ", "))
	}

	sources = append(sources, Source{Path: baseDir, Mode: ModeCopy})
	sources = appendAppendSource(sources, baseDir)

	if len(args) > 1 {
		kind := args[1]
		if kind == "base" {
			return nil, fmt.Errorf("%q is not a valid kind", kind)
		}

		kindDir := filepath.Join(stackDir, kind)
		if err := requireDir(kindDir); err != nil {
			available := AvailableKinds(stackDir)
			return nil, fmt.Errorf("kind %q not found in stack %q; available: %s", kind, stack, strings.Join(available, ", "))
		}

		sources = append(sources, Source{Path: kindDir, Mode: ModeCopy})
		sources = appendAppendSource(sources, kindDir)
	}

	return sources, nil
}

// resolveStackDir returns the directory for a stack, preferring the root of the
// template directory and falling back to the "stacks" subdirectory.
func resolveStackDir(templateDir, stack string) string {
	candidates := []string{
		filepath.Join(templateDir, stack),
		filepath.Join(templateDir, "stacks", stack),
	}

	for _, dir := range candidates {
		info, err := os.Stat(dir)
		if err == nil && info.IsDir() {
			return dir
		}
	}

	return ""
}

// appendAppendSource adds an append source for the given directory if its
// "append" subdirectory exists.
func appendAppendSource(sources []Source, dir string) []Source {
	appendDir := filepath.Join(dir, "append")
	if info, err := os.Stat(appendDir); err == nil && info.IsDir() {
		return append(sources, Source{Path: appendDir, Mode: ModeAppend})
	}
	return sources
}

// requireDir returns an error if the given path is not an existing directory.
func requireDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return fmt.Errorf("directory does not exist: %s", path)
		}
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("not a directory: %s", path)
	}
	return nil
}

// AvailableKinds returns the top-level kind directories inside a stack dir,
// excluding "append" and "base".
func AvailableKinds(stackDir string) []string {
	var kinds []string

	entries, err := os.ReadDir(stackDir)
	if err != nil {
		return kinds
	}

	for _, entry := range entries {
		if !entry.IsDir() || entry.Name() == "append" || entry.Name() == "base" {
			continue
		}
		kinds = append(kinds, entry.Name())
	}

	return kinds
}
