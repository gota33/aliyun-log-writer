package handlers

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiHandler(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		h0 := &MockHandler{enabled: true}
		h1 := &MockHandler{enabled: true}
		h := MultiHandler(h0, h1)

		ctx := context.Background()
		assert.True(t, h.Enabled(ctx, slog.LevelDebug))

		assert.NoError(t, h.Handle(ctx, slog.Record{}))
		assert.Equal(t, 1, h0.handleCount)
		assert.Equal(t, 1, h1.handleCount)

		h.WithAttrs(nil)
		assert.Equal(t, 1, h0.withAttrCount)
		assert.Equal(t, 1, h1.withAttrCount)

		h.WithGroup("")
		assert.Equal(t, 1, h0.withGroupCount)
		assert.Equal(t, 1, h1.withGroupCount)
	})

	t.Run("not enable", func(t *testing.T) {
		h0 := &MockHandler{enabled: false}
		h1 := &MockHandler{enabled: false}
		h := MultiHandler(h0, h1)

		ctx := context.Background()
		assert.False(t, h.Enabled(ctx, slog.LevelDebug))
		assert.NoError(t, h.Handle(ctx, slog.Record{}))
		assert.Equal(t, 0, h0.handleCount)
		assert.Equal(t, 0, h1.handleCount)
	})

	t.Run("error", func(t *testing.T) {
		err := errors.New("test error")
		h0 := &MockHandler{enabled: true, err: nil}
		h1 := &MockHandler{enabled: true, err: err}
		h := MultiHandler(h0, h1)

		ctx := context.Background()
		h.Enabled(ctx, slog.LevelDebug)
		assert.ErrorIs(t, h.Handle(ctx, slog.Record{}), err)
	})
}

type MockHandler struct {
	enabled        bool
	handleCount    int
	withAttrCount  int
	withGroupCount int
	err            error
}

func (h *MockHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.enabled
}

func (h *MockHandler) Handle(ctx context.Context, record slog.Record) error {
	h.handleCount++
	return h.err
}

func (h *MockHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.withAttrCount++
	return h
}

func (h *MockHandler) WithGroup(name string) slog.Handler {
	h.withGroupCount++
	return h
}
