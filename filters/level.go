package filters

import (
	"log/slog"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/gota33/aliyun-log-writer/internal/validator"
)

type LevelFilter struct {
	LevelKey string
	MinLevel slog.Leveler
}

func (f *LevelFilter) Filter(msg sls.Message) bool {
	key := validator.Coalesce(f.LevelKey, slog.LevelKey)
	if data, ok := msg.Contents[key]; ok {
		var level slog.Level
		if err := level.UnmarshalText([]byte(data)); err == nil {
			return level >= f.MinLevel.Level()
		}
	}
	return true
}

func InfoLevel() sls.MessageFilter {
	return &LevelFilter{
		LevelKey: slog.LevelKey,
		MinLevel: slog.LevelInfo,
	}
}
