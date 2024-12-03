package auth

import (
	"time"

	"gorm.io/gorm"
)

// Модель для акаунту
type User struct {
	UUID      string          `json:"uuid" gorm:"column:uuid"`
	FirstName string          `json:"first_name" gorm:"column:first_name"`
	LastName  string          `json:"last_name" gorm:"column:last_name"`
	Username  string          `json:"username" gorm:"column:username"`
	Password  string          `json:"password" gorm:"column:password"`
	Email     string          `json:"email" gorm:"column:email"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
}

func (User) TableName() string {
	return "user"
}
