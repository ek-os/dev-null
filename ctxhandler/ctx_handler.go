package ctxhandler

import (
	"context"
	"log/slog"
)

func New(h slog.Handler) slog.Handler {
	return CtxHandler{h: h}
}

type CtxHandler struct {
	h slog.Handler
}

// Enabled implements slog.Handler.
func (c CtxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return c.h.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (c CtxHandler) Handle(ctx context.Context, r slog.Record) error {
	attrs, ok := ctx.Value(ctxkey{}).(*[]slog.Attr)
	if ok {
		r.AddAttrs(*attrs...)
	}
	return c.h.Handle(ctx, r)
}

// WithAttrs implements slog.Handler.
func (c CtxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return c.h.WithAttrs(attrs)
}

// WithGroup implements slog.Handler.
func (c CtxHandler) WithGroup(name string) slog.Handler {
	return c.h.WithGroup(name)
}

type ctxkey struct{}

func CtxWithAttr(ctx context.Context, a slog.Attr) context.Context {
	attrs, ok := ctx.Value(ctxkey{}).(*[]slog.Attr)
	if ok {
		*attrs = append(*attrs, a)
		return ctx
	}

	newAttrs := make([]slog.Attr, 0, 1)
	newAttrs = append(newAttrs, a)
	return context.WithValue(ctx, ctxkey{}, &newAttrs)
}
