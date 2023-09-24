package validator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRequired(t *testing.T) {
	field, value := "f1", "v1"
	assert.NoError(t, Required(field, value))
	assert.Error(t, Required(field, ""))
}

func TestIllegalArgument(t *testing.T) {
	field, value := "f1", "v1"
	assert.Error(t, IllegalArgument(field, value))
}

func TestCoalesceStr(t *testing.T) {
	s0, s1 := "s0", "s1"
	assert.Equal(t, s0, Coalesce(s0, s1))
	assert.Equal(t, s1, Coalesce("", s1))
	assert.Equal(t, s1, Coalesce("\t", s1))
}

func TestCoalesceInt(t *testing.T) {
	n1, n2 := 1, 2
	assert.Equal(t, n1, Coalesce(n1, n2))
	assert.Equal(t, n2, Coalesce(0, n2))
	assert.Equal(t, n2, Coalesce(-1, n2))
}

func TestCoalesceDur(t *testing.T) {
	t1, t2 := time.Second, 2*time.Second
	assert.Equal(t, t1, Coalesce(t1, t2))
	assert.Equal(t, t2, Coalesce(0, t2))
	assert.Equal(t, t2, Coalesce(-1*time.Second, t2))
}
