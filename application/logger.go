package application

import (
	"fmt"
	"io"
	"time"
)

type Logger struct {
	out io.Writer
}

func NewLogger(out io.Writer) *Logger {
	return &Logger{out: out}
}

func (l Logger) LogInfo(info string) {
	l.print("INFO", info)
}

func (l Logger) LogError(err error) {
	l.print("ERROR", err.Error())
}

func (l Logger) print(level, message string) {
	record := fmt.Sprintf(
		"%s: [%s] %s \n",
		level,
		time.Now().UTC().Format(time.RFC3339),
		message,
	)

	l.out.Write([]byte(record))
}
