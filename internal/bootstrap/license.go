package bootstrap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AvailableLicenses returns the license names available in the given directory.
func AvailableLicenses(dir string) []string {
	var licenses []string

	entries, err := os.ReadDir(dir)
	if err != nil {
		return licenses
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if filepath.Ext(name) == ".txt" {
			licenses = append(licenses, strings.TrimSuffix(name, ".txt"))
		}
	}

	return licenses
}

// ValidateLicense checks whether the requested license exists in the licenses
// directory and returns the path to its template file.
func ValidateLicense(licensesDir, license string) (string, error) {
	available := AvailableLicenses(licensesDir)
	if len(available) == 0 {
		return "", fmt.Errorf("no license templates found in %s", licensesDir)
	}

	valid := false
	for _, l := range available {
		if l == license {
			valid = true
			break
		}
	}
	if !valid {
		return "", fmt.Errorf("unknown license %q, available: %s", license, strings.Join(available, ", "))
	}

	return filepath.Join(licensesDir, license+".txt"), nil
}
