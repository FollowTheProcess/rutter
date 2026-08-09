// Package migrations provides the SQL migrations for the history DB.
package migrations

import "embed"

// FS holds the SQL migrations using [goose], embedded
// so the program can upgrade the DB without any files
// on disk.
//
// [goose]: https://pressly.github.io/goose
//
//go:embed *.sql
var FS embed.FS
