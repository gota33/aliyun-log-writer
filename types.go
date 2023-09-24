package sls

import (
	"encoding/json"
)

type Secret []byte

func (s Secret) String() string { return "******" }

type AliyunError struct {
	HTTPCode  int32  `json:"-"`
	Code      string `json:"errorCode"`
	Message   string `json:"errorMessage"`
	RequestID string `json:"-"`
}

func (a AliyunError) Error() string {
	if data, err := json.Marshal(a); err != nil {
		return err.Error()
	} else {
		return string(data)
	}
}

type MessageModifier interface {
	Modify(msg Message) Message
}

type MessageFilter interface {
	Filter(msg Message) bool
}

type sender interface {
	Send(messages ...Message) error
}

type worker interface {
	Start()
	Stop()
	Submit(msg Message) error
}
