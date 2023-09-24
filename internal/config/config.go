package config

import (
	"github.com/joho/godotenv"
	"github.com/zenorachi/todo-service/pkg/database/postgres"
	"github.com/zenorachi/todo-service/pkg/logger"
	"sync"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

type Config struct {
	HTTP HTTPConfig
	Auth AuthConfig
	GIN  GINConfig
	DB   postgres.DBConfig
}

type (
	HTTPConfig struct {
		Host         string
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	AuthConfig struct {
		AccessTokenTTL  time.Duration
		RefreshTokenTTL time.Duration
		Salt            string
		Secret          string
	}

	GINConfig struct {
		Mode string
	}
)

var (
	config = &Config{}
	once   sync.Once
)

const (
	envFile   = ".env"
	directory = "configs"
	ymlFile   = "main"
)

func init() {
	if err := godotenv.Load(envFile); err != nil {
		logger.Fatal("config", ".env initialization failed")
	}

	viper.AddConfigPath(directory)
	viper.SetConfigName(ymlFile)
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("config", "viper initialization failed")
	}
}

func New() *Config {
	once.Do(func() {
		if err := viper.Unmarshal(config); err != nil {
			logger.Fatal("viper config", err.Error())
		}

		if err := envconfig.Process("db", &config.DB); err != nil {
			logger.Fatal("db config", err.Error())
		}

		if err := envconfig.Process("hash", &config.Auth); err != nil {
			logger.Fatal("hash envs", err.Error())
		}

		if err := envconfig.Process("gin", &config.GIN); err != nil {
			logger.Fatal("gin config", err.Error())
		}
	})

	return config
}
