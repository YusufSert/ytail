package common

import (
	"context"
	"log/slog"
)

type DiscardHandler struct{}

func (n *DiscardHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return false
}

func (n *DiscardHandler) Handle(_ context.Context, _ slog.Record) error {
	return nil
}

func (n *DiscardHandler) WithAttrs(_ []slog.Attr) slog.Handler {
	return n
}

func (n *DiscardHandler) WithGroup(_ string) slog.Handler {
	return n
}
