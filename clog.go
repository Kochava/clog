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

func sanitizeEnvRune(r rune) rune {
	r = unicode.ToUpper(r)
	if r == '_' || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return r
	}
	return -1
}

func envPrefix(procname string) string {
	if procname == "" && len(os.Args) > 0 {
		procname = filepath.Base(os.Args[0])
		procname = procname[:len(procname)-len(filepath.Ext(procname))]
	}
	procname = strings.Map(sanitizeEnvRune, procname)
	if procname == "" {
		return ""
	}
	return procname + "_"
}

// GetEnvConfig returns the environment-configured logging level and whether to use JSON and debug
// logging for procname. If procname is the empty string, os.Args[0] is used instead.
//
// If the PROCNAME_LEVEL value is invalid or not set, it defaults to zap.WarnLevel.
// If either the PROCNAME_JSON or PROCNAME_DEBUG, they default to false.
func GetEnvConfig(procname string) (level zapcore.Level, json, debug bool) {
	prefix := envPrefix(procname)
	json, _ = strconv.ParseBool(os.Getenv(prefix + "LOG_JSON"))
	debug, _ = strconv.ParseBool(os.Getenv(prefix + "LOG_DEBUG"))
	if txt, ok := os.LookupEnv(prefix + "LOG_LEVEL"); !ok || level.UnmarshalText([]byte(txt)) != nil {
		level = zap.WarnLevel
	}
	return
}

// NewFromEnv allocates a new zap.Logger using configuration from the environment. This looks for
// PROCNAME_LOG_JSON, PROCNAME_LOG_DEBUG, and PROCNAME_LOG_LEVEL to configure the logger. If JSON
// logs are used, the format is standard zap without modification.
func NewFromEnv(procname string, level zap.AtomicLevel) (*zap.Logger, error) {
	lvl, json, debug := GetEnvConfig(procname)
	if level != zeroLevel {
		level.SetLevel(lvl)
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
