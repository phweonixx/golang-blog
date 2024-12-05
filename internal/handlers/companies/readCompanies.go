package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"log"
)

func GetCompanies(limit, offset int, UUID, owner_uuid string) ([]models.Company, error) {
	var companies []models.Company
	// Запрос для пошуку компаній по введеним значенням
	query := database.DBGorm.Model(&models.Company{}).Where("deleted_at IS NULL")

	if UUID != "" {
		query = query.Where("uuid = ?", UUID)
	}
	if owner_uuid != "" {
		query = query.Where("owner_uuid = ?", owner_uuid)
	}

	err := query.Limit(limit).Offset(offset).Find(&companies).Error
	if err != nil {
		log.Println("Error getting companies")
		return nil, err
	}

	// Повернення отриманого списку структур для виводу користувачу
	return companies, nil
}
