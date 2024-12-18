package categories

import (
	"blogAPI/internal/database"
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

var db = database.New()

// Функція створення статті
func CreateCategory(w http.ResponseWriter, r *http.Request) {
	// Отримання тіла запросу в JSON
	var category models.Category
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}
	// Перевірка введеного значення для мови
	err = helpers.CheckLanguage(category.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Перевірка чи має користувач компанію
	var companyUUID string
	err = db.DBGorm.Model(&models.Company{}).
		Select("uuid").
		Where("owner_uuid = ?", middleware.User_UUID).
		Take(&companyUUID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "You don't have a company yet! Create it.", http.StatusForbidden)
		} else {
			http.Error(w, "Error retrieving company UUID!", http.StatusInternalServerError)
		}
		log.Println(err)
		return
	}

	if category.Title == "" || category.Slug == "" {
		http.Error(w, "Title and Slug are required fields!", http.StatusBadRequest)
		return
	}

	// Задання часу оновлення категорії
	timeNow := time.Now()
	category.CreatedAt = timeNow
	category.UpdatedAt = timeNow
	category.CompanyUUID = companyUUID
	category.User_uuid = middleware.User_UUID

	// Запрос для створення категорії
	result := db.DBGorm.Create(&category)
	if result.Error != nil {
		http.Error(w, "Error creating category!", http.StatusInternalServerError)
		log.Println(result.Error)
		return
	}

	// Оновлення перекладів
	translationType := "category"
	fields := []string{"title", "slug", "seo_title", "seo_description"}
	for _, field := range fields {
		var content string
		switch field {
		case "title":
			content = category.Title
		case "slug":
			content = category.Slug
		case "seo_title":
			content = category.SeoTitle
		case "seo_description":
			content = category.SeoDescription
		}

		if content != "" {
			translation := models.Translations{
				Type:     translationType,
				ObjectID: category.ID,
				Field:    field,
				Language: category.Language,
				Content:  content,
			}

			err := db.DBGorm.Create(&translation).Error
			if err != nil {
				log.Printf("Error creating translation for field %s: %v", field, err)
			}
		}
	}

	// Інформування про успішне створення категорії
	helpers.SendJSONResponse(w, http.StatusCreated, "Category created successfully", category)
}
