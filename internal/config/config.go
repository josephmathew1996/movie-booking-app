// internal/config/config.go
package config

import (
	"movie-booking-app/users-service/internal/ratelimiter"
	"time"

	"github.com/spf13/viper"
)

type (
	Config struct {
		App         AppConfig
		RateLimiter ratelimiter.Config
		// Db        DbConfig
	}

	// AppConfig holds application config settings.
	AppConfig struct {
		Name               string
		Port               int
		RunMode            string
		LogLevel           string
		LogFormat          string
		DBConnectionString string
	}

	// DbConfig holds settings to connect to the db.
	// DbConfig struct {
	// 	Host               string
	// 	DatabaseName       string
	// 	User               string
	// 	Password           string
	// 	UseSSL             *bool
	// 	SSLCertificate     string
	// 	SSLMode            string
	// 	MaxOpenConnections int
	// 	MaxIdleConnections int
	// 	MigrationsPath     string
	// }
)

// InitializeConfig sets up Viper to read from environment variables.
func InitializeConfig() Config {
	// Read environment variables
	viper.AutomaticEnv()

	return Config{
		App: AppConfig{
			Name:               viper.GetString("APP_NAME"),
			Port:               viper.GetInt("APP_PORT"),
			RunMode:            viper.GetString("SERVER_RUN_MODE"),
			LogLevel:           viper.GetString("LOG_LEVEL"),
			LogFormat:          viper.GetString("LOG_FORMAT"),
			DBConnectionString: viper.GetString("DATABASE_URL"),
		},
		RateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: viper.GetInt("RATELIMITER_REQUESTS_COUNT"),
			TimeFrame:            time.Minute * 1,
			Enabled:              viper.GetBool("RATELIMITER_ENABLED"),
		},

		// Db: DbConfig{
		// 	Host:               viper.GetString("DB_HOST"),
		// 	DatabaseName:       viper.GetString("DB_NAME"),
		// 	User:               viper.GetString("DB_USER"),
		// 	Password:           viper.GetString("DB_PASSWORD"),
		// 	SSLMode:            viper.GetString("DB_SSLMODE"),
		// 	UseSSL:             utils.InitBoolPtr(viper.GetBool("USE_SSL")),
		// 	SSLCertificate:     viper.GetString("DB_SSL_CA"),
		// 	MaxOpenConnections: viper.GetInt("DB_MAX_OPEN_CONNECTIONS"),
		// 	MaxIdleConnections: viper.GetInt("DB_MAX_IDLE_CONNECTIONS"),
		// 	MigrationsPath:     viper.GetString("DB_MIGRATIONS_PATH"),
		// },
	}
}
