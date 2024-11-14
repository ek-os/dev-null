package dbs

import (
	"context"
	"database/sql"
)

// Queries exposes all app queries via its receiver methods
type Queries struct {
	dbi dbi
}

// dbi is the common database access interface shared by sql.DB and sql.Tx
type dbi interface {
	Exec(query string, args ...any) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func (q *Queries) FindUser(ctx context.Context, id int64) (*User, error) {
	u := &User{ID: id}
	if err := q.dbi.QueryRowContext(
		ctx,
		"SELECT first_name, last_name FROM users WHERE id = @id",
		sql.Named("id", id),
	).Scan(&u.FirstName, &u.LastName); err != nil {
		return nil, err
	}
	return u, nil
}

func (q *Queries) SaveUser(ctx context.Context, firstName, lastName string) (int64, error) {
	res, err := q.dbi.ExecContext(
		ctx,
		"INSERT INTO users (first_name, last_name) VALUES (@first_name, @last_name)",
		sql.Named("first_name", firstName),
		sql.Named("last_name", lastName),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
