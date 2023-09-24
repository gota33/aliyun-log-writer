package modifiers

import sls "github.com/gota33/aliyun-log-writer"

type Chain []sls.MessageModifier

func (m Chain) Modify(msg sls.Message) sls.Message {
	for _, modifier := range m {
		msg = modifier.Modify(msg)
	}
	return msg
}
