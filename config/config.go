package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type App struct {
	Port        string
	Env         string
	MaxFileSize int64
}

type DB struct {
	Url string
}

type Store struct {
	BucketName string `mapstructure:"bucket_name"`
	Region     string `mapstructure:"region"`
}

type Config struct {
	App   App   `mapstructure:"app"`
	Store Store `mapstructure:"store"`
	DB    DB    `mapstructure:"db"`
}

func LoadConfig() (config *Config, err error) {
	if os.Getenv("APP_ENV") == "" {
		err := godotenv.Load(".env")
		if err != nil {
			return nil, err
		}
	}

	maxFileSizeMB, err := strconv.ParseInt(os.Getenv("APP_MAX_FILE_SIZE_MB"), 10, 64)
	if err != nil {
		return nil, err
	}
	appConfig := App{
		Port:        os.Getenv("APP_PORT"),
		Env:         os.Getenv("APP_ENV"),
		MaxFileSize: maxFileSizeMB,
	}

	storeConfig := Store{
		BucketName: os.Getenv("STORE_BUCKET_NAME"),
		Region:     os.Getenv("STORE_REGION"),
	}

	dbConfig := DB{
		Url: os.Getenv("DB_URL"),
	}

	return &Config{
		App:   appConfig,
		Store: storeConfig,
		DB:    dbConfig,
	}, nil
}

func (a *App) IsDevelopment() bool {
	return a.Env == "development"
}
