package sls

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMessage(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		timeStr := "2023-09-24T07:57:30+08:00"
		time0, _ := time.Parse(time.RFC3339, timeStr)
		raw := fmt.Sprintf(`{"time": "%s", "level": "INFO", "msg": "demo", "num": 1, "bool": true, "null": null, "obj": {"a": "b"}, "arr": [1]}`, timeStr)

		var msg Message
		err := json.Unmarshal([]byte(raw), &msg)

		assert.NoError(t, err)
		assert.Equal(t, time0, msg.Time)
		assert.Equal(t, map[string]string{
			"level": "INFO",
			"msg":   "demo",
			"num":   "1",
			"bool":  "true",
			"null":  "null",
			"obj":   `{"a":"b"}`,
			"arr":   `[1]`,
		}, msg.Contents)
	})
}
