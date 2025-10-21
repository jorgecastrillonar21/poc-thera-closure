package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Redis    RedisConfig    `mapstructure:"redis"`
	AWS      AWSConfig      `mapstructure:"aws"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name     string `mapstructure:"name"`
	Version  string `mapstructure:"version"`
	LogLevel string `mapstructure:"log_level"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret               string        `mapstructure:"secret"`
	AccessTokenDuration  time.Duration `mapstructure:"access_token_duration"`
	RefreshTokenDuration time.Duration `mapstructure:"refresh_token_duration"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AWSConfig holds AWS configuration
type AWSConfig struct {
	Region              string `mapstructure:"region"`
	CognitoUserPoolID   string `mapstructure:"cognito_user_pool_id"`
	CognitoClientID     string `mapstructure:"cognito_client_id"`
	CognitoClientSecret string `mapstructure:"cognito_client_secret"`
	AccessKeyID         string `mapstructure:"access_key_id"`
	SecretAccessKey     string `mapstructure:"secret_access_key"`
}

// Load loads configuration from environment variables and config files
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	// Set default values
	setDefaults()

	// Override with environment variables
	viper.AutomaticEnv()

	// Bind environment variables
	bindEnvVars()

	// Read config file (optional)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Validate required fields
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("app.name", "theraclosure-auth-service")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.log_level", "debug")

	// Server defaults
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", "3001")
	viper.SetDefault("server.mode", "debug")

	// Database defaults
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.ssl_mode", "disable")

	// JWT defaults
	viper.SetDefault("jwt.access_token_duration", "1h")
	viper.SetDefault("jwt.refresh_token_duration", "168h")

	// Redis defaults
	viper.SetDefault("redis.host", "localhost")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.db", 0)

	// AWS defaults
	viper.SetDefault("aws.region", "us-east-1")
}

func bindEnvVars() {
	// App
	viper.BindEnv("app.name", "APP_NAME")
	viper.BindEnv("app.version", "APP_VERSION")
	viper.BindEnv("app.log_level", "LOG_LEVEL")

	// Server
	viper.BindEnv("server.host", "SERVER_HOST")
	viper.BindEnv("server.port", "SERVER_PORT")
	viper.BindEnv("server.mode", "GIN_MODE")

	// Database
	viper.BindEnv("database.host", "DB_HOST")
	viper.BindEnv("database.port", "DB_PORT")
	viper.BindEnv("database.user", "DB_USER")
	viper.BindEnv("database.password", "DB_PASSWORD")
	viper.BindEnv("database.name", "DB_NAME")
	viper.BindEnv("database.ssl_mode", "DB_SSL_MODE")

	// JWT
	viper.BindEnv("jwt.secret", "JWT_SECRET")
	viper.BindEnv("jwt.access_token_duration", "JWT_ACCESS_TOKEN_DURATION")
	viper.BindEnv("jwt.refresh_token_duration", "JWT_REFRESH_TOKEN_DURATION")

	// Redis
	viper.BindEnv("redis.host", "REDIS_HOST")
	viper.BindEnv("redis.port", "REDIS_PORT")
	viper.BindEnv("redis.password", "REDIS_PASSWORD")
	viper.BindEnv("redis.db", "REDIS_DB")

	// AWS
	viper.BindEnv("aws.region", "AWS_REGION")
	viper.BindEnv("aws.cognito_user_pool_id", "AWS_COGNITO_USER_POOL_ID")
	viper.BindEnv("aws.cognito_client_id", "AWS_COGNITO_CLIENT_ID")
	viper.BindEnv("aws.cognito_client_secret", "AWS_COGNITO_CLIENT_SECRET")
	viper.BindEnv("aws.access_key_id", "AWS_ACCESS_KEY_ID")
	viper.BindEnv("aws.secret_access_key", "AWS_SECRET_ACCESS_KEY")
}

func validate(config *Config) error {
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT secret is required")
	}
	if config.Database.User == "" {
		return fmt.Errorf("database user is required")
	}
	if config.Database.Name == "" {
		return fmt.Errorf("database name is required")
	}
	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Database.Host,
		c.Database.User,
		c.Database.Password,
		c.Database.Name,
		c.Database.Port,
		c.Database.SSLMode,
	)
}

// GetRedisAddr returns the Redis address
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}
