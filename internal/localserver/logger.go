package localserver

import (
	"fmt"

	awswebproxy "github.com/tfmcdigital/aws-web-proxy/internal"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func newProductionZaplogger(service string) (*zap.SugaredLogger, error) {
	conf := zap.NewProductionConfig()
	conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	conf.DisableCaller = true
	conf.DisableStacktrace = true
	conf.OutputPaths = remove(conf.OutputPaths, "stdout")
	conf.OutputPaths = remove(conf.OutputPaths, "stderr")
	zapLogger, err := conf.Build(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapCore(c, service)
	}))

	return zapLogger.Sugar(), err
}

func zapCore(c zapcore.Core, service string) zapcore.Core {
	// lumberjack.Logger is already safe for concurrent use, so we don't need to
	// lock it.
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/logs/%s.log", awswebproxy.BaseAwpPath(), service),
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
