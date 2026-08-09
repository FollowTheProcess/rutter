// Package db contains sqlite db interactions via [sqlc].
//
// [sqlc]: https://sqlc.dev
package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pressly/goose/v3"
	"go.followtheprocess.codes/rutter/internal/db/migrations"

	_ "modernc.org/sqlite"
)

// defaultDirPermissions is the file mode permissions used for any directories
// created by this package.
const defaultDirPermissions = 0o700

// Store represents the history data storage backend and manages
// the lifecycle of the backing DB.
type Store struct {
	Query *Queries

	sql *sql.DB
}

// Open connects to the database at path, creating and performing
// migrations if needed.
//
// If the path does not exist, it's parent directories will be created
// as necessary.
//
// Once opened, a Store must be deferred closed with [Store.Close].
func Open(ctx context.Context, path string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(path), defaultDirPermissions); err != nil {
		return nil, fmt.Errorf("failed to create DB parent dir: %w", err)
	}

	pragmas := []string{
		"_pragma=journal_mode(WAL)",
		"_pragma=busy_timeout(5000)",
		"_pragma=foreign_keys(1)",
		"_pragma=synchronous(NORMAL)",
		"_pragma=secure_delete(1)",
		"_inttotime=1",
		"_time_integer_format=unix_nano",
	}

	dsn := fmt.Sprintf("file:%s?%s", path, strings.Join(pragmas, "&"))

	sqlDB, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB at path %s: %w", path, err)
	}

	if err := migrate(ctx, sqlDB); err != nil {
		_ = sqlDB.Close()

		return nil, err
	}

	return &Store{
		Query: New(sqlDB),
		sql:   sqlDB,
	}, nil
}

// Close releases the underlying database resource.
func (s *Store) Close() error {
	return s.sql.Close()
}

// migrate performs the DB migrations.
func migrate(ctx context.Context, sqlDB *sql.DB) error {
	provider, err := goose.NewProvider(goose.DialectSQLite3, sqlDB, migrations.FS)
	if err != nil {
		return fmt.Errorf("failed to build migration provider: %w", err)
	}

	if _, err := provider.Up(ctx); err != nil {
		return fmt.Errorf("failed to perform up migration(s): %w", err)
	}

	return nil
}
