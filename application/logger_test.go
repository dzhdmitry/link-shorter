package application

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type testWriter struct {
	messages []string
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	w.messages = append(w.messages, string(p))

	return len(p), nil
}

type testClock struct {
	//
}

func (t *testClock) Now() time.Time {
	location, _ := time.LoadLocation("Europe/London")

	return time.Date(2024, 2, 7, 12, 0, 0, 0, location)
}

func TestLogInfo(t *testing.T) {
	w := &testWriter{}
	l := NewLogger(w, &testClock{})

	l.LogInfo("test info")
	assert.Equal(t, "INFO: [2024-02-07T12:00:00Z] test info \n", w.messages[0])
}

func TestLogError(t *testing.T) {
	w := testWriter{}
	l := NewLogger(&w, &testClock{})

	l.LogError(errors.New("test error"))

	assert.Equal(t, "ERROR: [2024-02-07T12:00:00Z] test error \n", w.messages[0])
}
