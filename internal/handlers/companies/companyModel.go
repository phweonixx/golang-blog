package companies

import (
	"time"

	"gorm.io/gorm"
)

// Модель для компанії
type Company struct {
	UUID      string          `json:"uuid" gorm:"column:uuid"`
	Title     string          `json:"title" gorm:"column:title"`
	CreatedAt time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at" gorm:"column:deleted_at"`
	OwnerUUID string          `json:"owner_uuid" gorm:"column:owner_uuid"`
}

func (Company) TableName() string {
	return "company"
}
