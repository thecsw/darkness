package puck

import (
	"os"
	"time"

	l "github.com/charmbracelet/log"
	"github.com/thecsw/darkness/v3/yunyun"
)

const (
	// DebugLevel is the debug level.
	DebugLevel l.Level = iota - 1
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
func (s stopwatch) Record(loggers ...*l.Logger) time.Duration {
	logger := Logger
	if len(loggers) > 0 {
		logger = loggers[0]
	}
	elapsed := time.Since(s.start)
	logger.Info(s.msg, append(s.msgs, "elapsed", elapsed)...)
	return elapsed
}

// RecordWithFile records the elapsed time since the stopwatch was created and
// calls the given fileTimeRecorder with the given key and elapsed time.
func (s stopwatch) RecordWithFile(
	fileTimeRecorder func(yunyun.FullPathFile, time.Duration),
	key yunyun.FullPathFile,
	loggers ...*l.Logger) time.Duration {
	logger := Logger
	if len(loggers) > 0 {
		logger = loggers[0]
	}
	elapsed := time.Since(s.start)
	fileTimeRecorder(key, elapsed)
	logger.Info(s.msg, append(s.msgs, "elapsed", elapsed)...)
	return elapsed
}

// Stopwatch is a simple stopwatch that can be used to time operations.
func Stopwatch(msg any, msgs ...any) interface {
	Record(...*l.Logger) time.Duration
	RecordWithFile(func(yunyun.FullPathFile, time.Duration), yunyun.FullPathFile, ...*l.Logger) time.Duration
} {
	s := stopwatch{
		start: time.Now(),
		msg:   msg,
		msgs:  msgs,
	}
	return s
}
