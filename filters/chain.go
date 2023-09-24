package filters

import sls "github.com/gota33/aliyun-log-writer"

type Chain []sls.MessageFilter

func (m Chain) Filter(msg sls.Message) bool {
	for _, modifier := range m {
		if ok := modifier.Filter(msg); !ok {
			return false
		}
	}
	return true
}
