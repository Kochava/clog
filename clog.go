package clog

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogFunc is any logging function, such as methods on zap.Logger.
type LogFunc func(string, ...zapcore.Field)

var zeroLevel = zap.AtomicLevel{}

// NewFromEnv allocates a new zap.Logger using configuration from the environment. This looks for
// PROCNAME_LOG_JSON, PROCNAME_LOG_DEBUG, and PROCNAME_LOG_LEVEL to configure the logger. If JSON
// logs are used, the format is standard zap without modification.
func NewFromEnv(procname string, level zap.AtomicLevel) (*zap.Logger, error) {
	var json, debug bool
	prefix := ""
	if len(os.Args) > 0 || procname != "" {
		if procname == "" {
			if len(os.Args) == 0 {
				// Skip prefix
				goto configure
			}
			procname = filepath.Base(os.Args[0])
			procname = procname[:len(procname)-len(filepath.Ext(procname))]
		}
		prefix = strings.Map(func(r rune) rune {
			if r = unicode.ToUpper(r); r != '_' && !unicode.IsLetter(r) {
				r = -1
			}
			return r
		}, procname) + "_"
	}

configure:
	json, _ = strconv.ParseBool(os.Getenv(prefix + "LOG_JSON"))
	debug, _ = strconv.ParseBool(os.Getenv(prefix + "LOG_DEBUG"))
	if level != zeroLevel {
		var lvl zapcore.Level
		if txt, ok := os.LookupEnv(prefix + "LOG_LEVEL"); ok && lvl.UnmarshalText([]byte(txt)) == nil {
			level.SetLevel(lvl)
		}
	}

	return New(level, json, debug)
}

// New allocates a new zap.Logger using configuration based on the level given and the json and
// debug parameters, as interpreted by Config.
func New(level zap.AtomicLevel, json, debug bool) (*zap.Logger, error) {
	return Config(level, json, debug).Build()
}

// ShortTimeEncoder is a time encoder that records short, glog-like times. This is only intended for
// use with console-based logs. All standard options should be used for production configurations.
func ShortTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = `01-02 15:04:05.000`
	enc.AppendString(t.UTC().Format(layout))
}

// Config returns a zap.Config based on the level given and the json and debug parameters. If json
// is true, the config uses a JSON encoder. If debug is true, production limits on logging are
// removed and the development flag is set to true.
func Config(level zap.AtomicLevel, json, debug bool) zap.Config {
	cfg := zap.NewProductionConfig()
	if level != zeroLevel {
		cfg.Level = level
	}

	if debug {
		cfg.Sampling = nil
		cfg.Development = true
	}

	if json {
		cfg.Encoding = "json"
	} else {
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	}

	return cfg
}
