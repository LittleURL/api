package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger = zap.Config{
	Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
	Development: false,
	Sampling: &zap.SamplingConfig{
		Initial:    100,
		Thereafter: 100,
	},
	Encoding:         "json",
	EncoderConfig:    logEncoder,
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}

var logEncoder = zapcore.EncoderConfig{
	TimeKey:        "timestamp",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "caller",
	FunctionKey:    zapcore.OmitKey,
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.SecondsDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}
