package logger

import (
	"fmt"
	"log"
	"sync"

	"github.com/tfmcdigital/aws-web-proxy/internal/domain"
	"github.com/tfmcdigital/aws-web-proxy/internal/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var loggers sync.Map = sync.Map{}

func GetLogger(service string) *zap.SugaredLogger {
	actual, _ := loggers.LoadOrStore(service, newProductionZaplogger(service))
	return actual.(*zap.SugaredLogger)
}

func newProductionZaplogger(service string) *zap.SugaredLogger {
	conf := zap.NewProductionConfig()
	conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	conf.DisableCaller = true
	conf.DisableStacktrace = true
	conf.OutputPaths = utils.Remove(conf.OutputPaths, "stdout")
	conf.OutputPaths = utils.Remove(conf.OutputPaths, "stderr")
	zapLogger, err := conf.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapCore(c, service)
	}))
	if err != nil {
		log.Default().Fatal("Failed to init zap logger", err)
	}
	return zapLogger.Sugar()
}

func zapCore(c zapcore.Core, service string) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   GetLogFileLocation(service),
		MaxSize:    50, // megabytes
		MaxBackups: 30,
		MaxAge:     28, // days
	})
	conf := zap.NewProductionEncoderConfig()
	conf.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(conf),
		w,
		zap.DebugLevel,
	)
	cores := zapcore.NewTee(c, core)

	return cores
}

func GetLogFileLocation(service string) string {
	return fmt.Sprintf("%s/logs/%s.log", domain.BasePath, service)
}
