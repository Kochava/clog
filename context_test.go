package clog

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestContextLogging(t *testing.T) {
	levels := map[zapcore.Level]func(context.Context, string, ...zapcore.Field){
		zap.DPanicLevel: DPanic,
		zap.FatalLevel:  Fatal,
		zap.PanicLevel:  Panic,
		zap.DebugLevel:  Debug,
		zap.ErrorLevel:  Error,
		zap.InfoLevel:   Info,
		zap.WarnLevel:   Warn,
	}

	fstr := zap.String("f-string", "value")
	fint := zap.Int("f-int", 1234)

	const (
		msgNoFields   = "log message - no fields"
		msgWithFields = "log message - with fields"
	)

	root := context.Background()
	if Logger(root) != zap.L() {
		t.Error("default logger should be the global logger (zap.L)")
	}

	for lev, fun := range levels {
		lev, fun := lev, fun
		t.Run(lev.String(), func(t *testing.T) {
			core, logs := observer.New(zap.DebugLevel)
			logger := zap.New(core)
			ctx := WithLogger(root, logger)

			var (
				wantWithoutFields = 3
				wantWithFields    = 3
			)
			switch lev {
			case zap.DPanicLevel, zap.FatalLevel, zap.PanicLevel:
				wantWithoutFields = 2
				wantWithFields = 2
			default:
				fun(ctx, msgNoFields)
				fun(ctx, msgWithFields, fstr, fint)
			}
			wantTotal := wantWithoutFields + wantWithoutFields

			// Without fields
			check := Check(ctx, lev, msgNoFields)
			check.Should(check.Entry, zapcore.WriteThenNoop).Write()

			// With fields
			check = Check(ctx, lev, msgWithFields)
			check.Should(check.Entry, zapcore.WriteThenNoop).Write(fstr, fint)

			// With no extra fields
			with := With(ctx)
			check = Check(with, lev, msgNoFields)
			check.Should(check.Entry, zapcore.WriteThenNoop).Write()

			// With extra fields
			with = With(ctx, fstr, fint)
			check = Check(with, lev, msgWithFields)
			check.Should(check.Entry, zapcore.WriteThenNoop).Write()

			for _, l := range logs.All() {
				t.Log(l)
			}

			if got, want := logs.Len(), wantTotal; got != want {
				t.Fatalf("received %d log messages; want %d", got, want)
			}

			noFields := logs.FilterMessage(msgNoFields)
			if got, want := noFields.Len(), wantWithoutFields; got != want {
				t.Errorf("expected %d messages with message %q; got %d", got, msgNoFields, want)
			}

			withFields := logs.FilterMessage(msgWithFields).FilterField(fstr).FilterField(fint)
			if got, want := withFields.Len(), wantWithFields; got != want {
				t.Errorf("expected %d messages with message %q; got %d", got, msgWithFields, want)
			}
		})
	}
}
