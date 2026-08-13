package xdg_test

import (
	"os"
	"testing"

	"go.followtheprocess.codes/rutter/internal/xdg"
	"go.followtheprocess.codes/test"
)

func TestDataHome(t *testing.T) {
	// Isolate it from any actual set env vars
	prev, set := os.LookupEnv("XDG_DATA_HOME")
	if set {
		t.Cleanup(func() {
			// Put original value back
			//nolint:usetesting // This must be os
			os.Setenv("XDG_DATA_HOME", prev)
		})

		err := os.Unsetenv("XDG_DATA_HOME")
		test.Ok(t, err)
	}

	// Data home on windows is %LocalAppData%, not the config dir
	realDataHome := os.Getenv("LocalAppData")

	tests := []struct {
		name    string            // Name of the test case
		env     map[string]string // Env vars to set
		want    string            // Expected return path
		wantErr bool              // Whether we want an error
	}{
		{
			name:    "no env var uses windows default",
			env:     map[string]string{},
			want:    realDataHome,
			wantErr: false,
		},
		{
			name: "env var takes precedence",
			env: map[string]string{
				"XDG_DATA_HOME": `C:\some\dir`,
			},
			want:    `C:\some\dir`,
			wantErr: false,
		},
		{
			name: "relative path returns error",
			env: map[string]string{
				"XDG_DATA_HOME": `some\relative\dir`,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "no LocalAppData returns error",
			env: map[string]string{
				"LocalAppData": "",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			got, err := xdg.DataHome()
			test.WantErr(t, err, tt.wantErr)

			test.Equal(t, got, tt.want)
		})
	}
}
