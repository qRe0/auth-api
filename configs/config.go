package configs

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	errs "github.com/qRe0/auth-api/internal/errors"
)

type DBConfig struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

type RedisConfig struct {
	Host string
	Port string
}

type JWTConfig struct {
	SecretKey     string
	TTL           string
	BlacklistTime string
}

type Config struct {
	DB    DBConfig
	Redis RedisConfig
	JWT   JWTConfig
}

func LoadEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, errs.ErrLoadEnvVars
	}

	requiredEnvs := []string{
		"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT", "REDIS_HOST", "REDIS_PORT", "SECRET_KEY", "TTL", "BLACKLIST_TIME",
	}

	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			return nil, fmt.Errorf("environment variable `%s` is not set or is empty", env)
		}
	}

	config := Config{
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
		},
		Redis: RedisConfig{
			Host: os.Getenv("REDIS_HOST"),
			Port: os.Getenv("REDIS_PORT"),
		},
		JWT: JWTConfig{
			SecretKey:     os.Getenv("SECRET_KEY"),
			TTL:           os.Getenv("TTL"),
			BlacklistTime: os.Getenv("BLACKLIST_TIME"),
		},
	}

	return &config, nil
}
