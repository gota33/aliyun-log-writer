package modifiers

import (
	"log/slog"
	"testing"
	"time"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	m := RemapLevelToSysLog()

	mapper := map[slog.Level]string{
		slog.LevelDebug: "7",
		slog.LevelInfo:  "6",
		slog.LevelWarn:  "4",
		slog.LevelError: "3",
	}

	for level, value := range mapper {
		msg := m.Modify(sls.Message{
			Time:     time.Time{},
			Contents: map[string]string{"level": level.String()},
		})
		assert.Equal(t, msg.Contents["level"], value)
	}

	msg := m.Modify(sls.Message{
		Time:     time.Time{},
		Contents: map[string]string{"level": "unknown"},
	})
	assert.Equal(t, msg.Contents["level"], "unknown")
}
