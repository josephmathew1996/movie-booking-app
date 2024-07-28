package logger_test

import (
	"movie-booking-app/users-service/internal/logger"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestInitializeLogger(t *testing.T) {
	tests := []struct {
		level         string
		expectedLevel zapcore.Level
	}{
		{"debug", zapcore.DebugLevel},
		{"info", zapcore.InfoLevel},
		{"warn", zapcore.WarnLevel},
		{"error", zapcore.ErrorLevel},
		{"invalid", zapcore.InfoLevel}, // Default case
	}
	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			logger, err := logger.InitializeLogger(tt.level)
			assert.NoError(t, err)
			assert.NotNil(t, logger)
			assert.Equal(t, tt.expectedLevel, logger.Level())
		})
	}
}
