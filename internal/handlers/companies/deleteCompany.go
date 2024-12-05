package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"log"
)

func SoftDeleteCompany(uuid string) error {
	err := database.DBGorm.Where("uuid = ?", uuid).Delete(&models.Company{}).Error
	if err != nil {
		log.Println("Error deleting the company!")
		return err
	}

	return nil
}

func checkCompanyExistsByUUID(UUID string) (bool, error) {
	var count int64
	err := database.DBGorm.Model(&models.Company{}).
		Where("uuid = ? AND deleted_at IS NULL", UUID).
		Count(&count).Error
	if err != nil {
		log.Println("Error checking company existence:", err)
		return false, err
	}

	return count > 0, nil
}
