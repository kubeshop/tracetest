package global_decorators

import (
	"os"

	global_types "github.com/kubeshop/tracetest/cli/global/types"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logger struct {
	global_types.Command
	Logger     *zap.Logger
	Parameters LoggerParams
}

type Logger interface {
	global_types.Command
	SetLogger(*zap.Logger)
	GetLogger() *zap.Logger
}

type LoggerParams struct {
	Verbose bool
}

var _ Logger = &logger{}

func WithLogger(command global_types.Command) global_types.Command {
	logger := &logger{
		Command: command,
	}

	cmd := logger.Get()
	cmd.PreRun = logger.preRun(cmd.PreRun)
	cmd.PostRun = logger.postRun(cmd.PostRun)
	cmd.PersistentFlags().BoolVarP(&logger.Parameters.Verbose, "verbose", "v", false, "display debug information")

	logger.Set(cmd)

	return logger
}

func (d *logger) SetLogger(logger *zap.Logger) {
	d.Logger = logger
}

func (d *logger) GetLogger() *zap.Logger {
	return d.Logger
}

func (d *logger) preRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		atom := zap.NewAtomicLevel()
		if d.Parameters.Verbose {
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

		logger := zap.New(zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderCfg),
			zapcore.Lock(os.Stdout),
			atom,
		))

		d.SetLogger(logger)

		next(cmd, args)
	}
}

func (d *logger) postRun(next CobraFn) CobraFn {
	return func(cmd *cobra.Command, args []string) {
		defer d.Logger.Sync()

		next(cmd, args)
	}
}
