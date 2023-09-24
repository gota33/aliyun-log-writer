package sls

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWriter(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		mw := &MockWorker{}

		w := Writer{
			worker:   mw,
			filter:   &MockFilter{},
			modifier: &MockModifier{},
		}

		w.worker.Start()
		assert.True(t, mw.running)

		n, err := fmt.Fprintln(w, `{}`)
		assert.NoError(t, err)
		assert.Greater(t, n, 0)

		err = w.Close()
		assert.NoError(t, err)

		assert.False(t, mw.running)
	})

	t.Run("filter", func(t *testing.T) {
		mw := &MockWorker{}
		w := Writer{
			worker:   mw,
			filter:   &MockFilter{block: true},
			modifier: &MockModifier{},
		}

		for i := 0; i < 10; i++ {
			_, err := w.Write([]byte(`{}`))
			assert.NoError(t, err)
		}

		assert.Equal(t, 0, mw.count)
	})

	t.Run("modifier", func(t *testing.T) {
		msg := Message{
			Time:     time.Now(),
			Contents: map[string]string{"a": "q"},
		}
		mw := &MockWorker{}
		w := Writer{
			worker:   mw,
			filter:   &MockFilter{},
			modifier: &MockModifier{value: msg},
		}

		_, err := w.Write([]byte(`{}`))
		assert.NoError(t, err)
		assert.Equal(t, msg, mw.lastMessage)
	})

	t.Run("error", func(t *testing.T) {
		mw := &MockWorker{err: errors.New("test error")}
		w := Writer{
			worker:   mw,
			filter:   &MockFilter{},
			modifier: &MockModifier{},
		}
		_, err := w.Write([]byte(`{}`))
		assert.ErrorIs(t, err, mw.err)
	})
}

type MockWorker struct {
	running     bool
	count       int
	err         error
	lastMessage Message
}

func (w *MockWorker) Start()                   { w.running = true }
func (w *MockWorker) Stop()                    { w.running = false }
func (w *MockWorker) Submit(msg Message) error { w.count++; w.lastMessage = msg; return w.err }

type MockModifier struct {
	value Message
}

func (m *MockModifier) Modify(msg Message) Message {
	return m.value
}

type MockFilter struct {
	block bool
}

func (m *MockFilter) Filter(msg Message) bool {
	return !m.block
}
