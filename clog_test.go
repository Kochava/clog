package clog

import (
	"sync"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var arg0Lock sync.Mutex

func setArg0(t *testing.T, v string) func() {
	arg0Lock.Lock()
	cur := osArg0
	osArg0 = func() string {
		t.Logf("Requested os.Args[0] (%q)", v)
		return v
	}
	return func() { defer arg0Lock.Unlock(); osArg0 = cur }
}

var getenvLock sync.Mutex

func setGetenv(t *testing.T, values map[string]string) func() {
	getenvLock.Lock()
	cur := osGetenv
	osGetenv = func(k string) string {
		t.Logf("Requested os.Getenv(%q)", values[k])
		return values[k]
	}
	return func() { defer getenvLock.Unlock(); osGetenv = cur }
}

func TestEnvPrefix(t *testing.T) {
	cases := []struct {
		In, Want string
	}{
		{"./Foo_", "FOO__"},
		{"./Foo-", "FOO_"},
		{"/usr/sbin/health-checker", "HEALTHCHECKER_"},
		{"/usr/sbin/health_checker", "HEALTH_CHECKER_"},
		{"foo_bar", "FOO_BAR_"},
	}

	for _, c := range cases {
		t.Run("ProcnamePass="+c.In, func(t *testing.T) {
			c := c
			want, got := c.Want, envPrefix(c.In)
			if want != got {
				t.Errorf("envPrefix(%q) = %q; want %q", c.In, got, want)
			}
		})
	}

	for _, c := range cases {
		c := c
		t.Run("ProcnameArg0="+c.In, func(t *testing.T) {
			defer setArg0(t, c.In)()
			want, got := c.Want, envPrefix("")
			if want != got {
				t.Errorf("envPrefix(%q) = %q; want %q", c.In, got, want)
			}
		})
	}

	t.Run("EmptyPass", func(t *testing.T) {
		const (
			in   = "/usr/bin/clog-daemon"
			want = "CLOGDAEMON_"
		)
		defer setArg0(t, in)()
		if got := envPrefix(""); want != got {
			t.Errorf("envPrefix(%q) = %q; want %q", in, got, want)
		}
	})

	t.Run("NoValues", func(t *testing.T) {
		const (
			in   = ""
			want = ""
		)
		defer setArg0(t, in)()
		if got := envPrefix(""); want != got {
			t.Errorf("envPrefix(%q) = %q; want %q", in, got, want)
		}
	})
}

func TestGetEnvConfig(t *testing.T) {
	type values map[string]string
	cases := []struct {
		Case   string
		Arg0   string
		Values values
		Level  zapcore.Level
		IsDev  bool
	}{
		{
			"Defaults",
			"",
			nil,
			zap.InfoLevel,
			false,
		},
		{
			"EmptyName",
			"",
			values{"LOG_MODE": "dev", "LOG_LEVEL": "warn"},
			zap.WarnLevel,
			true,
		},
		{
			"InvalidLevel",
			"",
			values{"LOG_MODE": "prod", "LOG_LEVEL": "WARNING"},
			zap.InfoLevel,
			false,
		},
		{
			"NamedProc",
			"/usr/local/bin/daemon",
			values{"DAEMON_LOG_MODE": "dev", "DAEMON_LOG_LEVEL": "fatal"},
			zap.FatalLevel,
			true,
		},
	}

	for _, c := range cases {
		c := c
		t.Run(c.Case, func(t *testing.T) {
			defer setArg0(t, c.Arg0)()
			defer setGetenv(t, c.Values)()

			level, isDev := GetEnvConfig("")
			if level != c.Level {
				t.Errorf("level = %v; want %v", level, c.Level)
			}
			if isDev != c.IsDev {
				t.Errorf("isDev = %t; want %t", isDev, c.IsDev)
			}
		})
	}
}
