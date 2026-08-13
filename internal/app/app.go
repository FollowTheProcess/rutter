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
	io IO
}

type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// New create a new [App].
func New(io IO) App {
	return App{
		io: io,
	}
}

// Add implements the add subcommand.
//
// It inserts the command into the history DB.
func (a App) Add(command string) error {
	fmt.Fprintf(a.io.Stdout, "Adding command '%s'\n", command)

	return nil
}

// Finish implements the finish subcommand.
func (a App) Finish(id, exit int, duration time.Duration) error {
	fmt.Fprintf(a.io.Stdout, "Finishing command %d with exit %d and duration %s\n", id, exit, duration)

	return nil
}

// Init implements the init subcommand.
func (a App) Init(shell string) error {
	fmt.Fprintf(a.io.Stdout, "Printing shell script for %s\n", shell)

	return nil
}
