package xdg

import (
	"errors"
	"os"
)

func dataHome() (string, error) {
	dir := os.Getenv("LocalAppData")
	if dir == "" {
		return "", errors.New("%LocalAppData% is not defined")
	}

	return dir, nil
}
