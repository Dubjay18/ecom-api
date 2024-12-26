package config

import (
	"fmt"
	"log"

	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	// Redis   RedisConfig // For rate limiting and caching if needed
	APIKeys APIKeysConfig
}

type ServerConfig struct {
	Host            string        `mapstructure:"SERVER_HOST"`
	Port            string        `mapstructure:"PORT"`
	Mode            string        `mapstructure:"SERVER_MODE"`
	ReadTimeout     time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout    time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`
	ShutdownTimeout time.Duration `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
}

type BaseConfig struct {
	SERVER_PORT             string `mapstructure:"SERVER_PORT"`
	SERVER_HOST             string `mapstructure:"SERVER_HOST"`
	SERVER_MODE             string `mapstructure:"SERVER_MODE"`
	SERVER_READ_TIMEOUT     string `mapstructure:"SERVER_READ_TIMEOUT"`
	SERVER_WRITE_TIMEOUT    string `mapstructure:"SERVER_WRITE_TIMEOUT"`
	SERVER_SHUTDOWN_TIMEOUT string `mapstructure:"SERVER_SHUTDOWN_TIMEOUT"`
	SERVER_SECRET           string `mapstructure:"SERVER_SECRET"`

	DB_HOST string `mapstructure:"DB_HOST"`
	DB_PORT string `mapstructure:"DB_PORT"`
	DB_USER string `mapstructure:"DB_USER"`

	TIMEZONE    string `mapstructure:"TIMEZONE"`
	DB_SSL_MODE string `mapstructure:"DB_SSL_MODE"`

	DB_PASSWORD              string `mapstructure:"DB_PASSWORD"`
	DB_NAME                  string `mapstructure:"DB_NAME"`
	DB_MAX_IDLE_CONNS        int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DB_MAX_OPEN_CONNS        int    `mapstructure:"DB_MAX_OPEN_CONNS"`
	DB_MAX_LIFETIME          string `mapstructure:"DB_MAX_LIFETIME"`
	MIGRATE                  bool   `mapstructure:"MIGRATE"`
	CLOUDINARY_CLOUD_NAME    string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
	CLOUDINARY_KEY           string `mapstructure:"CLOUDINARY_KEY"`
	CLOUDINARY_SECRET        string `mapstructure:"CLOUDINARY_SECRET"`
	JWT_SECRET_KEY           string `mapstructure:"JWT_SECRET"`
	JWT_ACCESS_TOKEN_EXPIRY  string `mapstructure:"JWT_ACCESS_TOKEN_EXPIRY"`
	JWT_REFRESH_TOKEN_EXPIRY string `mapstructure:"JWT_REFRESH_TOKEN_EXPIRY"`
	IPSTACK_KEY              string `mapstructure:"IPSTACK_KEY"`
	IPSTACK_BASE_URL         string `mapstructure:"IPSTACK_BASE_URL"`

	MAIL_SERVER   string `mapstructure:"MAIL_SERVER"`
	MAIL_PASSWORD string `mapstructure:"MAIL_PASSWORD"`
	MAIL_USERNAME string `mapstructure:"MAIL_USERNAME"`
	MAIL_PORT     string `mapstructure:"MAIL_PORT"`

	REDIS_PORT string `mapstructure:"REDIS_PORT"`
	REDIS_HOST string `mapstructure:"REDIS_HOST"`
	REDIS_DB   string `mapstructure:"REDIS_DB"`
}
type DBConfig struct {
	Host         string        `mapstructure:"DB_HOST"`
	Port         string        `mapstructure:"DB_PORT"`
	User         string        `mapstructure:"DB_USER"`
	Password     string        `mapstructure:"DB_PASSWORD"`
	DBName       string        `mapstructure:"DB_NAME"`
	SSLMode      string        `mapstructure:"DB_SSL_MODE"`
	MaxIdleConns int           `mapstructure:"DB_MAX_IDLE_CONNS"`
	MaxOpenConns int           `mapstructure:"DB_MAX_OPEN_CONNS"`
	MaxLifetime  time.Duration `mapstructure:"DB_MAX_LIFETIME"`
}

type JWTConfig struct {
	SecretKey          string        `mapstructure:"JWT_SECRET_KEY"`
	AccessTokenExpiry  time.Duration `mapstructure:"JWT_ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry time.Duration `mapstructure:"JWT_REFRESH_TOKEN_EXPIRY"`
}

// type RedisConfig struct {
// 	Host     string
// 	Port     string
// 	Password string
// 	DB       int
// }

type APIKeysConfig struct {
	StripeKey           string `mapstructure:"STRIPE_KEY"`
	CloudinaryKey       string `mapstructure:"CLOUDINARY_KEY"`
	CloudinarySecret    string `mapstructure:"CLOUDINARY_SECRET"`
	CloudinaryCloudName string `mapstructure:"CLOUDINARY_CLOUD_NAME"`
}

// LoadConfig reads configuration from environment variables or config file
func LoadConfig(path string) (*Config, error) {
	v := viper.New()

	v.SetConfigFile(".env")

	// v.SetEnvPrefix("")
	// v.AutomaticEnv()
	// v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.ReadInConfig()
	// if err := ; err != nil {
	// 	log.Printf("Error reading config file, %s", err)
	// 	log.Printf("Falling back to environment variables")
	// }

	baseConfig := &BaseConfig{}
	if err := v.Unmarshal(baseConfig); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	log.Println(baseConfig.SERVER_PORT, "1")
	if err := BindKeys(v, baseConfig); err != nil {
		return nil, fmt.Errorf("unable to bind keys: %w", err)
	}

	config := &Config{
		Server: ServerConfig{
			Host:            baseConfig.SERVER_HOST,
			Port:            baseConfig.SERVER_PORT,
			Mode:            baseConfig.SERVER_MODE,
			ReadTimeout:     5 * time.Second,
			WriteTimeout:    5 * time.Second,
			ShutdownTimeout: 5 * time.Second,
		},
		DB: DBConfig{
			Host:         baseConfig.DB_HOST,
			Port:         baseConfig.DB_PORT,
			User:         baseConfig.DB_USER,
			Password:     baseConfig.DB_PASSWORD,
			DBName:       baseConfig.DB_NAME,
			SSLMode:      baseConfig.DB_SSL_MODE,
			MaxIdleConns: baseConfig.DB_MAX_IDLE_CONNS,
			MaxOpenConns: baseConfig.DB_MAX_OPEN_CONNS,
			MaxLifetime:  1 * time.Hour,
		},
		JWT: JWTConfig{
			SecretKey:          baseConfig.JWT_SECRET_KEY,
			AccessTokenExpiry:  24 * time.Hour,
			RefreshTokenExpiry: 7 * 24 * time.Hour,
		},
		APIKeys: APIKeysConfig{
			CloudinaryKey:       baseConfig.CLOUDINARY_KEY,
			CloudinarySecret:    baseConfig.CLOUDINARY_SECRET,
			CloudinaryCloudName: baseConfig.CLOUDINARY_CLOUD_NAME,
		},
	}

	return config, nil
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

func BindKeys(v *viper.Viper, input interface{}) error {

	envKeysMap := &map[string]interface{}{}
	if err := mapstructure.Decode(input, &envKeysMap); err != nil {
		return err
	}
	for k := range *envKeysMap {
		if bindErr := viper.BindEnv(k); bindErr != nil {
			return bindErr
		}
	}

	return nil
}
