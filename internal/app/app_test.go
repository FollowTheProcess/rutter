package app_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"go.followtheprocess.codes/rutter/internal/app"
	"go.followtheprocess.codes/test"
)

func TestAddFinish(t *testing.T) {
	stdin := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}

	appIO := app.IO{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}

	dir := t.TempDir()
	sessionID := "49583c11-0a56-4750-906e-afe6065a5e47"

	// db should be elsewhere
	dbPath := filepath.Join(t.TempDir(), "test.db")

	app, err := app.New(t.Context(), appIO, dir, sessionID, dbPath)
	test.Ok(t, err)
	t.Cleanup(func() { app.Close() })

	id, err := app.Add(t.Context(), "echo hello")
	test.Ok(t, err)

	// The DB should have been created
	_, err = os.Stat(dbPath)
	test.Ok(t, err)

	// Finish should be able to finalise the same command
	err = app.Finish(t.Context(), id, 0, 69*time.Millisecond)
	test.Ok(t, err)

	// Stdout and Stderr should be empty as they will interfere with the TUI
	test.Equal(t, stdout.Len(), 0, test.Context("stdout was written to: %s", stdout.String()))
	test.Equal(t, stderr.Len(), 0, test.Context("stderr was written to: %s", stderr.String()))

	// TODO: When we can search for history entries, use that here in the test
}
