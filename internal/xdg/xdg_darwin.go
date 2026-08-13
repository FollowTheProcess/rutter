package xdg

import (
	"errors"
	"os"
	"path/filepath"
)

func dataHome() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", errors.New("$HOME is not defined")
	}

	return filepath.Join(home, "Library", "Application Support"), nil
}
