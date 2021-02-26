package clog

import (
	"context"

	"go.uber.org/zap/zapcore"
)

// FieldGenerator describes a function that generates zap fields for a given context
type FieldGenerator func(ctx context.Context) []zapcore.Field

// FieldGenerators retrieves the FieldGenerators stored on a context
func FieldGenerators(ctx context.Context) (fieldGenerators []FieldGenerator) {
	if ctx != nil {
		fieldGenerators, _ = ctx.Value(ctxFields).([]FieldGenerator)
	} else {
		fieldGenerators = []FieldGenerator{}
	}

	return fieldGenerators
}

// WithFieldGenerators adds the supplied FieldGenerator to those found on the supplied context and surfaces a new context
func WithFieldGenerators(ctx context.Context, generators ...FieldGenerator) context.Context {
	if len(generators) == 0 {
		return ctx
	}

	return context.WithValue(ctx, ctxFields, append(FieldGenerators(ctx), generators...))
}

func generateFields(ctx context.Context) []zapcore.Field {
	var fields []zapcore.Field

	for _, gen := range FieldGenerators(ctx) {
		fields = append(fields, gen(ctx)...)
	}

	return fields
}
