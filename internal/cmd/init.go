package cmd

import (
	"context"

	"go.followtheprocess.codes/cli"
	"go.followtheprocess.codes/rutter/internal/app"
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
			io := app.IO{
				Stdin:  cmd.Stdin(),
				Stdout: cmd.Stdout(),
				Stderr: cmd.Stderr(),
			}
			app := app.New(io)

			return app.Init(shell)
		}),
	)
}
