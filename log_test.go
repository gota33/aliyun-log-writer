package sls

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLog(t *testing.T) {
	assert.Equal(t, io.Discard, logger.Writer())

	SetDebug(true)
	assert.Equal(t, os.Stdout, logger.Writer())

	SetDebug(false)
	assert.Equal(t, io.Discard, logger.Writer())
}
