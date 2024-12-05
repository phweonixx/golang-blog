package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Декларація структур, що завантажуються з environment файлу
type MainConfig struct {
	DatabaseConfig DBCfg
	Config         Cfg
	JWTConfig      JWTCfg
}

type DBCfg struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type Cfg struct {
	Port        string
	DefaultLang string
	Score       string
}

type JWTCfg struct {
	JWTSecret     []byte
	TokenLifetime time.Duration
}

var MainCfg MainConfig

func New() *MainConfig {
	// Завантаження environment файлу
	err := godotenv.Load("../pkg/.env")
	if err != nil {
		log.Fatal("Error loading .env file\n", err)
	}

	// Створення екземпляра конфігурації
	cfg := &MainConfig{}
	cfg.LoadAllConfigs()
	return cfg
}

func (mc *MainConfig) LoadAllConfigs() {
	mc.DatabaseConfig = mc.DatabaseConfig.LoadDatabaseConfig()
	mc.Config = mc.Config.LoadConfig()
	mc.JWTConfig = mc.JWTConfig.LoadJWTConfig()
}

func (dbCfg DBCfg) LoadDatabaseConfig() DBCfg {
	return DBCfg{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
}

func (cfg Cfg) LoadConfig() Cfg {
	return Cfg{
		Port:        os.Getenv("PORT"),
		DefaultLang: os.Getenv("DEFAULT_LANGUAGE"),
		Score:       os.Getenv("SCORE"),
	}
}

func (jc JWTCfg) LoadJWTConfig() JWTCfg {
	tokenLifetime, err := time.ParseDuration(os.Getenv("JWT_TOKEN_LIFETIME"))
	if err != nil {
		log.Fatal("Invalid JWT_TOKEN_LIFETIME format:", err)
	}

	return JWTCfg{
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		TokenLifetime: tokenLifetime,
	}
}

/*func loadDatabaseConfig() DBCfg {
	return DBCfg{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
	}
} */

/* func loadConfig() cfg {
	return cfg{
		Port:        os.Getenv("PORT"),
		DefaultLang: os.Getenv("DEFAULT_LANGUAGE"),
		Score:       os.Getenv("SCORE"),
	}
} */

/* func loadJWTConfig() JWTCfg {
	tokenLifetime, err := time.ParseDuration(os.Getenv("JWT_TOKEN_LIFETIME"))
	if err != nil {
		log.Fatal("Invalid JWT_TOKEN_LIFETIME format:", err)
	}

	return JWTCfg{
		JWTSecret:     []byte(os.Getenv("JWT_SECRET")),
		TokenLifetime: tokenLifetime,
	}
} */
