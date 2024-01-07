package lib

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var lumlog = &lumberjack.Logger{
	MaxSize:    10, // megabytes
	MaxBackups: 30, // number of log files
	MaxAge:     1,  // days
}

func lumberjackZapHook(e zapcore.Entry) error {
	lumlog.Write([]byte(fmt.Sprintf("%+v", e)))
	return nil
}

func SetUpLogger() *zap.SugaredLogger {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(getLogLevel("debug")),
		OutputPaths:      []string{"stdout", "logs/it_" + time.Now().Format("02012006") + ".log"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	_logger, _ := cfg.Build(zap.Hooks(lumberjackZapHook))
	return _logger.Sugar()
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "DEBUG", "debug":
		return zapcore.DebugLevel
	case "ERROR", "error":
		return zapcore.ErrorLevel
	case "WARN", "warn":
		return zapcore.WarnLevel
	case "INFO", "info":
		return zapcore.InfoLevel
	default:
		return zapcore.ErrorLevel
	}
}
