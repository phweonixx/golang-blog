package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

var db = database.New()

var ErrTitleTooShort = errors.New("title less than 3 characters")
var ErrCompanyExist = errors.New("user already has a company")

func CreateCompany(company *models.Company) error {
	timeNow := time.Now()
	company.CreatedAt = timeNow
	company.UpdatedAt = timeNow
	company.UUID = uuid.New().String()
	company.OwnerUUID = middleware.User_UUID

	// Перевірка на валідність введеного заголовку для компанії
	if company.Title == "" || len(company.Title) < 3 {
		return ErrTitleTooShort
	}

	// Перевірка на існування компанії у користувача

	exists, err := helpers.CheckExists(company.OwnerUUID, "company")
	if err != nil {
		log.Println("Error checking company existence:", err)
		return err
	}
	if exists {
		log.Println("User already has a company.")
		return ErrCompanyExist
	}

	// Запрос для створення компанії
	err = db.DBGorm.Create(&company).Error
	if err != nil {
		return err
	}

	return nil
}
