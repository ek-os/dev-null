package lmgsql

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	CHANGELOG_TABLE_NAME = "lmg_changelog"
	LOCK_TABLE_NAME      = "lmg_lock"
)

type DB interface {
	ChangelogTableExists(ctx context.Context) (bool, error)
	CreateChangelogTable(ctx context.Context) error
	LockTableExists(ctx context.Context) (bool, error)
	CreateLockTable(ctx context.Context) error
	Exec(ctx context.Context, query string) error
}

func Open(driver, dsn string) (DB, error) {
	switch driver {
	case "sqlite3":
		db, err := sql.Open(driver, dsn)
		if err != nil {
			return nil, fmt.Errorf("sql.Open: %w", err)
		}
		return &sqlite3DB{db: db}, nil
	default:
		return nil, &ErrUnknownDriver{Driver: driver}
	}
}

type ErrUnknownDriver struct {
	Driver string
}

func (e *ErrUnknownDriver) Error() string {
	return fmt.Sprintf("unknown driver: %s", e.Driver)
}
