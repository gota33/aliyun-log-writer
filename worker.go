package sls

import (
	"errors"
	"sync"
	"time"

	"github.com/gota33/aliyun-log-writer/internal/validator"
)

var (
	ErrClosed = errors.New("write to closed writer")
)

type ErrorListener func(err error)

type workerOption struct {
	bufferSize int
	interval   time.Duration
	onError    ErrorListener
	client     sender
}

type asyncWorker struct {
	chData    chan Message
	startFunc *sync.Once
	closeFunc *sync.Once
	wgRunning *sync.WaitGroup
	chQuit    chan struct{}
	chSubmit  chan struct{}
	workerOption
}

func newWorker(opt workerOption) worker {
	w := &asyncWorker{
		workerOption: opt,
		startFunc:    &sync.Once{},
		closeFunc:    &sync.Once{},
		wgRunning:    &sync.WaitGroup{},
		chSubmit:     make(chan struct{}),
		chQuit:       make(chan struct{}),
	}

	w.bufferSize = validator.Coalesce(opt.bufferSize, DefaultBufferSize)
	w.interval = validator.Coalesce(opt.interval, DefaultInterval)
	w.chData = make(chan Message, 2*w.bufferSize)
	return w
}

func (w *asyncWorker) Start() {
	w.startFunc.Do(func() {
		w.wgRunning.Add(1)
		go func() {
			defer w.wgRunning.Done()
			w.run()
		}()
	})
}

func (w *asyncWorker) run() {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for running := true; running; {
		select {
		case <-w.chQuit:
			running = false
		case <-w.chSubmit:
			w.flush(w.chData, w.bufferSize)
		case <-ticker.C:
			w.flush(w.chData, 1)
		}
	}

	ch := w.chData
	w.chData = make(chan Message)

	logger.Printf("Remain: %d", len(ch))
	w.flush(ch, 1)
}

func (w *asyncWorker) Submit(msg Message) (err error) {
	select {
	case <-w.chQuit:
		return ErrClosed
	case w.chData <- msg:
		select {
		case w.chSubmit <- struct{}{}:
		default:
		}
		logger.Printf("Submit: %v", msg)
	}
	return
}

func (w *asyncWorker) Stop() {
	w.closeFunc.Do(func() {
		close(w.chQuit)
	})
	w.wgRunning.Wait()
}

func (w *asyncWorker) flush(chData <-chan Message, atLeast int) {
	for {
		messages := w.takeAtLeast(chData, atLeast)
		if size := len(messages); size == 0 {
			return
		}

		if err := w.client.Send(messages...); err != nil && w.onError != nil {
			w.onError(err)
		}

		if debugMode {
			logger.Printf("Flush %d messages", len(messages))
			for i, message := range messages {
				logger.Printf("Flush[%d]: %v", i, message)
			}
		}
	}
}

func (w *asyncWorker) takeAtLeast(chData <-chan Message, atLeast int) (msgs []Message) {
	size := len(chData)
	if size < atLeast {
		return
	}

	if size > w.bufferSize {
		size = w.bufferSize
	}
loop:
	for i := 0; i < size; i++ {
		select {
		case msg := <-chData:
			logger.Printf("Pull: %v", msg)
			msgs = append(msgs, msg)
		default:
			break loop
		}
	}
	logger.Printf("Pull %d messages", len(msgs))
	return
}
