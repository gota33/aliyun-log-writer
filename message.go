package sls

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"
)

type Message struct {
	Time     time.Time
	Contents map[string]string
}

func (msg *Message) UnmarshalJSON(data []byte) (err error) {
	var m map[string]json.RawMessage
	if err = json.Unmarshal(data, &m); err != nil {
		return
	}

	*msg = Message{
		Time:     time.Now(),
		Contents: make(map[string]string, len(m)),
	}

	if raw, ok := m[slog.TimeKey]; ok {
		if err = json.Unmarshal(raw, &msg.Time); err != nil {
			return
		}
		delete(m, slog.TimeKey)
	}

	for k, v := range m {
		var value string
		if value, err = formatJsonValue(v); err != nil {
			return
		}
		msg.Contents[k] = value
	}
	return
}

func formatJsonValue(input json.RawMessage) (output string, err error) {
	// bool, for JSON booleans
	// float64, for JSON numbers
	// string, for JSON strings
	// []interface{}, for JSON arrays
	// map[string]interface{}, for JSON objects
	// nil for JSON null

	switch input[0] {
	case '"':
		output = string(bytes.Trim(input, `"`))
	case '[', '{':
		var data []byte
		if data, err = json.Marshal(input); err == nil {
			output = string(data)
		}
	default:
		output = fmt.Sprintf(`%s`, input)
	}
	return
}
