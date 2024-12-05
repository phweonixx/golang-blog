package models

import "time"

// Модель для категорії
type Category struct {
	ID          int       `json:"id" gorm:"column:id"`
	CompanyUUID string    `json:"company_uuid" gorm:"column:company_uuid;type:varchar(36)"`
	Language    string    `json:"language" gorm:"column:language;type:enum('en','uk')"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:datetime"`
	User_uuid   string    `json:"user_uuid" gorm:"column:user_uuid;type:varchar(36)"`
	Parent_id   int       `json:"parent_id" gorm:"column:parent_id"`

	// Динамічні поля
	Title          string `json:"title" gorm:"-"`
	Slug           string `json:"slug" gorm:"-"`
	SeoTitle       string `json:"seo_title" gorm:"-"`
	SeoDescription string `json:"seo_description" gorm:"-"`
}

func (Category) TableName() string {
	return "category"
}
