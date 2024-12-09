package database

import (
	"blogAPI/internal/config"
	"blogAPI/internal/models"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Змінна бази даних
type DB struct {
	DBGorm *gorm.DB
}

// Підключення до бази даних
func New() *DB {
	// Настройка з використанням змінних з .env файлу
	cfg := config.New()

	dbConnectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		cfg.DatabaseConfig.User,
		cfg.DatabaseConfig.Password,
		cfg.DatabaseConfig.Host,
		cfg.DatabaseConfig.Port,
		cfg.DatabaseConfig.Name,
	)

	// Відкриття бази даних GORM
	gormDB, err := gorm.Open(mysql.Open(dbConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Error opening database connection:\n", err)
	}

	rawDB, err := gormDB.DB()
	if err != nil {
		log.Fatal("Error getting raw DB from GORM:\n", err)
	}

	// Перевірка підключення
	err = rawDB.Ping()
	if err != nil {
		log.Fatal("Error connecting to database via GORM:\n", err)
	}

	err = gormDB.AutoMigrate(
		&models.Article{},
		&models.Category{},
		&models.Company{},
		/*&models.RelatedArticles{},*/
		&models.Translations{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Migrations failed.")
	}

	return &DB{
		DBGorm: gormDB,
	}
}
