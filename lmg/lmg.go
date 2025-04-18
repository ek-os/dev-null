package lmg

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ek-os/lmg/internal/lmgsql"
)

const (
	ENV_CHANGELOG = "LMG_CHANGELOG_PATH"
	ENV_DRIVER    = "LMG_DRIVER"
	ENV_DSN       = "LMG_DSN"
)

func Run() {
	if err := run(context.Background(), realSystem{}); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

type system interface {
	Getenv(key string) string
	Stdout() io.Writer
}

type realSystem struct{}

func (realSystem) Getenv(key string) string {
	return os.Getenv(key)
}

func (realSystem) Stdout() io.Writer {
	return os.Stdout
}

func run(ctx context.Context, sys system) error {
	var (
		changelogPath = sys.Getenv(ENV_CHANGELOG)
		driver        = sys.Getenv(ENV_DRIVER)
		dsn           = sys.Getenv(ENV_DSN)
	)

	db, err := lmgsql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("lmgsql.Open: %w", err)
	}

	if err := ensureLockTableExists(ctx, db); err != nil {
		return err
	}

	migrations, err := readChangelog(changelogPath)
	if err != nil {
		return fmt.Errorf("read changelog: %w", err)
	}

	for _, migration := range migrations {
		if err := executeMigration(ctx, db, migration); err != nil {
			return fmt.Errorf("execute %s: %w", migration, err)
		}
	}

	return nil
}

func ensureChangelogTableExists(ctx context.Context, db lmgsql.DB) error {
	ok, err := db.ChangelogTableExists(ctx)
	if err != nil {
		return fmt.Errorf("check if changelog table exists: %w", err)
	}

	if !ok {
		if err := db.CreateChangelogTable(ctx); err != nil {
			return fmt.Errorf("create changelog table: %w", err)
		}
	}

	return nil
}

func ensureLockTableExists(ctx context.Context, db lmgsql.DB) error {
	ok, err := db.LockTableExists(ctx)
	if err != nil {
		return fmt.Errorf("check if lock table exists: %w", err)
	}

	if !ok {
		if err := db.CreateLockTable(ctx); err != nil {
			return fmt.Errorf("create lock table: %w", err)
		}
	}

	return nil
}

func readChangelog(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var (
		migrationPaths []string
		s              = bufio.NewScanner(f)
		dir            = filepath.Dir(path)
	)
	for s.Scan() {
		if len(s.Text()) == 0 {
			continue
		}
		migrationPath := filepath.Join(dir, s.Text())
		migrationPaths = append(migrationPaths, migrationPath)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return migrationPaths, nil
}

func executeMigration(ctx context.Context, db lmgsql.DB, path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if err := db.Exec(ctx, string(query)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
