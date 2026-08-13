// Package xdg implements convenience functions for returning common XDG
// directories, with platform-specific fallbacks.
//
// If the XDG env var for that directory type is set, it will use that,
// otherwise it will use the fallback for the platform.
package xdg

import (
	"errors"
	"os"
	"path/filepath"
)

// ConfigHome returns the path for XDG_CONFIG_HOME.
//
// If the XDG_CONFIG_HOME env var is not set, it returns
// [os.UserConfigDir].
func ConfigHome() (string, error) {
	if dir := os.Getenv("XDG_CONFIG_HOME"); dir != "" {
		return dir, nil
	}

	return os.UserConfigDir()
}

// DataHome returns the path for XDG_DATA_HOME.
//
// If the XDG_DATA_HOME env var is not set, it returns
// a per-platform fallback:
//
//   - Unix: ~/.local/share
//   - MacOS: ~/Library/Application Support
//   - Windows: LocalAppData
func DataHome() (string, error) {
	if dir := os.Getenv("XDG_DATA_HOME"); dir != "" {
		if !filepath.IsAbs(dir) {
			return "", errors.New("path in $XDG_DATA_HOME is relative")
		}

		return dir, nil
	}

	return dataHome()
}
