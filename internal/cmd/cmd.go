// Package cmd implements the rutter CLI.
package cmd

import (
	"go.followtheprocess.codes/cli"
)

//nolint:gochecknoglobals // These have to be here for ldflags to set them.
var (
	version = "dev"
	commit  = ""
	date    = ""
)

// Build returns the root rutter CLI command.
func Build() (*cli.Command, error) {
	return cli.New(
		"rutter",
		cli.Short("Sail through your shell history ⚓"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(date),
		cli.SubCommands(uuidCmd, add, finish),
	)
}
