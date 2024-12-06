package categories

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID категорії
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting category id!", http.StatusBadRequest)
		log.Println("Invalid category ID:", err)
		return
	}

	// Перевірка на існування вказаної категорії
	exists, err := helpers.CheckExists(id, "category")
	if err != nil {
		http.Error(w, "Error checking category existance.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Category not found.", http.StatusNotFound)
		return
	}

	// Перевірка чи є користувач автором категорії
	var categoryAuthorUUID string

	err = db.DBGorm.Model(&models.Category{}).
		Select("user_uuid").
		Where("id = ?", id).
		Limit(1).
		Scan(&categoryAuthorUUID).Error
	if err != nil {
		http.Error(w, "Error checking category", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if categoryAuthorUUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to update this category! You are not its author!", http.StatusForbidden)
		return
	}

	// Отримання тіла запросу в JSON
	var category models.Category
	err = json.NewDecoder(r.Body).Decode(&category)
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

	// Задання часу оновлення категорії
	category.UpdatedAt = time.Now()

	// Запрос для оновлення категорії
	err = db.DBGorm.Model(&models.Category{}).
		Where("id = ?", id).
		Updates(models.Category{
			Language:  category.Language,
			Parent_id: category.Parent_id,
			UpdatedAt: category.UpdatedAt,
		}).Error
	if err != nil {
		http.Error(w, "Error updating category!", http.StatusInternalServerError)
		log.Println("Error updating category:", err)
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
				ObjectID: id,
				Field:    field,
				Language: category.Language,
				Content:  content,
			}

			var existingTranslation models.Translations
			err = db.DBGorm.Model(&models.Translations{}).
				Where("type = ? AND object_id = ? AND language = ? AND field = ?", translation.Type, translation.ObjectID, translation.Language, translation.Field).
				First(&existingTranslation).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = db.DBGorm.Create(&translation).Error
				if err != nil {
					http.Error(w, "Error creating translation!", http.StatusInternalServerError)
					log.Println(err)
					return
				}
			} else if err == nil {
				err = db.DBGorm.Model(&models.Translations{}).
					Where("type = ? AND object_id = ? AND language = ? AND field = ?", translation.Type, translation.ObjectID, translation.Language, translation.Field).
					Updates(map[string]interface{}{
						"content": content,
					}).Error

				if err != nil {
					http.Error(w, "Error updating translation!", http.StatusInternalServerError)
					log.Println(err)
					return
				}
			} else {
				http.Error(w, "Error checking translation existence!", http.StatusInternalServerError)
				log.Println(err)
				return
			}
		}
	}

	// Інформування про успішне оновлення категорії
	helpers.SendJSONResponse(w, http.StatusOK, "Category updated successfully", helpers.Empty{})
}
