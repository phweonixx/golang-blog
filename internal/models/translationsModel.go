package models

// Модель для перекладів
type Translations struct {
	ID       int    `json:"id" gorm:"column:id"`
	Type     string `json:"type" gorm:"column:type;type:enum('article','category')"`
	ObjectID int    `json:"object_id" gorm:"column:object_id"`
	Field    string `json:"field" gorm:"column:field;varchar(255)"`
	Language string `json:"language" gorm:"column:language;type:enum('en','uk')"`
	Content  string `json:"content" gorm:"column:content;type:longtext"`
}
