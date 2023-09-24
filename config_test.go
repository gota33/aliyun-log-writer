package sls

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	raw := Config{
		Endpoint:     "http://test-project.regionid.example.com/logstores/test-logstore",
		AccessKey:    "123",
		AccessSecret: "321",
		Project:      "test-project",
		Store:        "test-store",
		Topic:        "test-topic",
		Source:       "127.0.0.1",
		BufferSize:   10,
		Timeout:      1 * time.Second,
		Interval:     1 * time.Second,
	}

	t.Run("required", func(t *testing.T) {
		c := raw
		c.Endpoint = " "
		assert.Error(t, c.validate())

		c = raw
		c.AccessKey = " "
		assert.Error(t, c.validate())

		c = raw
		c.AccessSecret = " "
		assert.Error(t, c.validate())

		c = raw
		c.Project = " "
		assert.Error(t, c.validate())

		c = raw
		c.Store = " "
		assert.Error(t, c.validate())

		c = raw
		c.Store = " "
		assert.Error(t, c.validate())

		c = raw
		c.Topic = " "
		assert.Error(t, c.validate())
	})

	t.Run("default", func(t *testing.T) {
		c := raw
		c.Source = " "
		if assert.NoError(t, c.validate()) {
			assert.NotEmpty(t, c.Source)
		}

		c = raw
		c.BufferSize = 0
		if assert.NoError(t, c.validate()) {
			assert.Equal(t, DefaultBufferSize, c.BufferSize)
		}

		c = raw
		c.Timeout = 0
		if assert.NoError(t, c.validate()) {
			assert.Equal(t, DefaultTimeout, c.Timeout)
		}

		c = raw
		c.Interval = 0
		if assert.NoError(t, c.validate()) {
			assert.Equal(t, DefaultInterval, c.Interval)
		}

		c = raw
		c.HttpClient = nil
		if assert.NoError(t, c.validate()) {
			assert.Equal(t, http.DefaultClient, c.HttpClient)
		}
	})
}
