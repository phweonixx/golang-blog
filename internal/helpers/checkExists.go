package helpers

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"errors"
	"log"
)

func CheckExists[generic int | string](idOrUUID generic, typeOfElement string) (bool, error) {
	db := database.New()
	var count int64

	if typeOfElement == "article" {
		err := db.DBGorm.Model(&models.Article{}).
			Where("id = ?", idOrUUID).
			Count(&count).Error
		if err != nil {
			return false, err
		}

		return count > 0, nil
	} else if typeOfElement == "category" {
		err := db.DBGorm.Model(&models.Category{}).
			Where("id = ?", idOrUUID).
			Count(&count).Error
		if err != nil {
			return false, err
		}
		return count > 0, nil
	} else if typeOfElement == "company" {
		err := db.DBGorm.Model(&models.Company{}).
			Where("uuid = ? AND deleted_at IS NULL", idOrUUID).
			Count(&count).Error
		if err != nil {
			log.Println("Error checking company existence:", err)
			return false, err
		}

		return count > 0, nil
	}

	return false, errors.New("wrong name of element")
}

func CheckCompanyExistsByOwner(ownerUUID string) (bool, error) {
	db := database.New()
	var count int64
	err := db.DBGorm.Model(&models.Company{}).
		Where("owner_uuid = ? AND deleted_at IS NULL", ownerUUID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func CheckUserExistsLogin(username, email string) (bool, error) {
	db := database.New()
	var exists bool
	err := db.DBGorm.Model(&models.User{}).
		Select("1").
		Where("(username = ? OR email = ?) AND deleted_at IS NULL", username, email).
		Limit(1).
		Find(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}
