package lmgsql

import (
	"context"
	"database/sql"
)

var _ DB = (*sqlite3DB)(nil)

type sqlite3DB struct {
	db *sql.DB
}

// ChangelogExists implements DB.
func (s *sqlite3DB) ChangelogExists(ctx context.Context) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(
		ctx,
		"SELECT true FROM sqlite_master WHERE name = :name",
		sql.Named("name", CHANGELOG_TABLE_NAME),
	).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}

// CreateChangelogTable implements DB.
func (s *sqlite3DB) CreateChangelogTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS lmg_changelog (
			filename TEXT NOT NULL,
			executed TEXT NOT NULL,
			order    INTEGER NOT NULL
		);
	`)
	return err
}

// MigrateContext implements DB.
func (s *sqlite3DB) MigrateContext(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	return err
}
