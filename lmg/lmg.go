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

	db, err := lmgsql.Open(driver, dsn)
	if err != nil {
		return fmt.Errorf("lmgsql.Open: %w", err)
	}

	for _, migration := range migrations {
		if err := executeMigration(ctx, db, migration); err != nil {
			return fmt.Errorf("execute %s: %w", migration, err)
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

	if err := db.MigrateContext(ctx, string(query)); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	return nil
}
