package application

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun(t *testing.T) {
	w := &testWriter{messages: []string{}}
	b := Background{
		logger: Logger{out: w, clock: &testClock{}},
	}

	b.Run(func() {
		panic("panic in background task")
	})

	b.Wait()
	assert.Equal(t, "ERROR: [2024-02-07T12:00:00Z] panic in background task \n", w.messages[0])
}
