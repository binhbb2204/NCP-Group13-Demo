package utils

import (
	"os"
	"path/filepath"
)

// GetProjectConfigPath returns a project-local path for storing config or runtime files.
// It uses the current working directory rather than the user's home directory.
func GetProjectConfigPath(name string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Join(cwd, name), nil
}
