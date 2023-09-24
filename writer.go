package sls

import "encoding/json"

type Writer struct {
	worker   worker
	filter   MessageFilter
	modifier MessageModifier
}

func New(c Config) (writer *Writer, err error) {
	if err = c.validate(); err != nil {
		return
	}

	client := &sls{
		Client:    c.HttpClient,
		AppKey:    c.AccessKey,
		AppSecret: Secret(c.AccessSecret),
		Uri:       c.uri,
		Host:      c.uri.Host,
		Topic:     c.Topic,
		Source:    c.Source,
		Timeout:   c.Timeout,
	}

	option := workerOption{
		bufferSize: c.BufferSize,
		interval:   c.Interval,
		onError:    c.OnError,
		client:     client,
	}

	w := newWorker(option)
	w.Start()

	writer = &Writer{
		worker:   w,
		filter:   c.MessageFilter,
		modifier: c.MessageModifier,
	}
	return
}

func (w Writer) Write(data []byte) (n int, err error) {
	var msg Message
	if err = json.Unmarshal(data, &msg); err != nil {
		return
	}

	if w.filter == nil || !w.filter.Filter(msg) {
		return
	}

	if w.modifier != nil {
		msg = w.modifier.Modify(msg)
	}

	if err = w.worker.Submit(msg); err == nil {
		n = len(data)
	}
	return
}

func (w Writer) Close() (err error) {
	w.worker.Stop()
	return
}
