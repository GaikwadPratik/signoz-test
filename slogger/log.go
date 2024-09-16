package slogger

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"
	"time"
)

type ConfigLogger struct {
	// ProcessTitle title of the process will be required only for bunyan logger
	ProcessTitle string
	// LogLevel minimum level of the slogger, this can be used to change level in runtime
	LogLevel *slog.LevelVar
}

// ConfigureLogger updates option on logger and returns instance
// this instance needs to be set using slog.SetDefault(logger)
func ConfigureLogger(config ConfigLogger) *slog.Logger {
	if config.LogLevel == nil {
		panic("log level is empty")
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("unable to determine working directory: %v", err))
	}

	homeDirname, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("unable to get home directory of current user: %v", err))
	}

	//Check if it is in Docker env
	if _, err := os.Stat(path.Join(homeDirname, ".dockerenv")); !errors.Is(err, os.ErrNotExist) {
		replacer := func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.SourceKey {
				source := a.Value.Any().(*slog.Source)
				if file, ok := strings.CutPrefix(source.File, wd); ok {
					source.File = file
				}
			}

			if a.Key == slog.TimeKey {
				t := a.Value.Time()
				a.Value = slog.StringValue(t.Format(time.RFC3339))
			}

			return a
		}

		options := &slog.HandlerOptions{
			Level:       config.LogLevel,
			ReplaceAttr: replacer,
			AddSource:   true,
		}

		return slog.New(slog.NewJSONHandler(os.Stderr, options))
	}

	// if not in Docker env, then config the logger to work with Bunyan
	replacer := func(groups []string, a slog.Attr) slog.Attr {
		// Use relative path instead of absolute
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			if file, ok := strings.CutPrefix(source.File, wd); ok {
				source.File = file
			}
		}

		if a.Key == slog.TimeKey {
			t := a.Value.Time()
			a.Value = slog.StringValue(t.Format(time.RFC3339))
		}

		// Convert level to Bunyan integer
		if a.Key == slog.LevelKey {
			return slog.Int(a.Key, bunyanLevel(a.Value.Any().(slog.Level)))
		}

		return a
	}

	opttions := &slog.HandlerOptions{
		Level:       config.LogLevel,
		ReplaceAttr: replacer,
		AddSource:   true,
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("unable to get hostname of the system: %v", err))
	}

	logAttribs := []slog.Attr{
		slog.Int("pid", os.Getpid()),
		slog.Int("v", 0),
		slog.String("hostname", hostname),
		slog.String("name", config.ProcessTitle),
	}

	return slog.New(slog.NewJSONHandler(os.Stderr, opttions).WithAttrs(logAttribs))
}

// bunyanLevel maps slog levels to Bunyan levels
func bunyanLevel(level slog.Level) int {
	switch level {
	case slog.LevelDebug:
		return 20
	case slog.LevelInfo:
		return 30
	case slog.LevelWarn:
		return 40
	case slog.LevelError:
		return 50
	default:
		return 30 //deafult to info level
	}
}
