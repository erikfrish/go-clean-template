package slog

import (
	"go-clean-template/pkg/logger/common"
	"go-clean-template/pkg/logger/slog/handlers/prettyslog"
	"log/slog"
	"os"
)

const (
	defaultLevel = "INFO"
)

type SlogLoggerOpts struct {
	Enabled bool
	Level   string
	JSON    bool
}

type Logger struct {
	env    string
	level  int
	logger *slog.Logger
}

func NewLogger(selfOpts *SlogLoggerOpts, opts *common.GeneralOpts) *Logger {
	levelText := selfOpts.Level
	if levelText == "" {
		levelText = defaultLevel
	}

	var logger *slog.Logger
	var level slog.Level
	err := level.UnmarshalText([]byte(levelText))
	if err != nil {
		level = slog.LevelInfo
	}

	if selfOpts.JSON {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
	} else {
		logger = prettyslog.SetupPrettySlog(os.Stdout, level)
	}

	l := &Logger{
		env:    opts.Env,
		level:  int(level),
		logger: logger,
	}
	return l
}

func (l *Logger) Close() {}

func (l *Logger) Debug(v ...interface{}) {
	l.logger.Debug(v[0].(string), "caller", common.GetFuncName())
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.Info(v[0].(string))
}

func (l *Logger) Warning(v ...interface{}) {
	l.logger.Warn(v[0].(string))
}

func (l *Logger) Error(v ...interface{}) {
	l.logger.Error(v[0].(string), "caller", common.GetFuncName())
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logger.Error(v[0].(string))
}
