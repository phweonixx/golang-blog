package auth

// Модель для прийому даних користувача при вході в акаунт
type Credentials struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
