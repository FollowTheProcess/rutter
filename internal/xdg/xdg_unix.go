//go:build unix && !darwin

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

	return filepath.Join(home, ".local", "share"), nil
}
