package std

import (
	"go-clean-template/pkg/logger/common"
	"log"
	"os"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	dtMask     = time.RFC3339Nano
	basePrefix = ""
	baseFlag   = 0
)

const (
	defaultLevel   = infoLevel
	defaultLogFile = "/var/log/service.log"
)

const (
	debugLevel   = "DEBUG"
	infoLevel    = "INFO"
	warningLevel = "WARNING"
	errLevel     = "ERROR"
	fatalLevel   = "FATAL"
)

const (
	DEBUG   = 40
	INFO    = 30
	WARNING = 20
	ERROR   = 10
	FATAL   = 0
)

type StdLoggerOpts struct {
	Enabled bool
	Level   string
	LogFile string
	Stdout  bool
}

type Logger struct {
	env    string
	level  int
	logger *log.Logger
}

func NewLogger(selfOpts *StdLoggerOpts, opts *common.GeneralOpts) *Logger {
	level := selfOpts.Level
	if level == "" {
		level = defaultLevel
	}
	logFile := selfOpts.LogFile
	if logFile == "" {
		logFile = defaultLogFile
	}

	var logLevelMap = map[string]int{
		debugLevel:   DEBUG,
		infoLevel:    INFO,
		warningLevel: WARNING,
		errLevel:     ERROR,
		fatalLevel:   FATAL,
	}

	var logger *log.Logger
	if !selfOpts.Stdout {
		logger = log.New(&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    5,  //nolint:mnd //.
			MaxBackups: 20, //nolint:mnd //.
			MaxAge:     60, //nolint:mnd //.
		}, basePrefix, baseFlag)
	} else {
		logger = log.New(os.Stdout, basePrefix, baseFlag)
	}

	l := &Logger{
		env:    opts.Env,
		level:  logLevelMap[level],
		logger: logger,
	}

	return l
}

func (l *Logger) Close() {}

func (l *Logger) Debug(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= DEBUG {
		msg, _ := v[0].(string)
		l.logger.Print(now+" DEBUG ", common.GetFuncName(), " ", msg)
	}
}

func (l *Logger) Info(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= INFO {
		msg, _ := v[0].(string)
		l.logger.Print(now+" INFO ", common.GetFuncName(), " ", msg)
	}
}

func (l *Logger) Warning(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= WARNING {
		msg, _ := v[0].(string)
		l.logger.Print(now+" WARNING ", common.GetFuncName(), " ", msg)
	}
}

func (l *Logger) Error(v ...interface{}) {
	now := time.Now().Format(dtMask)
	if l.level >= ERROR {
		msg, _ := v[0].(string)
		l.logger.Print(now+" ERROR ", common.GetFuncName(), " ", msg)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	now := time.Now().Format(dtMask)
	msg, _ := v[0].(string)
	l.logger.Print(now+" FATAL ", common.GetFuncName(), " ", msg)
}
