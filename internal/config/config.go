package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Декларація структур, що завантажуються з environment файлу
type DBCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type cfg struct {
	Port        string
	DefaultLang string
	Score       string
}

type JWTCfg struct {
	JWTSecret     []byte
	TokenLifetime time.Duration
}

var DatabaseConfig DBCfg
var Config cfg
var JWTConfig JWTCfg

func init() {
	// Завантаження environment файлу
	err := godotenv.Load("../pkg/.env")
	if err != nil {
		log.Fatal("Error loading .env file\n", err)
	}

	// Змінні бази даних
	DatabaseConfig = loadDatabaseConfig()
	Config = loadConfig()
	JWTConfig = loadJWTConfig()
}

func loadDatabaseConfig() DBCfg {
	return DBCfg{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func loadConfig() cfg {
	return cfg{
		Port:        os.Getenv("PORT"),
		DefaultLang: os.Getenv("DEFAULT_LANGUAGE"),
		Score:       os.Getenv("SCORE"),
	}
}

func loadJWTConfig() JWTCfg {
	tokenLifetime, err := time.ParseDuration(os.Getenv("JWT_TOKEN_LIFETIME"))
	if err != nil {
		log.Fatal("Invalid JWT_TOKEN_LIFETIME format:", err)
	}

	return JWTCfg{
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		TokenLifetime: tokenLifetime,
	}
}
