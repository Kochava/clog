package clogdatadog_test

import (
	"context"

	"go.opencensus.io/trace"
	"go.uber.org/zap"

	"github.com/Kochava/clog"
	"github.com/Kochava/clog/clogopencensus/clogdatadog"
	"github.com/Kochava/clog/opencensus/clogdatadog"
)

func Example() {
	logger, _ := zap.NewProduction()

	ctx := clog.WithLogger(context.Background(), logger)

	// inform clog to always add datadog tracing fields
	ctx = clog.WithFieldGenerators(ctx, clogdatadog.TraceFieldGenerator)

	// start a span which will add a OpenCensus span to the context
	ctx, span := trace.StartSpan(ctx, "example")
	defer span.End()

	clog.Info(ctx, "some message")
	// {"level":"info","ts":1613673168.4960983,"caller":"clogdatadog/trace_test.go:22","msg":"some message","dd.span_id":"23eacfbb31546f15","dd.trace_id":"999c1ef5eb47a7082b1e2ec22213df6b"}
}
