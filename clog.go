// Package clog is a convenience package for passing Zap loggers through
// contexts.
package clog

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zeroLevel = zap.AtomicLevel{}

var osGetenv = os.Getenv

var osArg0 = func() string {
	return os.Args[0]
}

func sanitizeEnvRune(r rune) rune {
	r = unicode.ToUpper(r)
	if r == '_' || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
		return r
	}
	return -1
}

func envPrefix(procname string) string {
	if procname == "" && len(os.Args) > 0 {
		procname = filepath.ToSlash(osArg0())
	}
	procname = filepath.Base(procname)
	if ext := path.Ext(procname); ext != "" {
		procname = procname[:len(procname)-len(ext)]
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
// If PROCNAME_LOG_MODE is set to "dev" (case-insensitive) then log output will be formatted for
// reading on a console. Otherwise, logging defaults to a production configuration.
//
// If PROCNAME_LOG_LEVEL is set to a valid Zap logging level (info, warn, etc.), then that logging
// level will be returned. Otherwise, the logging level defaults to zap.InfoLevel for production and
// zap.DebugLevel for development.
func GetEnvConfig(procname string) (level zapcore.Level, isDev bool) {
	const devEnvironment = "dev"
	prefix := envPrefix(procname)
	isDev = strings.EqualFold(osGetenv(prefix+"LOG_MODE"), devEnvironment)
	levelText := osGetenv(prefix + "LOG_LEVEL")
	if levelText != "" && level.UnmarshalText([]byte(levelText)) == nil {
		// nop
	} else if isDev {
		level = zap.DebugLevel // development
	} else {
		level = zap.InfoLevel // production

	}
	return
}

// NewFromEnv allocates a new zap.Logger using configuration from the environment.
// This looks for PROCNAME_LOG_MODE and PROCNAME_LOG_LEVEL to configure the logger.
// If LOG_MODE is not "dev", the development configuration of Zap is used.
// Otherwise, logging is configured for production.
func NewFromEnv(procname string, level zap.AtomicLevel) (*zap.Logger, error) {
	lvl, isDev := GetEnvConfig(procname)
	if level != zeroLevel {
		level.SetLevel(lvl)
	}
	return New(level, isDev)
}

// New allocates a new zap.Logger using configuration based on the level given and the json and
// debug parameters, as interpreted by Config.
func New(level zap.AtomicLevel, isDev bool) (*zap.Logger, error) {
	return Config(level, isDev).Build()
}

// Config returns a zap.Config based on the level given and the json and debug parameters. If json
// is true, the config uses a JSON encoder. If debug is true, production limits on logging are
// removed and the development flag is set to true.
func Config(level zap.AtomicLevel, isDev bool) (conf zap.Config) {
	if isDev {
		conf = zap.NewDevelopmentConfig()
	} else {
		conf = zap.NewProductionConfig()
	}
	conf.Level = level
	return conf
}
