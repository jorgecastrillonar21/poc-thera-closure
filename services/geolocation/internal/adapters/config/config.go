package config


import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	// Server configuration
	Server struct {
		Port string `mapstructure:"port"`
		Host string `mapstructure:"host"`
	} `mapstructure:"server"`

	// Database configuration
	Database struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
	} `mapstructure:"database"`

	// CORS configuration
	CORS struct {
		AllowedOrigins []string `mapstructure:"allowed_origins"`
		AllowedMethods []string `mapstructure:"allowed_methods"`
		AllowedHeaders []string `mapstructure:"allowed_headers"`
	} `mapstructure:"cors"`

	// Application configuration
	App struct {
		Name     string `mapstructure:"name"`
		Version  string `mapstructure:"version"`
		LogLevel string `mapstructure:"log_level"`
	} `mapstructure:"app"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}

	// Set default values
	viper.SetDefault("server.port", "3003")
	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "theraclosure")
	viper.SetDefault("database.password", "password123")
	viper.SetDefault("database.dbname", "theraclosure_geolocation")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("cors.allowed_origins", []string{"http://localhost:3000", "http://localhost:3001", "http://localhost:3002", "http://localhost:3003"})
	viper.SetDefault("cors.allowed_methods", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	viper.SetDefault("cors.allowed_headers", []string{"Origin", "Content-Type", "Accept", "Authorization"})
	viper.SetDefault("app.name", "theraclosure-geolocation-service")
	viper.SetDefault("app.version", "1.0.0")
	viper.SetDefault("app.log_level", "info")

	// Configure viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/theraclosure/")

	// Enable environment variable reading
	viper.AutomaticEnv()

	// Try to read configuration file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error since we have defaults
	}

	// Unmarshal configuration
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Override with environment variables
	if serverHost := os.Getenv("SERVER_HOST"); serverHost != "" {
		config.Server.Host = serverHost
	}

	if serverPort := os.Getenv("SERVER_PORT"); serverPort != "" {
		config.Server.Port = serverPort
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		config.Database.Host = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		config.Database.Port = dbPort
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		config.Database.User = dbUser
	}

	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		config.Database.Password = dbPassword
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		config.Database.DBName = dbName
	}

	if dbSSLMode := os.Getenv("DB_SSL_MODE"); dbSSLMode != "" {
		config.Database.SSLMode = dbSSLMode
	}

	if appName := os.Getenv("APP_NAME"); appName != "" {
		config.App.Name = appName
	}

	if appVersion := os.Getenv("APP_VERSION"); appVersion != "" {
		config.App.Version = appVersion
	}

	if logLevel := os.Getenv("APP_LOG_LEVEL"); logLevel != "" {
		config.App.LogLevel = logLevel
	}

	return config, nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Server.Host, c.Server.Port)
}