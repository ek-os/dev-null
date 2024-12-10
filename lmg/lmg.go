package lmg

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	ENV_CHANGELOG = "LMG_CHANGELOG_PATH"
	ENV_DRIVER    = "LMG_DRIVER"
	ENV_DSN       = "LMG_DSN"
)

func Run() {
	if err := run(
		context.Background(),
		os.Getenv,
		os.Stdout,
	); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(
	ctx context.Context,
	getenv func(string) string,
	stdout io.Writer,
) error {
	var (
		changelogPath = getenv(ENV_CHANGELOG)
		driver        = getenv(ENV_DRIVER)
		dsn           = getenv(ENV_DSN)
	)

	migrations, err := readChangelog(changelogPath)
	if err != nil {
		return fmt.Errorf("read changelog: %w", err)
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("sql open: %w", err)
	}

	dir := filepath.Dir(changelogPath)
	for _, migration := range migrations {
		migrationPath := filepath.Join(dir, migration)
		if err := executeMigration(ctx, db, migrationPath); err != nil {
			return fmt.Errorf("execute %s: %w", migrationPath, err)
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

	var migrationPaths []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		migrationPaths = append(migrationPaths, s.Text())
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return migrationPaths, nil
}

func executeMigration(ctx context.Context, db *sql.DB, path string) error {
	query, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if _, err := db.ExecContext(ctx, string(query)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
