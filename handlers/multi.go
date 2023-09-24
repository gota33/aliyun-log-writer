package handlers

import (
	"context"
	"errors"
	"log/slog"
)

func MultiHandler(handlers ...slog.Handler) slog.Handler {
	return &multiHandler{
		handlers: handlers,
		enables:  make([]bool, len(handlers)),
	}
}

type multiHandler struct {
	handlers []slog.Handler
	enables  []bool
}

func (h *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for i, handler := range h.handlers {
		h.enables[i] = handler.Enabled(ctx, level)
	}
	for _, enable := range h.enables {
		if enable {
			return true
		}
	}
	return false
}

func (h *multiHandler) Handle(ctx context.Context, record slog.Record) error {
	var errs []error
	for i, handler := range h.handlers {
		if !h.enables[i] {
			continue
		}
		if err := handler.Handle(ctx, record); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (h *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	arr := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		arr[i] = handler.WithAttrs(attrs)
	}
	return MultiHandler(arr...)
}

func (h *multiHandler) WithGroup(name string) slog.Handler {
	arr := make([]slog.Handler, len(h.handlers))
	for i, handler := range h.handlers {
		arr[i] = handler.WithGroup(name)
	}
	return MultiHandler(arr...)
}
