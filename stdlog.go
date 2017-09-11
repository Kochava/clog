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

func (l *StdLogger) With(fields ...zapcore.Field) *StdLogger {
	return &StdLogger{l.log.With(fields...)}
}

func (l *StdLogger) WithOptions(opts ...zap.Option) *StdLogger {
	return &StdLogger{l.log.WithOptions(opts...)}
}

func (l *StdLogger) Panic(args ...interface{}) {
	l.log.Panic(fmt.Sprint(args...))
}

func (l *StdLogger) Panicln(args ...interface{}) {
	l.log.Panic(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

func (l *StdLogger) Panicf(format string, args ...interface{}) {
	l.log.Panic(fmt.Sprintf(format, args...))
}

func (l *StdLogger) Fatal(args ...interface{}) {
	l.log.Fatal(fmt.Sprint(args...))
}

func (l *StdLogger) Fatalln(args ...interface{}) {
	l.log.Fatal(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

func (l *StdLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal(fmt.Sprintf(format, args...))
}

func (l *StdLogger) Print(args ...interface{}) {
	l.log.Info(fmt.Sprint(args...))
}

func (l *StdLogger) Println(args ...interface{}) {
	l.log.Info(strings.TrimSuffix(fmt.Sprintln(args...), "\n"))
}

func (l *StdLogger) Printf(format string, args ...interface{}) {
	l.log.Info(fmt.Sprintf(format, args...))
}
