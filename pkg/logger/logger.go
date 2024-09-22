package logger

import (
	"bytes"
	"fmt"
	"go-clean-template/config"
	"go-clean-template/pkg/logger/common"
	"go-clean-template/pkg/logger/slog"
	"go-clean-template/pkg/logger/std"
	"go-clean-template/pkg/logger/telegram"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LoggerOpts struct {
	Opts               *common.GeneralOpts
	TelegramLoggerOpts *telegram.TelegramLoggerOpts
	StdLoggerOpts      *std.StdLoggerOpts
	SlogLoggerOpts     *slog.SlogLoggerOpts
}

type Logger interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	// Fatal writes log message with fatal level and os.Exit(1) after
	Fatal(v ...interface{})
	Close()
}

type logger struct {
	lgs []Logger
}

func New(opts *LoggerOpts) Logger {
	l := &logger{}
	if opts.StdLoggerOpts.Enabled {
		l.lgs = append(l.lgs, std.NewLogger(opts.StdLoggerOpts, opts.Opts))
	}
	if opts.TelegramLoggerOpts.Enabled {
		l.lgs = append(l.lgs, telegram.NewLogger(opts.TelegramLoggerOpts, opts.Opts))
	}
	if opts.SlogLoggerOpts.Enabled {
		l.lgs = append(l.lgs, slog.NewLogger(opts.SlogLoggerOpts, opts.Opts))
	}
	return l
}

func (l *logger) Close() {
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Close()
		}
	}
}

func (l *logger) Debug(v ...interface{}) {
	msg := concat(v...)
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Debug(msg)
		}
	}
}

func (l *logger) Info(v ...interface{}) {
	msg := concat(v...)
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Info(msg)
		}
	}
}

func (l *logger) Warning(v ...interface{}) {
	msg := concat(v...)
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Warning(msg)
		}
	}
}

func (l *logger) Error(v ...interface{}) {
	msg := concat(v...)
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Error(msg)
		}
	}
}

func (l *logger) Fatal(v ...interface{}) {
	msg := concat(v...)
	for i := range l.lgs {
		if l.lgs[i] != nil {
			l.lgs[i].Fatal(msg)
		}
	}
	time.Sleep(2 * time.Second)
	os.Exit(1)
}

func concat(v ...interface{}) string {
	var buffer bytes.Buffer
	for i, s := range v {
		if i == len(v)-1 {
			buffer.WriteString(fmt.Sprintf("%v", s))
		} else {
			buffer.WriteString(fmt.Sprintf("%v ", s))
		}
	}
	return buffer.String()
}

func MakeLoggerOpts(c *config.Config) *LoggerOpts {
	return &LoggerOpts{
		Opts: &common.GeneralOpts{
			AppVersion: c.AppVersion,
			InstanceID: c.InstanceID,
			Env:        c.Env,
			AppName:    c.AppName,
		},
		TelegramLoggerOpts: &telegram.TelegramLoggerOpts{
			Enabled:      c.Logger.LoggerTelegram.Enabled,
			Level:        c.Logger.LoggerStd.Level,
			TargetChatID: c.Logger.LoggerTelegram.TargetChatID,
			BotAPIToken:  c.Logger.LoggerTelegram.BotAPIToken,
		},
		StdLoggerOpts: &std.StdLoggerOpts{
			Enabled: c.Logger.LoggerStd.Enabled,
			Level:   c.Logger.LoggerStd.Level,
			LogFile: c.Logger.LoggerStd.LogFile,
			Stdout:  c.Logger.LoggerStd.Stdout,
		},
		SlogLoggerOpts: &slog.SlogLoggerOpts{
			Enabled: c.Logger.LoggerSlog.Enabled,
			Level:   c.Logger.LoggerSlog.Level,
			JSON:    c.Logger.LoggerSlog.JSON,
		},
	}
}

func GetFuncName() string {
	var buffer bytes.Buffer
	const pcSize = 10
	pc := make([]uintptr, pcSize)
	const skip = 4
	runtime.Callers(skip, pc)
	frame, _ := runtime.CallersFrames(pc).Next()
	function := frame.Function
	line := frame.Line
	buffer.WriteString(function)
	buffer.WriteString(fmt.Sprintf(":%d", line))

	return filepath.Base(buffer.String())
}
