package puck

import (
	"os"
	"time"

	l "github.com/charmbracelet/log"
)

const (
	// DebugLevel is the debug level.
	DebugLevelLevel l.Level = iota - 1
	// InfoLevel is the info level.
	InfoLevel
	// WarnLevel is the warn level.
	WarnLevel
	// ErrorLevel is the error level.
	ErrorLevel
	// FatalLevel is the fatal level.
	FatalLevel
)

// Logger is darkness' logger.
var Logger = NewLogger("Darkness 🥬 ")

// NewLogger returns a new logger with the given prefix.
func NewLogger(prefix string, levels ...l.Level) *l.Logger {
	level := WarnLevel
	if len(levels) > 0 {
		level = levels[0]
	}
	return l.NewWithOptions(os.Stderr, l.Options{
		Prefix:          prefix,
		TimeFormat:      time.DateTime,
		ReportTimestamp: true,
		ReportCaller:    false,
		Level:           level,
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
