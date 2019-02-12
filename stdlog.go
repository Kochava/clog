package clog

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// StdLogger is a logger that implements most methods used by "log" stdlib users. Print-ed logs are
// recorded at Info level.
type StdLogger struct{ log *zap.Logger }

// NewStdLogger creats a new StdLogger attached to the given zap.Logger. If logger is nil,
// NewStdLogger returns nil as well.
func NewStdLogger(logger *zap.Logger) *StdLogger {
	if logger == nil {
		return nil
	}
	return &StdLogger{logger.WithOptions(zap.AddCallerSkip(1))}
}

// With returns a copy of this StdLogger with additional Zap fields.
func (l *StdLogger) With(fields ...zapcore.Field) *StdLogger {
	return &StdLogger{l.log.With(fields...)}
}

// WithOptions returns a copy of this StdLogger with additional Zap options.
func (l *StdLogger) WithOptions(opts ...zap.Option) *StdLogger {
	return &StdLogger{l.log.WithOptions(opts...)}
}

// Panic writes a panic-level fmt.Sprint-formatted log message.
func (l *StdLogger) Panic(args ...interface{}) {
	l.log.Panic(fmt.Sprint(args...))
}

// Panicln writes a panic-level fmt.Sprintln-formatted log message.
func (l *StdLogger) Panicln(args ...interface{}) {
	l.log.Panic(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Panicf writes a panic-level fmt.Sprintf-formatted log message.
func (l *StdLogger) Panicf(format string, args ...interface{}) {
	l.log.Panic(fmt.Sprintf(format, args...))
}

// Fatal writes a fatal-level fmt.Sprint-formatted log message.
func (l *StdLogger) Fatal(args ...interface{}) {
	l.log.Fatal(fmt.Sprint(args...))
}

// Fatalln writes a fatal-level fmt.Sprintln-formatted log message.
func (l *StdLogger) Fatalln(args ...interface{}) {
	l.log.Fatal(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Fatalf writes a fatal-level fmt.Sprintf-formatted log message.
func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, args...))
}

// Print writes an info-level fmt.Sprint-formatted log message.
func (l *StdLogger) Print(args ...interface{}) {
	l.log.Info(fmt.Sprint(args...))
}

// Println writes an info-level fmt.Sprintln-formatted log message.
func (l *StdLogger) Println(args ...interface{}) {
	l.log.Info(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

// Printf writes an info-level fmt.Sprintf-formatted log message.
func (l *StdLogger) Printf(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}
