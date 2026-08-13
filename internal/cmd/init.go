package cmd

import (
	"context"
	"fmt"

	"go.followtheprocess.codes/cli"
)

// initCmd returns the init subcommand.
//
// The init subcommand prints the shell startup script to be sourced
// in e.g. .zshrc.
func initCmd() (*cli.Command, error) {
	var shell string

	return cli.New(
		"init",
		cli.Short("Print rutter's shell init script"),
		cli.Arg(&shell, "shell", "Choice of shell, one of 'zsh', 'fish' or 'nu'"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			// Not hooked up to app as it's a static script that
			// doesn't need the DB
			fmt.Fprintf(cmd.Stdout(), "A script here...\n")

			return nil
		}),
	)
}
