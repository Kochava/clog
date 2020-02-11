clog
===

[![GoDoc](https://godoc.org/github.com/Kochava/clog?status.svg)](https://godoc.org/github.com/Kochava/clog)

    go get github.com/Kochava/clog

clog is a simple package for initializing a [Zap][] logger and attaching it to
a context, along with functions for logging from the context-attached logger or
associating new fields to the logger.

Generally speaking this is a bad use of the context package, but utility won out
over passing both a context and a logger around all the time. In particular,
this is useful for passing a request-scoped logger through different
http.Handler implementations that otherwise do not support Zap.

[Zap]: https://go.uber.org/zap


Usage
---

A few examples of basic usage follow.

### Initialize a logger

```go
// Create a logger at info level with a production configuration.
level := zap.NewAtomicLevelAt(zap.InfoLevel)
l, err := clog.New(level, false)
if err != nil {
    panic(err)
}
l.Info("Ready")
```

### Attach a logger to a context

```go
// var l *zap.Logger

// Attach the logger, l, to a context:
ctx := clog.WithLogger(context.Background(), l)

// Attach fields to the logger:
ctx = clog.With(ctx, zap.Int("field", 1234))

// Log at info level:
clog.Info(ctx, "Log message")
```


License
---

clog is made available under the ISC license. A copy of it can be found in the
repository in the `COPYING` file.
