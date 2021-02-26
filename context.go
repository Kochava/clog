package clog

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ctxKey is the internal context key type for clog.
type ctxKey int

const (
	// ctxLogger is the context key for attaching a logger to a context.Context. Although this
	// is generally on the side of being fairly gross, attaching a logger to a context tends to
	// work fairly well in a sealed application (i.e., logging is tightly controller).
	ctxLogger ctxKey = iota
	// ctxFields is the context key for attaching FieldGenerators to a context.Context.
	ctxFields
)

// WithLogger returns a context parented to ctx with the given logger attached as a value.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger, logger)
}

// With returns a context parented to ctx with a logger that has the given fields appended to it.
// The logger of ctx is the one returned by Logger.
func With(ctx context.Context, fields ...zapcore.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	return WithLogger(ctx, Logger(ctx).With(fields...))
}

// Logger returns the zap.Logger for the given context.
func Logger(ctx context.Context) (logger *zap.Logger) {
	if ctx != nil {
		logger, _ = ctx.Value(ctxLogger).(*zap.Logger)
	}
	if logger == nil {
		logger = zap.L()
	}

	logger = logger.With(generateFields(ctx)...)
	return logger
}

var skipLoggerFrame = []zap.Option{zap.AddCallerSkip(1)}

// Check is a convenience function for calling Logger(ctx).Check(lvl, msg).
func Check(ctx context.Context, lvl zapcore.Level, msg string) *zapcore.CheckedEntry {
	return Logger(ctx).WithOptions(skipLoggerFrame...).Check(lvl, msg)
}

// DPanic is a convenience function for calling Logger(ctx).DPanic(msg, fields...).
func DPanic(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).DPanic(msg, fields...)
}

// Debug is a convenience function for calling Logger(ctx).Debug(msg, fields...).
func Debug(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Debug(msg, fields...)
}

// Error is a convenience function for calling Logger(ctx).Error(msg, fields...).
func Error(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Error(msg, fields...)
}

// Fatal is a convenience function for calling Logger(ctx).Fatal(msg, fields...).
func Fatal(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Fatal(msg, fields...)
}

// Info is a convenience function for calling Logger(ctx).Info(msg, fields...).
func Info(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Info(msg, fields...)
}

// Panic is a convenience function for calling Logger(ctx).Panic(msg, fields...).
func Panic(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Panic(msg, fields...)
}

// Warn is a convenience function for calling Logger(ctx).Warn(msg, fields...).
func Warn(ctx context.Context, msg string, fields ...zapcore.Field) {
	Logger(ctx).WithOptions(skipLoggerFrame...).Warn(msg, fields...)
}

// Sugar is a convenience function for calling Logger(ctx).Sugar().
func Sugar(ctx context.Context) *zap.SugaredLogger {
	return Logger(ctx).Sugar()
}

// WithOptions is a convenience function for calling WithLogger(ctx, Logger(ctx).WithOptions(opts...)).
func WithOptions(ctx context.Context, opts ...zap.Option) context.Context {
	return WithLogger(ctx, Logger(ctx).WithOptions(opts...))
}
