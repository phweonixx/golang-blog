package models

import (
	"time"

	"gorm.io/gorm"
)

// Модель для акаунту
type User struct {
	UUID      string          `json:"uuid" gorm:"column:uuid;type:varchar(36)"`
	FirstName string          `json:"first_name" gorm:"column:first_name;type:varchar(255)"`
	LastName  string          `json:"last_name" gorm:"column:last_name;type:varchar(255)"`
	Username  string          `json:"username" gorm:"column:username;type:varchar(45)"`
	Password  string          `json:"password" gorm:"column:password;type:varchar(255)"`
	Email     string          `json:"email" gorm:"column:email;type:varchar(255)"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
}

func (User) TableName() string {
	return "user"
}
