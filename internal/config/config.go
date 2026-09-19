package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manifest is the parsed contents of templates.yml.
type Manifest struct {
	Init      InitManifest `yaml:"init"`
	AppCommit Commit       `yaml:"app_commit"`
}

// InitManifest describes the ordered bootstrap commits and the initial tag.
type InitManifest struct {
	Commits []Commit `yaml:"commits"`
	Tag     Tag      `yaml:"tag"`
}

// Commit describes a single commit message and the files/directories it covers.
type Commit struct {
	Message string   `yaml:"message"`
	Files   []string `yaml:"files"`
}

// Tag describes the initial tag to create.
type Tag struct {
	Name    string `yaml:"name"`
	Message string `yaml:"message"`
}

// Load reads and validates templates.yml from the given template directory.
func Load(templateDir string) (*Manifest, error) {
	path := filepath.Join(templateDir, "templates.yml")

	//nolint:gosec // path is resolved within the configured templates directory
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading templates manifest: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing templates manifest: %w", err)
	}

	if err := m.validate(); err != nil {
		return nil, err
	}

	return &m, nil
}

func (m *Manifest) validate() error {
	if m.Init.Tag.Name == "" {
		return fmt.Errorf("templates manifest: init.tag.name is required")
	}

	for i, c := range m.Init.Commits {
		if c.Message == "" {
			return fmt.Errorf("templates manifest: init.commits[%d].message is required", i)
		}
	}

	return nil
}
