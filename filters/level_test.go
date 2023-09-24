package filters

import (
	"testing"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/stretchr/testify/assert"
)

func TestLevel(t *testing.T) {
	filter := InfoLevel()
	var msg sls.Message

	msg = sls.Message{Contents: map[string]string{"level": "DEBUG"}}
	assert.False(t, filter.Filter(msg))

	msg = sls.Message{Contents: map[string]string{"level": "INFO"}}
	assert.True(t, filter.Filter(msg))

	msg = sls.Message{Contents: map[string]string{"level": "WARN"}}
	assert.True(t, filter.Filter(msg))

	msg = sls.Message{Contents: make(map[string]string)}
	assert.True(t, filter.Filter(msg))
}
