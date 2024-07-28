package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitializeLogger sets up and returns a zap.SugaredLogger instance based on the provided log level.
func InitializeLogger(level string) (*zap.SugaredLogger, error) {
	var logger *zap.Logger
	var err error

	// Set the log level
	logLevel := zap.NewAtomicLevel()
	switch level {
	case "debug":
		logLevel.SetLevel(zapcore.DebugLevel)
	case "info":
		logLevel.SetLevel(zapcore.InfoLevel)
	case "warn":
		logLevel.SetLevel(zapcore.WarnLevel)
	case "error":
		logLevel.SetLevel(zapcore.ErrorLevel)
	default:
		logLevel.SetLevel(zapcore.InfoLevel)
	}

	config := zap.Config{
		Level:    logLevel,
		Encoding: "json",
	}
	logger, err = config.Build()
	if err != nil {
		return nil, err
	}
	return logger.Sugar(), nil
}
