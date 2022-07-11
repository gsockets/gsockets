package log

import (
	"os"

	kitlog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type Fields map[string]interface{}

type Logger interface {
	With(keyvals ...interface{}) Logger

	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
}

type logger struct {
	log kitlog.Logger
}

func New() Logger {
	log := kitlog.NewJSONLogger(kitlog.NewSyncWriter(os.Stdout))
	log = kitlog.With(log, "ts", kitlog.DefaultTimestampUTC)

	return &logger{log: log}
}

func (l *logger) With(keyvals ...interface{}) Logger {
	newLogger := &logger{}
	newLogger.log = kitlog.With(l.log, keyvals...)

	return newLogger
}

func (l *logger) Debug(args ...interface{}) {
	level.Debug(l.log).Log(args...)
}

func (l *logger) Info(args ...interface{}) {
	level.Info(l.log).Log(args...)
}

func (l *logger) Warn(args ...interface{}) {
	level.Warn(l.log).Log(args...)
}

func (l *logger) Error(args ...interface{}) {
	level.Error(l.log).Log(args...)
}

func (l *logger) Fatal(args ...interface{}) {
	level.Error(l.log).Log(args...)
	os.Exit(1)
}
