package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

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

	exists, err := checkCompanyExists(company.OwnerUUID)
	if err != nil {
		log.Println("Error checking company existence:", err)
		return err
	}
	if exists {
		log.Println("User already has a company.")
		return ErrCompanyExist
	}

	// Запрос для створення компанії
	err = database.DBGorm.Create(&company).Error
	if err != nil {
		return err
	}

	return nil
}

func checkCompanyExists(ownerUUID string) (bool, error) {
	var count int64
	err := database.DBGorm.Model(&models.Company{}).
		Where("owner_uuid = ? AND deleted_at IS NULL", ownerUUID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
