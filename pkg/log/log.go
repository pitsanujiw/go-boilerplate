package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pitsanujiw/go-boilerplate/config"
)

type Logger struct {
	*zap.Logger
}

func New(cfg *config.App) (*Logger, error) {
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	conf := zap.NewProductionEncoderConfig()
	conf.MessageKey = "message"
	consoleEncoder := zapcore.NewJSONEncoder(conf)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	logger := zap.New(core, zap.Fields(zap.String("service", cfg.Log.Name)))

	return &Logger{
		Logger: logger,
	}, nil
}

func (l *Logger) Close() error {
	return l.Sync()
}
