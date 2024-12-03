package auth

import (
	"blogAPI/internal/config"
	"blogAPI/internal/database"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Claims struct {
	Username string `json:"username"`
	UserUUID string `json:"user_uuid"`
	jwt.StandardClaims
}

func LoginUser(credentials *Credentials) (string, error) {
	var storedHash string
	var user_uuid string
	err := database.DBGorm.Model(&User{}).
		Select("password, uuid").
		Where("username = ? OR email = ?", credentials.Username, credentials.Email).
		Row().
		Scan(&storedHash, &user_uuid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("User not found.")
			return "", nil
		}
		log.Println("Error retrieving user data:", err)
		return "", err
	}

	// Перевірка на правильність паролю
	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(credentials.Password))
	if err != nil {
		return "", errors.New("invalid password")
	}

	// Створення токену
	expirationTime := time.Now().Add(config.JWTConfig.TokenLifetime)
	claims := &Claims{
		Username: credentials.Username,
		UserUUID: user_uuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Підпис токену
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.JWTConfig.JWTSecret)
	if err != nil {
		log.Println("Error signing token:", err)
		return "", err
	}

	return signedToken, nil
}
