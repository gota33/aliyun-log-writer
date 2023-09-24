package sls

import (
	"errors"
	"strconv"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWorker(t *testing.T) {
	SetDebug(true)
	t.Run("normal", func(t *testing.T) {
		const (
			bufferSize  = 10
			deliverSize = int(3.5 * bufferSize)
			batchSize   = 2 + 3 + 1
			interval    = 20 * time.Millisecond
		)

		client := &MockSender{}

		w := newWorker(workerOption{
			bufferSize: bufferSize,
			interval:   interval,
			onError:    func(err error) { assert.NoError(t, err) },
			client:     client,
		})

		w.Start()

		for i, msg := range makeMessages(deliverSize) {
			err := w.Submit(msg)
			assert.NoError(t, err)

			if i == 1 || i == 3 {
				time.Sleep(interval * 2)
			}
		}

		w.Stop()

		assert.EqualValues(t, deliverSize, client.total)
		assert.GreaterOrEqual(t, client.batch, batchSize)
	})

	t.Run("onError", func(t *testing.T) {
		const total = 10
		count := &atomic.Int64{}

		client := &MockSender{err: errors.New("test error")}

		w := newWorker(workerOption{
			bufferSize: 2,
			interval:   1 * time.Second,
			onError:    func(err error) { count.Add(1) },
			client:     client,
		}).(*asyncWorker)

		w.Start()

		for _, msg := range makeMessages(total) {
			err := w.Submit(msg)
			assert.NoError(t, err)
		}

		w.Stop()
		assert.EqualValues(t, client.batch, count.Load())
	})

	t.Run("write after closed", func(t *testing.T) {
		w := newWorker(workerOption{
			bufferSize: 1,
			interval:   time.Second,
			onError:    nil,
			client:     &MockSender{},
		})
		w.Start()
		w.Stop()

		err := w.Submit(Message{})
		assert.ErrorIs(t, err, ErrClosed)
	})
}

type MockSender struct {
	total int
	batch int
	err   error
}

func (s *MockSender) Send(messages ...Message) error {
	s.total += len(messages)
	s.batch++
	logger.Printf("Total: %d", s.total)
	return s.err
}

func makeMessages(num int) (msgs []Message) {
	msgs = make([]Message, num)
	for i := 0; i < num; i++ {
		msgs[i] = Message{Contents: map[string]string{"no": strconv.Itoa(i)}}
	}
	return
}
