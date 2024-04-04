package cmdutil

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

type loggerConfig struct {
	Verbose bool
}

type loggerOption func(*loggerConfig)

func WithVerbose(verbose bool) loggerOption {
	return func(c *loggerConfig) {
		c.Verbose = verbose
	}
}

func GetLogger(opts ...loggerOption) *zap.Logger {
	if logger != nil {
		return logger
	}

	loggerConfig := loggerConfig{}
	for _, opt := range opts {
		opt(&loggerConfig)
	}

	atom := zap.NewAtomicLevel()
	if loggerConfig.Verbose {
		atom.SetLevel(zap.DebugLevel)
	} else {
		atom.SetLevel(zap.WarnLevel)
	}

	encoderCfg := zapcore.EncoderConfig{
		TimeKey:        zapcore.OmitKey,
		LevelKey:       "level",
		NameKey:        zapcore.OmitKey,
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	return zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))
}
