package articles

import "time"

// Модель для категорії
type Article struct {
	ID          int       `json:"id" gorm:"column:id;primaryKey"`
	CategoryID  int       `json:"category_id" gorm:"column:category_id"`
	CompanyUUID string    `json:"company_uuid" gorm:"column:company_uuid"`
	Language    string    `json:"language" gorm:"column:language"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
	UserUUID    string    `json:"user_uuid" gorm:"column:user_uuid"`

	// Динамічні поля
	Title          string `json:"title" gorm:"-"`
	Description    string `json:"description" gorm:"-"`
	Slug           string `json:"slug" gorm:"-"`
	SeoTitle       string `json:"seo_title" gorm:"-"`
	SeoDescription string `json:"seo_description" gorm:"-"`

	// ID схожих статей
	RelatedArticlesID []int `json:"related_articles_id" gorm:"-"`
}

func (Article) TableName() string {
	return "article"
}
