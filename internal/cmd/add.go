package cmd

import (
	"context"
	"fmt"

	"go.followtheprocess.codes/cli"
)

// add returns the add subcommand.
//
// The add command is responsible for inserting a new command
// in the history. A command added with add is not complete and
// must be finished with the finish command to update it's exit code
// and duration.
//
// The stored ID of the command is printed so the shell hook can
// finalise it with the finish command.
func add() (*cli.Command, error) {
	var command string

	return cli.New(
		"add",
		cli.Short("Adds a new command to the history, returning it's ID"),
		cli.Arg(&command, "command", "The command to add"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			fmt.Fprintf(cmd.Stdout(), "Adding %s\n", command)

			return nil
		}),
	)
}
