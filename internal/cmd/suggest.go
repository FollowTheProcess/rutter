package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/rutter/internal/app"
	"go.followtheprocess.codes/rutter/internal/xdg"
)

func suggest() (*cli.Command, error) {
	var query string

	return cli.New(
		"suggest",
		cli.Short("Suggest a command based on a fuzzy match, used for shell autosuggest"),
		cli.Arg(&query, "query", "Search term"),
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

			return app.Suggest(ctx, query)
		}),
	)
}
