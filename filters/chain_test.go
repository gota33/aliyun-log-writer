package filters

import (
	"testing"

	sls "github.com/gota33/aliyun-log-writer"
	"github.com/stretchr/testify/assert"
)

func TestChain(t *testing.T) {
	chain := Chain{&MockFilter{value: true}}
	assert.True(t, chain.Filter(sls.Message{}))

	chain = Chain{
		&MockFilter{value: true},
		&MockFilter{value: false},
	}
	assert.False(t, chain.Filter(sls.Message{}))
}

type MockFilter struct {
	value bool
}

func (m *MockFilter) Filter(msg sls.Message) bool {
	return m.value
}
