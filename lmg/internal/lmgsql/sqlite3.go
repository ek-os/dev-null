package lmgsql

import (
	"context"
	"database/sql"
)

var _ DB = (*sqlite3DB)(nil)

type sqlite3DB struct {
	db *sql.DB
}

// ExecuteContext implements DB.
func (s *sqlite3DB) ExecContext(ctx context.Context, query string) error {
	_, err := s.db.ExecContext(ctx, query)
	return err
}
