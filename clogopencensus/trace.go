package clogopencensus

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// TraceFieldGenerator generates dd.trace_id and dd.span_id fields.
//
// Adding dd.trace_id and dd.span_id to a json formatted log message allows DataDog APM and DataDog Logs
// to links traces and log messages.
func TraceFieldGenerator(ctx context.Context) []zapcore.Field {
	var (
		spanCtx = trace.FromContext(ctx).SpanContext()
		traceID = spanCtx.TraceID.String()
		spanID  = spanCtx.SpanID.String()
	)

	return []zapcore.Field{
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	}
}
