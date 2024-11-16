package ctxhandler

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

func TestCtxHandler(t *testing.T) {
	h := slog.New(New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx := context.Background()

	h.InfoContext(ctx, "Hey there!")

	ctx = CtxWithAttr(ctx, slog.String("foo", "bar"))

	h.InfoContext(ctx, "Hey there!")
}
