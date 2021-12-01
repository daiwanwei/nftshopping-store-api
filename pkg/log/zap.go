package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"nftshopping-store-api/pkg/config"
	"os"
)

var logInstance Logger

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Panic(args ...interface{})
	Fatal(args ...interface{})
	DebugF(template string, args ...interface{})
	InfoF(template string, args ...interface{})
	WarnF(template string, args ...interface{})
	ErrorF(template string, args ...interface{})
	PanicF(template string, args ...interface{})
	FatalF(template string, args ...interface{})
}

type logger struct {
	sugarLogger *zap.SugaredLogger
}

func GetLog() (instance Logger, err error) {
	if logInstance == nil {
		instance, err = newLog()
		if err != nil {
			panic(err)
		}
		logInstance = instance
	}
	return logInstance, nil
}

func newLog() (Logger, error) {
	c, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	loggerConfig := c.Logger
	level := getLoggerLevel(loggerConfig)
	logWriter := zapcore.AddSync(os.Stderr)
	var encoderCfg zapcore.EncoderConfig
	if loggerConfig.Mode == "dev" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	var encoder zapcore.Encoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MESSAGE"

	if loggerConfig.Encoding == "console" {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		fmt.Println("QQQQQ")
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(encoder, logWriter, zap.NewAtomicLevelAt(level))
	sugarLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	defer sugarLogger.Sync()
	log := &logger{sugarLogger: sugarLogger}
	return log, nil
}

func (l *logger) Debug(args ...interface{}) {
	l.sugarLogger.Debug(args...)
}

func (l *logger) DebugF(template string, args ...interface{}) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.sugarLogger.Info(args...)
}

func (l *logger) InfoF(template string, args ...interface{}) {
	l.sugarLogger.Infof(template, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.sugarLogger.Warn(args...)
}

func (l *logger) WarnF(template string, args ...interface{}) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.sugarLogger.Error(args...)
}

func (l *logger) ErrorF(template string, args ...interface{}) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *logger) DPanic(args ...interface{}) {
	l.sugarLogger.DPanic(args...)
}

func (l *logger) DPanicF(template string, args ...interface{}) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *logger) Panic(args ...interface{}) {
	l.sugarLogger.Panic(args...)
}

func (l *logger) PanicF(template string, args ...interface{}) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.sugarLogger.Fatal(args...)
}

func (l *logger) FatalF(template string, args ...interface{}) {
	l.sugarLogger.Fatalf(template, args...)
}

var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dPanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(loggerConfig *config.Logger) zapcore.Level {
	level, exist := loggerLevelMap[loggerConfig.Level]
	if !exist {
		return zapcore.DebugLevel
	}
	return level
}
