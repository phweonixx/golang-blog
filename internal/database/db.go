package database

import (
	"blogAPI/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Змінна бази даних
var DB *sql.DB
var DBGorm *gorm.DB

// Підключення до бази даних
func InitDB() {
	// Настройка з використанням змінних з .env файлу

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.DatabaseConfig.User,
		config.DatabaseConfig.Password,
		config.DatabaseConfig.Host,
		config.DatabaseConfig.Port,
		config.DatabaseConfig.Name,
	)

	// Відкриття бази даних GORM
	var err error
	DBGorm, err = gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database connection:\n", err)
	} else {
		log.Println("GORM database connected successfully!")
	}

	gormDB, err := DBGorm.DB()
	if err != nil {
		log.Fatal("Error getting raw DB from GORM:\n", err)
	}

	// Перевірка підключення
	err = gormDB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database via GORM:\n", err)
	}

}
