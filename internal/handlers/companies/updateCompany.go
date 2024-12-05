package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"log"
	"time"
)

func UpdateCompany(company *models.Company, UUID string) error {
	// Задання часу оновлення компанії
	company.UpdatedAt = time.Now()

	// Запрос для оновлення компанії
	err := database.DBGorm.Model(&models.Company{}).
		Where("uuid = ?", UUID).
		Updates(models.Company{
			UpdatedAt: company.UpdatedAt,
			Title:     company.Title,
		}).Error

	if err != nil {
		log.Println("Error updating company:", err)
		return err
	}

	return nil
}
