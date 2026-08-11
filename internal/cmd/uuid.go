package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"go.followtheprocess.codes/cli"
)

// uuidCmd builds the uuid subcommand.
//
// The uuid command's job is simply to generate a RUTTER_SESSION_ID
// from the shell init script.
func uuidCmd() (*cli.Command, error) {
	return cli.New(
		"uuid",
		cli.Short("Generate a rutter session id"),
		cli.Run(func(ctx context.Context, cmd *cli.Command) error {
			id, err := uuid.NewRandom()
			if err != nil {
				return fmt.Errorf("failed to generate uuid: %w", err)
			}

			_, err = io.WriteString(cmd.Stdout(), id.String()+"\n")

			return err
		}),
	)
}
