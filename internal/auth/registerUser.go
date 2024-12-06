package auth

import (
	"blogAPI/internal/models"
	"errors"
	"log"
	"regexp"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailOrUsernameUsed = errors.New("email or username is already used")
	ErrUsernameTooShort    = errors.New("username less than 3 characters")
	ErrPasswordTooShort    = errors.New("password less than 8 characters")
	ErrFirstNameTooShort   = errors.New("first name less than 2 characters")
	ErrLastNameTooShort    = errors.New("last name less than 2 characters")
	ErrEmailTooShort       = errors.New("email less than 6 characters")
	ErrEmailInvalidFormat  = errors.New("invalid email format")
)

func RegisterUser(user *models.User) error {
	// Перевірка існування акаунту
	exists, err := checkUserExists(user.Username, user.Email)
	if err != nil {
		return err
	}
	if exists {
		return ErrEmailOrUsernameUsed
	}

	err = ValidateUserInput(user)
	if err != nil {
		return err
	}

	// Хешування паролю
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.UUID = uuid.New().String()
	timeNow := time.Now()
	user.CreatedAt = timeNow
	user.UpdatedAt = timeNow

	result := db.DBGorm.Create(&user)
	if result.Error != nil {
		log.Println("Error creating user:", result.Error)
		return result.Error
	}

	return nil
}

func checkUserExists(username, email string) (bool, error) {
	var countOfSimilarUsers int64
	err := db.DBGorm.Model(&models.User{}).
		Where("username = ? OR email = ?", username, email).
		Count(&countOfSimilarUsers).Error
	if err != nil {
		log.Println("Error executing query:", err)
		return false, err
	}

	return countOfSimilarUsers > 0, nil
}

func ValidateUserInput(user *models.User) error {
	// Перевірка валідності введених значень
	if len(user.Username) < 3 {
		return ErrUsernameTooShort
	}
	if len(user.Password) < 8 {
		return ErrPasswordTooShort
	}
	if len(user.FirstName) < 2 {
		return ErrFirstNameTooShort
	}
	if len(user.LastName) < 2 {
		return ErrLastNameTooShort
	}
	if len(user.Email) < 6 {
		return ErrEmailTooShort
	}
	if !isValidEmail(user.Email) {
		return ErrEmailInvalidFormat
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}
