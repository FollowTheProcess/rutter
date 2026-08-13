package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/rutter/internal/app"
	"go.followtheprocess.codes/rutter/internal/xdg"
)

// finish builds the finish subcommand.
//
// The finish command finalises a history entry in the DB by
// updating it with it's exit code and duration after it has run.
func finish() (*cli.Command, error) {
	var (
		id       int
		exit     int
		duration time.Duration
	)

	return cli.New(
		"finish",
		cli.Short("Finalises a command in the history, adding exit code and duration"),
		cli.Arg(&id, "id", "The id of the command to update"),
		cli.Arg(&exit, "exit", "The exit code of the command"),
		cli.Arg(&duration, "duration", "The duration of the command's execution"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			io := app.IO{
				Stdin:  cmd.Stdin(),
				Stdout: cmd.Stdout(),
				Stderr: cmd.Stderr(),
			}

			data, err := xdg.DataHome()
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}

			sessionID := os.Getenv("RUTTER_SESSION_ID")
			if sessionID == "" {
				return errors.New("RUTTER_SESSION_ID is not set, ensure you source the shell init script")
			}

			app, err := app.New(ctx, io, cwd, sessionID, app.DBPath(data))
			if err != nil {
				return err
			}
			defer app.Close()

			return app.Finish(ctx, id, exit, duration)
		}),
	)
}
