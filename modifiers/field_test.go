package modifiers

import (
	"testing"
	"time"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	m := RenameMessageField()

	output := m.Modify(sls.Message{
		Time:     time.Now(),
		Contents: map[string]string{"msg": "1"},
	})

	_, ok := output.Contents["msg"]
	assert.False(t, ok)

	value := output.Contents["message"]
	assert.Equal(t, value, "1")
}
