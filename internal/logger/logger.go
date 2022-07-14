package logger

import (
	"fmt"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/vessel-app/vessel-cli/internal/util"
	"io"
	"os"
	"path/filepath"
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
	Close() error
}

var logger Logs

type LogWriterCLoser interface {
	io.Writer
	io.Closer
}

func GetLogger() Logs {
	if logger == nil {
		lvl := os.Getenv("LOG_LEVEL")

		if lvl == "" {
			lvl = "warn"
		}

		writer, err := logWriter()

		if err != nil {
			writer = os.Stderr
		}

		w := log.NewSyncWriter(writer)
		kitLogger := log.NewLogfmtLogger(w)
		kitLogger = level.NewFilter(kitLogger, level.Allow(level.ParseDefault(lvl, level.DebugValue())))
		kitLogger = log.With(kitLogger, "ts", log.DefaultTimestampUTC)

		logger = &Logger{
			Level:  lvl,
			Base:   kitLogger,
			Writer: writer,
		}
	}

	return logger
}

func logWriter() (LogWriterCLoser, error) {
	storagePath, err := util.MakeStorageDir()

	if err != nil {
		return nil, fmt.Errorf("could not create storage directory: %w", err)
	}

	file, err := os.OpenFile(filepath.FromSlash(storagePath+"/debug.log"), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)

	if err != nil {
		return nil, fmt.Errorf("could not open ~/.vessel/debug.log for writing: %w", err)
	}

	return file, nil
}

type Logger struct {
	Level  string
	Base   log.Logger
	Writer LogWriterCLoser
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

func (l *Logger) Close() error {
	return l.Writer.Close()
}
