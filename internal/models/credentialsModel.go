package models

// Модель для прийому даних користувача при вході в акаунт
type Credentials struct {
	Username string `json:"username" gorm:"-"`
	Email    string `json:"email" gorm:"-"`
	Password string `json:"password" gorm:"-"`
}
