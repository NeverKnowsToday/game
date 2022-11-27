package logging

import (
	"github.com/op/go-logging"
	"os"
	"sync"
)

type Logger struct {
	Logger         *logging.Logger
	Module         string `json:"module"`
	Level          string `json:"level"`
	LeveledBackend logging.LeveledBackend
}

const (
	CRITICAL = "critical"
	ERROR    = "error"
	WARNING  = "warning"
	NOTICE   = "notice"
	INFO     = "info"
	DEBUG    = "debug"
)

var DEFAULT_LEVEL = DEBUG

var levelMap = map[string]logging.Level{
	CRITICAL: logging.CRITICAL,
	ERROR:    logging.ERROR,
	WARNING:  logging.WARNING,
	NOTICE:   logging.NOTICE,
	INFO:     logging.INFO,
	DEBUG:    logging.DEBUG,
}

var loggerMap = make(map[string]*Logger)
var locker = new(sync.Mutex)

func GetLogFile(filename string) (*os.File, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		f, err := os.Create(filename)
		return f, err
	} else {
		f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0600)
		return f, err
	}
}

func GetLogger(module, level string) *Logger {
	locker.Lock()
	defer locker.Unlock()

	if myLogger, ok := loggerMap[module]; ok {
		return myLogger
	}

	logger := logging.MustGetLogger(module)

	logger.ExtraCalldepth = 1
	file, _ := GetLogFile("./game.log")
	format := logging.MustStringFormatter("%{shortfile} %{time:15:04:05.000} [%{module}] %{level:.4s} : %{message}")
	backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
	backendStderrFormatter := logging.NewBackendFormatter(backendStderr, format)
	levelBackendStderr := logging.AddModuleLevel(backendStderrFormatter)
	levelBackendStderr.SetLevel(levelMap[level], module)

	backendFile := logging.NewLogBackend(file, "", 0)
	backendFileFormatter := logging.NewBackendFormatter(backendFile, format)
	levelBackendFile := logging.AddModuleLevel(backendFileFormatter)
	levelBackendFile.SetLevel(levelMap[level], module)

	multi := logging.SetBackend(levelBackendStderr, levelBackendFile)
	multi.SetLevel(levelMap[level], module)

	logging.SetBackend(multi).SetLevel(levelMap[level], module)

	myLogger := &Logger{
		Logger: logger,
		Module: module,
		Level:  level,
	}

	loggerMap[module] = myLogger

	return myLogger
}

func GetLoggerLevel(module string) string {
	level := ""
	if module == "" {
		for _, log := range loggerMap {
			level = log.Level
		}
	} else {
		log, ok := loggerMap[module]
		if ok {
			level = log.Level
		}
	}

	return level
}

func SetLevel(module, level string) {
	if module != "" {
		file, _ := GetLogFile("./game.log")
		format := logging.MustStringFormatter("%{shortfile} %{time:15:04:05.000} [%{module}] %{level:.4s} : %{message}")
		backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
		backendStderrFormatter := logging.NewBackendFormatter(backendStderr, format)
		levelBackendStderr := logging.AddModuleLevel(backendStderrFormatter)
		levelBackendStderr.SetLevel(levelMap[level], module)

		backendFile := logging.NewLogBackend(file, "", 0)
		backendFileFormatter := logging.NewBackendFormatter(backendFile, format)
		levelBackendFile := logging.AddModuleLevel(backendFileFormatter)
		levelBackendFile.SetLevel(levelMap[level], module)

		multi := logging.MultiLogger(levelBackendStderr, levelBackendFile)
		multi.SetLevel(levelMap[level], module)

		logging.SetBackend(multi).SetLevel(levelMap[level], module)
		log := loggerMap[module]
		log.Level = level
	} else {
		for tmpModule, _ := range loggerMap {
			file, _ := GetLogFile("./game.log")
			format := logging.MustStringFormatter("%{shortfile} %{time:15:04:05.000} [%{module}] %{level:.4s} : %{message}")
			backendStderr := logging.NewLogBackend(os.Stderr, "", 0)
			backendStderrFormatter := logging.NewBackendFormatter(backendStderr, format)
			levelBackendStderr := logging.AddModuleLevel(backendStderrFormatter)
			levelBackendStderr.SetLevel(levelMap[level], module)

			backendFile := logging.NewLogBackend(file, "", 0)
			backendFileFormatter := logging.NewBackendFormatter(backendFile, format)
			levelBackendFile := logging.AddModuleLevel(backendFileFormatter)
			levelBackendFile.SetLevel(levelMap[level], module)

			multi := logging.MultiLogger(levelBackendStderr, levelBackendFile)
			multi.SetLevel(levelMap[level], module)

			logging.SetBackend(multi).SetLevel(levelMap[level], module)

			log := loggerMap[tmpModule]
			log.Level = level
		}
	}
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Logger.Debug(args...)
}

func (logger *Logger) Debugf(template string, args ...interface{}) {
	logger.Logger.Debugf(template, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Logger.Info(args...)
}

func (logger *Logger) Infof(template string, args ...interface{}) {
	logger.Logger.Infof(template, args...)
}

func (logger *Logger) Notice(args ...interface{}) {
	logger.Logger.Notice(args...)
}

func (logger *Logger) Noticef(template string, args ...interface{}) {
	logger.Logger.Noticef(template, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.Logger.Warning(args...)
}

func (logger *Logger) Warningf(template string, args ...interface{}) {
	logger.Logger.Warningf(template, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Logger.Error(args...)
}

func (logger *Logger) Errorf(template string, args ...interface{}) {
	logger.Logger.Errorf(template, args...)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Logger.Panic(args...)
}

func (logger *Logger) Panicf(template string, args ...interface{}) {
	logger.Logger.Panicf(template, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Logger.Fatal(args...)
}

func (logger *Logger) Fatalf(template string, args ...interface{}) {
	logger.Logger.Fatalf(template, args...)
}
