package lmgsql

import (
	"context"
	"database/sql"
	"errors"
)

var _ DB = (*sqlite3DB)(nil)

type sqlite3DB struct {
	db *sql.DB
}

// ChangelogTableExists implements DB.
func (s *sqlite3DB) ChangelogTableExists(ctx context.Context) (bool, error) {
	return s.tableExists(ctx, CHANGELOG_TABLE_NAME)
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

// LockTableExists implements DB.
func (s *sqlite3DB) LockTableExists(ctx context.Context) (bool, error) {
	return s.tableExists(ctx, LOCK_TABLE_NAME)
}

func (s *sqlite3DB) tableExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	if err := s.db.QueryRowContext(
		ctx,
		"SELECT true FROM sqlite_master WHERE name = :name",
		sql.Named("name", name),
	).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

// CreateLockTable implements DB.
func (s *sqlite3DB) CreateLockTable(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS lmg_lock (
			id 		INT NOT NULL PRIMARY KEY,
			locked 		BOOLEAN NOT NULL,
			locked_at 	TIMESTAMP,
			locked_by 	TEXT
		);
	`)
	return err
}

// Exec implements DB.
func (s *sqlite3DB) Exec(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	return err
}
