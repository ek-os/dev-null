package dbs

import (
	"context"
	"database/sql"
)

// Queries exposes all app queries via its receiver methods
type Queries struct {
	stmts stmts
}

type stmts interface {
	stmt(ctx context.Context, query string) (*sql.Stmt, error)
}

type User struct {
	ID        int64
	FirstName string
	LastName  string
}

func (q *Queries) FindUser(ctx context.Context, id int64) (*User, error) {
	stmt, err := q.stmts.stmt(ctx, "SELECT first_name, last_name FROM users WHERE id = @id")
	if err != nil {
		return nil, err
	}

	u := &User{ID: id}
	if err := stmt.QueryRowContext(ctx, sql.Named("id", id)).
		Scan(&u.FirstName, &u.LastName); err != nil {
		return nil, err
	}

	return u, nil
}

func (q *Queries) SaveUser(ctx context.Context, firstName, lastName string) (int64, error) {
	stmt, err := q.stmts.stmt(ctx, "INSERT INTO users (first_name, last_name) VALUES (@first_name, @last_name)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, sql.Named("first_name", firstName), sql.Named("last_name", lastName))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
