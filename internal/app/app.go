// Package app implements the functionality of the rutter CLI.
//
// The commands and subcommands simply delegate to the exposed members
// of this package.
package app

import (
	"fmt"
	"io"
	"time"
)

// App represents the rutter program.
type App struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// New create a new [App].
func New(stdin io.Reader, stdout, stderr io.Writer) App {
	// TODO: This will obviously need the DB and all that
	return App{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Add implements the add subcommand.
//
// It inserts the command into the history DB.
func (a App) Add(command string) error {
	fmt.Fprintf(a.stdout, "Adding command '%s'\n", command)

	return nil
}

// Finish implements the finish subcommand.
func (a App) Finish(id, exit int, duration time.Duration) error {
	fmt.Fprintf(a.stdout, "Finishing command %d with exit %d and duration %s\n", id, exit, duration)

	return nil
}
