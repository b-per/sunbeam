package utils

import (
	"os"
	"path"
	"strings"
)

func ResolvePath(filepath string) (string, error) {
	if strings.HasPrefix(filepath, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return strings.Replace(filepath, "~", homeDir, 1), nil
	}

	if !path.IsAbs(filepath) {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return path.Join(cwd, filepath), nil
	}

	return filepath, nil

}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)

	return !os.IsNotExist(err)
}
