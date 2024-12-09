package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

// Queries exposes all app queries via its receiver methods
type Queries struct {
	stmts stmts
}

type stmts interface {
	stmt(ctx context.Context, query string) (*sql.Stmt, error)
}

type Key struct {
	ID             int64
	CreatedAt      time.Time
	IdempotencyKey string
	LastRunAt      time.Time
	LockedAt       sql.Null[time.Time]
	RequestMethod  string
	RequestParams  json.RawMessage
	RequestPath    string
	ResponseCode   sql.Null[int]
	ResponseBody   sql.Null[json.RawMessage]
	RecoveryPoint  string
	UserID         int64
}

type CreateKeyParams struct {
}

func (q *Queries) CreateKey(ctx context.Context, p CreateKeyParams) (*Key, error) {
	stmt, err := q.stmts.stmt(ctx, `
		INSERT INTO keys 
					(idempotency_key
					,locked_at
					,recovery_point
					,request_method
					,request_params
					,request_path
					,user_id)
		     VALUES 
			 		(@idempotency_key
					,@locked_at
					,@recovery_point
					,@request_method
					,@request_params
					,@request_path
					,@user_id)
		  RETURNING
		  			id,
					idempotency_key,
					locked_at,
					recovery_point,
					request_method,
					request_params,
					request_path,
					user_id

	`)
	if err != nil {
		return nil, err
	}

	stmt.QueryRowContext(
		ctx, 
		sql.Named("idempotency_key", p.IdempotencyKey)
	)
}

func (q *Queries) FindKey(ctx context.Context, userID int64, idempotencyKey string) (*Key, error) {
	stmt, err := q.stmts.stmt(ctx, `
		SELECT id,
			   created_at,
			   last_run_at,
			   locked_at,
			   request_method,
			   request_params,
			   request_path,
			   response_code,
			   response_body,
			   recovery_point,
		  FROM idempotency_keys 
		 WHERE user_id         = @user_id
		   AND idempotency_key = @idempotency_key
	`)
	if err != nil {
		return nil, err
	}

	k := &Key{
		IdempotencyKey: idempotencyKey,
		UserID:         userID,
	}

	if err := stmt.QueryRowContext(
		ctx,
		sql.Named("user_id", userID),
		sql.Named("idempotency_key", idempotencyKey),
	).Scan(
		&k.ID,
		&k.CreatedAt,
		&k.LastRunAt,
		&k.LockedAt,
		&k.RequestMethod,
		&k.RequestParams,
		&k.RequestPath,
		&k.ResponseCode,
		&k.ResponseBody,
		&k.RecoveryPoint,
	); err != nil {
		return nil, err
	}

	return k, nil
}

func (q *Queries) LockKey(ctx context.Context, keyID int64, t time.Time) error {
	stmt, err := q.stmts.stmt(ctx, `
		UPDATE idempotency_keys
		   SET last_run_at = @last_run_at,
		   	   locked_at   = @locked_at
		 WHERE id = @id
	`)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(
		ctx,
		sql.Named("last_run_at", t),
		sql.Named("locked_at", t),
		sql.Named("id", keyID),
	); err != nil {
		return err
	}

	return nil
}
