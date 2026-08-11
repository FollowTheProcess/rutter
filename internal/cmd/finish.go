package cmd

import (
	"context"
	"time"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/rutter/internal/app"
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
			app := app.New(cmd.Stdin(), cmd.Stdout(), cmd.Stderr())

			return app.Finish(id, exit, duration)
		}),
	)
}
