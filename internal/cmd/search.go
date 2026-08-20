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

func search() (*cli.Command, error) {
	var (
		query   string
		options app.SearchOptions
	)

	return cli.New(
		"search",
		cli.Short("Search through the shell history"),
		cli.Arg(&query, "query", "History search term"),
		cli.Flag(&options.Directory, "directory", 'd', "Filter results by directory"),
		cli.Flag(&options.Session, "session", 's', "Filter results by session"),
		cli.Flag(&options.Limit, "limit", 'l', "Maximum number of results to return", cli.FlagDefault(20)),
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

			return app.Search(ctx, query, options)
		}),
	)
}
