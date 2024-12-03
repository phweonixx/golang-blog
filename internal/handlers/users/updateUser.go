package users

import (
	"blogAPI/internal/auth"
	"blogAPI/internal/database"
	"log"
	"time"
)

func UpdateUser(user auth.User, UUID string) error {
	user.UpdatedAt = time.Now()

	updateUser := auth.User{
		UpdatedAt: user.UpdatedAt,
	}

	if user.FirstName != "" {
		updateUser.FirstName = user.FirstName
	}
	if user.LastName != "" {
		updateUser.LastName = user.LastName
	}
	if user.Username != "" {
		updateUser.Username = user.Username
	}
	if user.Password != "" {
		updateUser.Password = user.Password
	}
	if user.Email != "" {
		updateUser.Email = user.Email
	}

	err := database.DBGorm.Model(&auth.User{}).
		Where("uuid = ?", UUID).
		Updates(updateUser).Error

	if err != nil {
		return err
	}

	return nil
}

func CheckUserUsernameOrEmailExists(username, email string) (bool, error) {
	var count int64
	err := database.DBGorm.Model(&auth.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error

	if err != nil {
		log.Println("Error checking user existence:", err)
		return false, err
	}

	return count > 0, nil
}
