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
	io        IO
	store     *db.Store
	cwd       string
	sessionID string
}

// New creates a new [App].
func New(ctx context.Context, io IO, cwd, sessionID, dbPath string) (App, error) {
	store, err := db.Open(ctx, dbPath)
	if err != nil {
		return App{}, err
	}

	return App{
		io:        io,
		store:     store,
		cwd:       cwd,
		sessionID: sessionID,
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
// It inserts the command into the history DB, and returns
// its ID.
func (a App) Add(ctx context.Context, command string) (int, error) {
	id, err := a.store.Query.StartCommand(ctx, db.StartCommandParams{
		Cmd:       command,
		Cwd:       a.cwd,
		Session:   a.sessionID,
		StartedAt: time.Now().UTC(),
	})
	if err != nil {
		return 0, fmt.Errorf("failed to start command: %w", err)
	}

	return int(id), nil
}

// Finish implements the finish subcommand.
func (a App) Finish(ctx context.Context, id, exit int, duration time.Duration) error {
	err := a.store.Query.FinishCommand(ctx, db.FinishCommandParams{
		Duration: duration,
		Exit:     int64(exit),
		ID:       int64(id),
	})
	if err != nil {
		return fmt.Errorf("failed to finish command: %w", err)
	}

	return nil
}

// SearchOptions holds the flags for the search command.
type SearchOptions struct {
	Directory string // Filter results to this directory
	Session   string // Filter results to this session ID
	Limit     int    // Maximum number of results to return
}

func (a App) Search(ctx context.Context, query string, options SearchOptions) error {
	// TODO: This
	return nil
}

// Suggest returns the best match for a command based on the text, used for
// shell auto-suggestion.
func (a App) Suggest(ctx context.Context, query string) error {
	// TODO: This
	return nil
}
