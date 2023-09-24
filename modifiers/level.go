package modifiers

import (
	"log/slog"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/gota33/aliyun-log-writer/internal/validator"
)

type RemapLevel struct {
	LevelKey string
	Mapper   func(level string) string
}

func (m *RemapLevel) Modify(msg sls.Message) sls.Message {
	key := validator.Coalesce(m.LevelKey, slog.LevelKey)
	if level, ok := msg.Contents[key]; ok {
		msg.Contents[key] = m.Mapper(level)
	}
	return msg
}

func RemapLevelToSysLog() sls.MessageModifier {
	return &RemapLevel{
		LevelKey: slog.LevelKey,
		Mapper:   mapToSysLog,
	}
}

func mapToSysLog(level string) string {
	switch level {
	case "DEBUG":
		return "7"
	case "INFO":
		return "6"
	case "WARN":
		return "4"
	case "ERROR":
		return "3"
	default:
		return level
	}
}
