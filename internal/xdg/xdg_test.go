package xdg_test

import (
	"os"
	"testing"

	"go.followtheprocess.codes/rutter/internal/xdg"
	"go.followtheprocess.codes/test"
)

func TestConfigHome(t *testing.T) {
	// Isolate it from any actual set env vars
	prev, set := os.LookupEnv("XDG_CONFIG_HOME")
	if set {
		t.Cleanup(func() {
			// Put original value back
			//nolint:usetesting // This must be os
			os.Setenv("XDG_CONFIG_HOME", prev)
		})

		err := os.Unsetenv("XDG_CONFIG_HOME")
		test.Ok(t, err)
	}

	realConfigHome, err := os.UserConfigDir()
	test.Ok(t, err)

	tests := []struct {
		name string            // Name of the test case
		env  map[string]string // env vars to set for the test
		want string            // Expected config home directory
	}{
		{
			name: "no env var uses std lib",
			env:  map[string]string{},
			want: realConfigHome,
		},
		{
			name: "env var takes precedence",
			env: map[string]string{
				"XDG_CONFIG_HOME": "/some/dir",
			},
			want: "/some/dir",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			got, err := xdg.ConfigHome()
			test.Ok(t, err)

			test.Equal(t, got, tt.want)
		})
	}
}
