package modifiers

import (
	"log/slog"

	sls "github.com/gota33/aliyun-log-writer"
)

type RenameField struct {
	Before string
	After  string
}

func (m *RenameField) Modify(msg sls.Message) sls.Message {
	if value, ok := msg.Contents[m.Before]; ok {
		msg.Contents[m.After] = value
		delete(msg.Contents, m.Before)
	}
	return msg
}

func RenameMessageField() sls.MessageModifier {
	return &RenameField{
		Before: slog.MessageKey,
		After:  "message",
	}
}
