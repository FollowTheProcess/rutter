// Package app implements the functionality of the rutter CLI.
//
// The commands and subcommands simply delegate to the exposed members
// of this package.
package app

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"go.followtheprocess.codes/rutter/internal/db"
)

// IO holds the input/output streams for an App.
type IO struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// App represents the rutter program.
type App struct {
	io    IO
	store *db.Store
}

// New create a new [App].
func New(ctx context.Context, io IO, dbPath string) (App, error) {
	store, err := db.Open(ctx, dbPath)
	if err != nil {
		return App{}, err
	}

	return App{
		io:    io,
		store: store,
	}, nil
}

// DBPath returns the canonical path of the history SQLite DB, given
// the path of XDG_DATA_HOME.
func DBPath(dataHome string) string {
	return filepath.Join(dataHome, "rutter", "history.db")
}

// Close closes the DB connection.
func (a App) Close() error {
	return a.store.Close()
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
