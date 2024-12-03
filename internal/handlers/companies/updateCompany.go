package companies

import (
	"blogAPI/internal/database"
	"log"
	"time"
)

func UpdateCompany(company *Company, UUID string) error {
	// Задання часу оновлення компанії
	company.UpdatedAt = time.Now()

	// Запрос для оновлення компанії
	err := database.DBGorm.Model(&Company{}).
		Where("uuid = ?", UUID).
		Updates(Company{
			UpdatedAt: company.UpdatedAt,
			Title:     company.Title,
		}).Error

	if err != nil {
		log.Println("Error updating company:", err)
		return err
	}

	return nil
}
