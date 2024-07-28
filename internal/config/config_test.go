// internal/config/config_test.go
package config_test

import (
	"movie-booking-app/users-service/internal/config"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitializeConfig(t *testing.T) {
	// Set up environment variables
	os.Setenv("APP_NAME", "test-app")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("SERVER_RUN_MODE", "debug")
	os.Setenv("LOG_LEVEL", "info")
	os.Setenv("LOG_FORMAT", "json")
	os.Setenv("RATELIMITER_REQUESTS_COUNT", "100")
	os.Setenv("RATELIMITER_ENABLED", "true")

	// Call the function
	cfg := config.InitializeConfig()

	// Verify the results
	assert.Equal(t, "test-app", cfg.App.Name)
	assert.Equal(t, 8080, cfg.App.Port)
	assert.Equal(t, "debug", cfg.App.RunMode)
	assert.Equal(t, "info", cfg.App.LogLevel)
	assert.Equal(t, "json", cfg.App.LogFormat)
	assert.Equal(t, 100, cfg.RateLimiter.RequestsPerTimeFrame)
	assert.Equal(t, true, cfg.RateLimiter.Enabled)
	assert.Equal(t, time.Minute*1, cfg.RateLimiter.TimeFrame)
}
