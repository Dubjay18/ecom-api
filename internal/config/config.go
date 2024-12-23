package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Server  ServerConfig
	DB      DBConfig
	JWT     JWTConfig
	Redis   RedisConfig // For rate limiting and caching if needed
	APIKeys APIKeysConfig
}

type ServerConfig struct {
	Host            string
	Port            string
	Mode            string // development, production, test
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

type DBConfig struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxIdleConns int
	MaxOpenConns int
	MaxLifetime  time.Duration
}

type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type APIKeysConfig struct {
	StripeKey     string // For payment processing
	CloudinaryKey string // For image storage if needed
}

// LoadConfig reads configuration from environment variables or config file
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(path)
	v.AutomaticEnv()

	// Default values
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found; ignore error if desired
		fmt.Println("No config file found. Using environment variables and defaults.")
	}

	config := &Config{}
	err := v.Unmarshal(config)
	if err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	return config, nil
}

// setDefaults sets default values for configuration
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", "8080")
	v.SetDefault("server.mode", "development")
	v.SetDefault("server.readTimeout", time.Second*5)
	v.SetDefault("server.writeTimeout", time.Second*5)
	v.SetDefault("server.shutdownTimeout", time.Second*5)

	// Database defaults
	v.SetDefault("db.host", "localhost")
	v.SetDefault("db.port", "5432")
	v.SetDefault("db.sslmode", "disable")

	// Redis defaults
	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", "6379")
	v.SetDefault("redis.db", 0)

	// JWT defaults
	v.SetDefault("jwt.accessTokenExpiry", time.Hour*24)    // 24 hours
	v.SetDefault("jwt.refreshTokenExpiry", time.Hour*24*7) // 7 days
}

// GetDSN returns database connection string
func (c *DBConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode)
}

// GetServerAddress returns formatted server address
func (c *ServerConfig) GetServerAddress() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
