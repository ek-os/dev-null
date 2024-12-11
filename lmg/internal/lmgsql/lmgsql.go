package lmgsql

import (
	"context"
	"database/sql"
	"fmt"
)

const CHANGELOG_TABLE_NAME = "lmg_changelog"

type ErrUnknownDriver struct {
	Driver string
}

func (e *ErrUnknownDriver) Error() string {
	return fmt.Sprintf("Unknown driver: %s", e.Driver)
}

func Open(driver string, dsn string) (DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	switch driver {
	case "sqlite3":
		return &sqlite3DB{db: db}, nil
	default:
		return nil, &ErrUnknownDriver{Driver: driver}
	}
}

type DB interface {
	ChangelogExists(ctx context.Context) (bool, error)
	CreateChangelogTable(ctx context.Context) error
	MigrateContext(ctx context.Context, query string) error
}
