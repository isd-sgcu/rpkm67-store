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
	Endpoint   string
	AccessKey  string
	SecretKey  string
	UseSSL     bool
	BucketName string
	Region     string
	Token      string
}

type Config struct {
	App   App   `mapstructure:"app"`
	Store Store `mapstructure:"store"`
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
		Endpoint:   os.Getenv("STORE_ENDPOINT"),
		AccessKey:  os.Getenv("STORE_ACCESS_KEY"),
		SecretKey:  os.Getenv("STORE_SECRET_KEY"),
		UseSSL:     os.Getenv("STORE_USE_SSL") == "true",
	}

	return &Config{
		App:   appConfig,
		Store: storeConfig,
	}, nil
}

func (a *App) IsDevelopment() bool {
	return a.Env == "development"
}
