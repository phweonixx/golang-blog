package models

import (
	"time"

	"gorm.io/gorm"
)

// Модель для компанії
type Company struct {
	UUID      string          `json:"uuid" gorm:"column:uuid;type:varchar(36)"`
	Title     string          `json:"title" gorm:"column:title;type:varchar(255)"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at;type:datetime"`
	OwnerUUID string          `json:"owner_uuid" gorm:"column:owner_uuid;type:varchar(36)"`
}

func (Company) TableName() string {
	return "company"
}
