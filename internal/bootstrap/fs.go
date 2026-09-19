package bootstrap

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CopyDir copies a directory tree from src to dst, preserving permissions.
// Any top-level subdirectory names listed in excludes are skipped.
func CopyDir(src, dst string, excludes ...string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		if rel == "." {
			return nil
		}

		first := strings.SplitN(rel, string(filepath.Separator), 2)[0]
		for _, ex := range excludes {
			if first == ex {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		return CopyFile(path, target)
	})
}

// AppendFiles appends files from src into existing files at dst. If the target
// file does not exist it is created. Duplicate content is not appended twice.
func AppendFiles(src, dst string) (retErr error) {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if err := EnsureDirExists(target); err != nil {
			return err
		}

		//nolint:gosec // target path is part of the bootstrap destination
		existing, err := os.ReadFile(target)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		//nolint:gosec // path is resolved within the configured templates directory
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		if bytes.Equal(existing, content) || bytes.Contains(existing, content) {
			return nil
		}

		mode := info.Mode()
		if mode == 0 {
			mode = 0o644
		}

		//nolint:gosec // target path is part of the bootstrap destination
		f, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, mode)
		if err != nil {
			return err
		}
		defer func() {
			if err := f.Close(); err != nil && retErr == nil {
				retErr = fmt.Errorf("closing %s: %w", target, err)
			}
		}()

		if len(existing) > 0 {
			if _, err := f.WriteString("\n"); err != nil {
				return err
			}
		}

		_, err = f.Write(content)
		return err
	})
}

// CopyFile copies a single file from src to dst, preserving permissions.
// .gitkeep files are ignored.
func CopyFile(src, dst string) (retErr error) {
	if filepath.Base(src) == ".gitkeep" {
		return nil
	}

	if err := EnsureDirExists(dst); err != nil {
		return err
	}

	//nolint:gosec // source path is resolved within the configured templates directory
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}

	defer func() {
		if err := srcFile.Close(); err != nil && retErr == nil {
			retErr = fmt.Errorf("closing source %s: %w", src, err)
		}
	}()

	//nolint:gosec // destination path is produced by the bootstrap operation
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		if err := dstFile.Close(); err != nil && retErr == nil {
			retErr = fmt.Errorf("closing destination %s: %w", dst, err)
		}
	}()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// EnsureDirExists creates the parent directory of the given path if it does
// not exist.
func EnsureDirExists(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o750); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}
	return nil
}
