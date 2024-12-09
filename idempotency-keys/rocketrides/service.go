package rocketrides

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ek-os/idempotency-keys/db"
)

func NewHandler(db *db.DB) *Handler {
	return &Handler{
		db: db,
	}
}

type Handler struct {
	db *db.DB
}

func (h *Handler) RegisterRide(w http.ResponseWriter, r *http.Request) {
	var (
		ctx            = r.Context()
		idempotencyKey = r.Header.Get("Idempotency-Key")
	)

	userID, err := strconv.ParseInt(r.Header.Get("User-ID"), 10, 64)
	if err != nil {
		http.Error(w, "User-ID must be a valid integer", http.StatusBadRequest)
		return
	}

	var requestParams json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&requestParams); err != nil {
		http.Error(w, "Body must be valid JSON", http.StatusBadRequest)
		return
	}

	var key *db.Key

	if err := h.atomicPhase(ctx, func(tx *db.Tx) (err error) {
		key, err = tx.FindKey(ctx, userID, idempotencyKey)
		if err != nil {
			return err
		}

		if key != nil {
			// Programs sending multiple requests with different parameters but the
			// same idempotency key is a bug.
			if !bytes.Equal(key.RequestParams, requestParams) {
				return responseError{
					Status:  http.StatusConflict,
					Message: messageMismatch,
				}
			}

			// Only acquire a lock if the key is unlocked or its lock has expired
			// because the original request was long enough ago.
			if key.LockedAt.Valid && key.LockedAt.Value.Add(idempotencyKeyLockTimeoutSeconds).Before(time.Now()) {
				return responseError{
					Status:  http.StatusConflict,
					Message: messageRequestInProgress,
				}
			}

			// Lock the key and update latest run unless the request is already
			// finished.
			if key.RecoveryPoint != recoveryPointFinished {
				if err := tx.LockKey(ctx, key.ID, time.Now()); err != nil {
					return err
				}
			}
		} else {
			key, err = 
		}

		return nil
	}); err != nil {

	}
}

func (h *Handler) atomicPhase(ctx context.Context, fn func(*db.Tx) error) error {
	tx, err := h.db.BeginTx(ctx, &db.TxOptions{
		Isolation: db.LevelSerializable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

type responseError struct {
	Status  int
	Message string
}

func (e responseError) Error() string {
	return fmt.Sprintf("%d: %s", e.Status, e.Message)
}

const (
	messageMismatch          = "There was a mismatch between this request's parameters and the parameters of a previously stored request with the same Idempotency-Key."
	messageRequestInProgress = "An API request with the same Idempotency-Key is already in progress."
)

const idempotencyKeyLockTimeoutSeconds = 90 * time.Second

const (
	recoveryPointFinished = "finished"
)
