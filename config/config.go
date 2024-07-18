package config

import (
	"log"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

var once sync.Once

type (
	Database struct {
		Host     string `mapstructure:"host" validate:"required"`
		Port     int    `mapstructure:"port" validate:"required"`
		User     string `mapstructure:"user" validate:"required"`
		Password string `mapstructure:"password" validate:"required"`
		DBName   string `mapstructure:"dbname" validate:"required"`
		SSLMode  string `mapstructure:"sslmode" validate:"required"`
		Schema   string `mapstructure:"schema" validate:"required"`
	}

	AppInfo struct {
		Name    string `mapstructure:"name" validate:"required"`
		Version string `mapstructure:"version" validate:"required"`
		Env     string `mapstructure:"env" validate:"required"`
	}

	Server struct {
		Port         int           `mapstructure:"port" validate:"required"`
		AllowOrigins []string      `mapstructure:"allowOrigins" validate:"required"`
		Timeout      time.Duration `mapstructure:"timeout" validate:"required"`
		BodyLimit    string        `mapstructure:"bodyLimit" validate:"required"`
	}

	Config struct {
		Database *Database `mapstructure:"database" validate:"required"`
		Server   *Server   `mapstructure:"server" validate:"required"`
		AppInfo  *AppInfo  `mapstructure:"appinfo" validate:"required"`
	}
)

var configInstance *Config

// Can refactor later it easy when change .env to .yaml
func ConfigGetting() *Config {
	once.Do(func() {
		viper.SetConfigFile(".env")
		viper.AddConfigPath(".")
		viper.AutomaticEnv()

		viper.SetEnvPrefix("APP")
		viper.BindEnv("HOST", "APP_HOST")
		viper.BindEnv("PORT", "APP_PORT")
		viper.BindEnv("NAME", "APP_NAME")
		viper.BindEnv("VERSION", "APP_VERSION")
		viper.BindEnv("ENV", "APP_ENV")

		viper.SetEnvPrefix("DB")
		viper.BindEnv("HOST", "DB_HOST")
		viper.BindEnv("PORT", "DB_PORT")
		viper.BindEnv("USER", "DB_USER")
		viper.BindEnv("PASSWORD", "DB_PASSWORD")
		viper.BindEnv("DBNAME", "DB_NAME")
		viper.BindEnv("SSLMODE", "DB_SSLMODE")
		viper.BindEnv("SCHEMA", "DB_SCHEMA")

		viper.SetEnvPrefix("SERVER")
		viper.BindEnv("PORT", "SERVER_PORT")
		viper.BindEnv("ALLOWORIGINS", "SERVER_ALLOW_ORIGINS")
		viper.BindEnv("TIMEOUT", "SERVER_TIMEOUT")
		viper.BindEnv("BODYLIMIT", "SERVER_BODY_LIMIT")

		if err := viper.ReadInConfig(); err != nil {
			log.Fatalf("Error reading config file: %v", err)
		}

		configInstance = &Config{
			Database: &Database{
				Host:     viper.GetString("DB_HOST"),
				Port:     viper.GetInt("DB_PORT"),
				User:     viper.GetString("DB_USER"),
				Password: viper.GetString("DB_PASSWORD"),
				DBName:   viper.GetString("DB_NAME"),
				SSLMode:  viper.GetString("DB_SSLMODE"),
				Schema:   viper.GetString("DB_SCHEMA"),
			},
			Server: &Server{
				Port:         viper.GetInt("SERVER_PORT"),
				AllowOrigins: viper.GetStringSlice("SERVER_ALLOW_ORIGINS"),
				Timeout:      viper.GetDuration("SERVER_TIMEOUT"),
				BodyLimit:    viper.GetString("SERVER_BODY_LIMIT"),
			},
			AppInfo: &AppInfo{
				Name:    viper.GetString("APP_NAME"),
				Version: viper.GetString("APP_VERSION"),
				Env:     viper.GetString("APP_ENV"),
			},
		}

		validate := validator.New()
		if err := validate.Struct(configInstance); err != nil {
			log.Fatalf("Validation failed: %v", err)
		}
	})
	return configInstance
}
