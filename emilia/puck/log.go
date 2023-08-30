package puck

import (
	"os"
	"time"

	l "github.com/charmbracelet/log"
)

// Logger is darkness' logger.
var Logger = NewLogger("Darkness 🥬 ")

// NewLogger returns a new logger with the given prefix.
func NewLogger(prefix string) *l.Logger {
	return l.NewWithOptions(os.Stderr, l.Options{
		Prefix:          prefix,
		TimeFormat:      time.DateTime,
		ReportTimestamp: true,
		ReportCaller:    false,
		Level:           l.WarnLevel,
	})
}

// Stopwatch is a stopwatch.
type stopwatch struct {
	start time.Time
	msg   any
	msgs  []any
}

// Record records the elapsed time since the stopwatch was created.
func (s stopwatch) Record(loggers ...*l.Logger) {
	logger := Logger
	if len(loggers) > 0 {
		logger = loggers[0]
	}
	logger.Info(s.msg, append(s.msgs, "elapsed", time.Since(s.start))...)
}

// Stopwatch is a simple stopwatch that can be used to time operations.
func Stopwatch(msg any, msgs ...any) interface {
	Record(...*l.Logger)
} {
	s := stopwatch{
		start: time.Now(),
		msg:   msg,
		msgs:  msgs,
	}
	return s
}
