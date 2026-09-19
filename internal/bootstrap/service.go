package bootstrap

import (
	"bytes"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// Service provides project bootstrapping operations.
type Service struct{}

// NewService creates a new bootstrap service.
func NewService() *Service {
	return &Service{}
}

// Apply resolves and applies the template sources to the destination directory.
// It returns the list of sources that were effectively applied. If license is
// non-empty, the matching license file is copied into "dst/LICENSE".
func (s *Service) Apply(templateDir string, args []string, license string, dst string) ([]AppliedSource, error) {
	sources, err := ResolveSources(templateDir, args)
	if err != nil {
		return nil, err
	}

	return s.ApplySources(templateDir, sources, license, dst)
}

// ApplySources applies an already-resolved list of template sources to the
// destination directory. It returns the list of sources that were effectively
// applied. If license is non-empty, the matching license file is copied into
// "dst/LICENSE".
func (s *Service) ApplySources(templateDir string, sources []Source, license string, dst string) ([]AppliedSource, error) {
	var applied []AppliedSource

	for _, src := range sources {
		switch src.Mode {
		case ModeCopy:
			if err := CopyDir(src.Path, dst, "append"); err != nil {
				return applied, fmt.Errorf("copying %s templates: %w", src.Path, err)
			}
		case ModeAppend:
			if err := AppendFiles(src.Path, dst); err != nil {
				return applied, fmt.Errorf("appending %s files: %w", src.Path, err)
			}
		}
		applied = append(applied, AppliedSource(src))
	}

	if license != "" {
		licensesDir := filepath.Join(templateDir, "licenses")

		licenseSrc, err := ValidateLicense(licensesDir, license)
		if err != nil {
			return applied, err
		}

		if err := CopyFile(licenseSrc, filepath.Join(dst, "LICENSE")); err != nil {
			return applied, fmt.Errorf("copying license: %w", err)
		}
	}

	return applied, nil
}

// Verify compares the files produced by the resolved template sources against
// the destination directory and returns a summary result.
func (s *Service) Verify(templateDir string, args []string, dst string) (*VerifyResult, error) {
	sources, err := ResolveSources(templateDir, args)
	if err != nil {
		return nil, err
	}

	files := collectFiles(sources)
	result := &VerifyResult{}

	for _, rel := range files {
		expected := composeExpected(sources, rel)
		actual, err := readFileOrError(filepath.Join(dst, rel))
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", rel, err)
		}

		var status FileStatus
		switch {
		case actual == nil:
			status = StatusMissing
		case bytes.Equal(expected, actual):
			status = StatusOK
		default:
			status = StatusDiff
		}

		result.Reports = append(result.Reports, FileReport{Path: rel, Status: status})

		switch status {
		case StatusOK:
			result.OK++
		case StatusDiff:
			result.Diff++
		case StatusMissing:
			result.Missing++
		}
	}

	return result, nil
}

// collectFiles returns the relative file paths produced by all sources,
// without duplicates. Copy sources skip any "append" subdirectory because those
// files are handled by a separate append source.
func collectFiles(sources []Source) []string {
	seen := make(map[string]struct{})
	var files []string

	for _, s := range sources {
		_ = filepath.Walk(s.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if info.IsDir() {
				rel, _ := filepath.Rel(s.Path, path)
				if s.Mode == ModeCopy && rel == "append" {
					return filepath.SkipDir
				}
				return nil
			}

			rel, err := filepath.Rel(s.Path, path)
			if err != nil {
				return nil
			}

			if _, ok := seen[rel]; ok {
				return nil
			}
			seen[rel] = struct{}{}
			files = append(files, rel)

			return nil
		})
	}

	return files
}

// composeExpected builds the expected content for a target file by combining the
// contributions of each source in order. Copy sources overwrite the previous
// content; append sources add their content with a leading blank line when
// previous content exists.
func composeExpected(sources []Source, relPath string) []byte {
	var result []byte

	for _, s := range sources {
		src := filepath.Join(s.Path, relPath)
		//nolint:gosec // source path is resolved within the configured templates directory
		content, err := os.ReadFile(src)
		if err != nil {
			continue
		}

		switch s.Mode {
		case ModeCopy:
			result = content
		case ModeAppend:
			if len(result) > 0 {
				result = append(result, '\n')
			}
			result = append(result, content...)
		}
	}

	return result
}

func readFileOrError(path string) ([]byte, error) {
	//nolint:gosec // path is part of the project being verified
	content, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	return content, nil
}
