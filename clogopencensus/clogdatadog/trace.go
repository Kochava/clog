// Package clogdatadog helps clog interact with the Open Census Datadog exporter
package clogdatadog

import (
	"context"
	"encoding/binary"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TraceFieldGenerator generates dd.trace_id and dd.span_id fields.
//
// Adding dd.trace_id and dd.span_id to a json formatted log message allows DataDog APM and DataDog Logs
// to links traces and log messages.
func TraceFieldGenerator(ctx context.Context) []zapcore.Field {
	spanCtx := trace.FromContext(ctx).SpanContext()

	return []zapcore.Field{
		zap.Uint64("dd.trace_id", binary.BigEndian.Uint64(spanCtx.TraceID[8:])),
		zap.Uint64("dd.span_id", binary.BigEndian.Uint64(spanCtx.SpanID[:])),
	}
}
