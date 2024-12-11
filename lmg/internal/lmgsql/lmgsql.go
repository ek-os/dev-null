package lmgsql

import (
	"context"
	"database/sql"
	"fmt"
)

const CHANGELOG_TABLE_NAME = "lmg_changelog"

type DB interface {
	ChangelogExists(ctx context.Context) (bool, error)
	CreateChangelogTable(ctx context.Context) error
	MigrateContext(ctx context.Context, query string) error
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
	return fmt.Sprintf("Unknown driver: %s", e.Driver)
}
