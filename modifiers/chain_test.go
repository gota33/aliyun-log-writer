package modifiers

import (
	"fmt"
	"testing"
	"time"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	chain := make(Chain, 0)

	for i := 0; i < 10; i++ {
		chain = append(chain, &MockModifier{
			k: fmt.Sprintf("%d", i),
			v: fmt.Sprintf("%d", i),
		})
	}
	out := chain.Modify(sls.Message{
		Time:     time.Now(),
		Contents: make(map[string]string),
	})
	for i := 0; i < 10; i++ {
		k := fmt.Sprintf("%d", i)
		v := fmt.Sprintf("%d", i)
		assert.Equal(t, v, out.Contents[k])
	}
}

type MockModifier struct {
	k string
	v string
}

func (m *MockModifier) Modify(msg sls.Message) sls.Message {
	msg.Contents[m.k] = m.v
	return msg
}
