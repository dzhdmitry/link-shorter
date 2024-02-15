package utils

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"link-shorter.dzhdmitry.net/test"
	"testing"
)

func TestLogInfo(t *testing.T) {
	w := &test.Writer{}
	l := NewLogger(w, &test.Clock{})

	l.LogInfo("test info")
	assert.Equal(t, "INFO: [2024-02-07T12:00:00Z] test info \n", w.Messages[0])
}

func TestLogError(t *testing.T) {
	w := test.Writer{}
	l := NewLogger(&w, &test.Clock{})

	l.LogError(errors.New("test error"))
	assert.Equal(t, "ERROR: [2024-02-07T12:00:00Z] test error \n", w.Messages[0])
}
