package users

import (
	"blogAPI/internal/auth"
	"blogAPI/internal/database"
	"log"

	"gorm.io/gorm"
)

func SoftDeleteUser(uuid string) error {
	err := database.DBGorm.Model(&auth.User{}).
		Where("uuid = ?", uuid).
		Update("deleted_at", gorm.Expr("NOW()")).Error

	if err != nil {
		log.Println("Error deleting account")
		return err
	}
	return nil
}

func checkUserExists(uuid string) (bool, error) {
	var exists bool
	err := database.DBGorm.Model(&auth.User{}).
		Select("1").
		Where("uuid = ? AND deleted_at IS NULL", uuid).
		Limit(1).
		Find(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}
