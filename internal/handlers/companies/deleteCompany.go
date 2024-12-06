package companies

import (
	"blogAPI/internal/models"
	"log"
)

func SoftDeleteCompany(uuid string) error {
	err := db.DBGorm.Where("uuid = ?", uuid).Delete(&models.Company{}).Error
	if err != nil {
		log.Println("Error deleting the company!")
		return err
	}

	return nil
}
