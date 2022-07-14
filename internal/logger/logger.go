package logger

import (
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"os"
)

// Create interface for logging
// Implement a test logger / real logger
// NOTE: "Logging" in this case is output for the user, except for Debug logging which can be more like log.Println
// This package might be good for that: https://github.com/logrusorgru/aurora
// Ideas on implementation: https://gogoapps.io/blog/passing-loggers-in-go-golang-logging-best-practices/ (context passing!?)

type Logs interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warn(v ...interface{})
	Error(v ...interface{})
}

var logger Logs

func GetLogger() Logs {
	if logger == nil {
		lvl := os.Getenv("LOG_LEVEL")
		if lvl == "" {
			lvl = "warn"
		}

		// todo: Make log level an option in vessel.yml
		// todo: Use a file writer to write to $HOME/.vessel/debug.log
		w := log.NewSyncWriter(os.Stderr)
		kitLogger := log.NewLogfmtLogger(w)
		kitLogger = level.NewFilter(kitLogger, level.Allow(level.ParseDefault(lvl, level.DebugValue())))
		kitLogger = log.With(kitLogger, "caller", log.DefaultCaller)

		logger = &Logger{
			Level: lvl,
			Base:  kitLogger,
		}
	}

	return logger
}

type Logger struct {
	Level string
	Base  log.Logger
}

func (l *Logger) Info(v ...interface{}) {
	level.Info(l.Base).Log(v...)
}

func (l *Logger) Warn(v ...interface{}) {
	level.Warn(l.Base).Log(v...)
}

func (l *Logger) Error(v ...interface{}) {
	level.Error(l.Base).Log(v...)
}

func (l *Logger) Debug(v ...interface{}) {
	level.Debug(l.Base).Log(v...)
}
