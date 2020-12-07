package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var Logger *zap.SugaredLogger
var LogLevel zap.AtomicLevel

func init() {
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		LogLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		LogLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		LogLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		LogLevel = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	default:
		LogLevel = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	config := zap.Config{
		Level:             LogLevel,
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Encoding:          "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "log_message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "@timestamp",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			NameKey:    "name",
			EncodeName: zapcore.FullNameEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,

			StacktraceKey: "stacktrace",
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	unSugaredlogger, _ := config.Build()
	Logger = unSugaredlogger.Sugar()
}
