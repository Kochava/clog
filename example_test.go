package clog_test

import (
	"context"

	"github.com/Kochava/clog"
	"go.uber.org/zap"
)

func ExampleNew() {
	// Create a logger at info level with a production configuration.
	level := zap.NewAtomicLevelAt(zap.InfoLevel)
	l, err := clog.New(level, false)
	if err != nil {
		panic(err)
	}
	l.Info("Ready")

	// Attach the logger, l, to a context:
	ctx := clog.WithLogger(context.Background(), l)

	// Attach fields to the logger:
	ctx = clog.With(ctx, zap.Int("field", 1234))

	// Log at info level:
	clog.Info(ctx, "Log message")
}
