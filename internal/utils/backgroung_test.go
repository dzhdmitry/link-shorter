package utils

import (
	"github.com/dzhdmitry/link-shorter/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	w := &test.Writer{Messages: []string{}}
	b := Background{
		logger: *NewLogger(w, &test.Clock{}),
	}

	b.Run(func() {
		panic("panic in background task")
	})

	b.Wait()
	assert.Equal(t, "ERROR: [2024-02-07T12:00:00Z] panic in background task \n", w.Messages[0])
}
