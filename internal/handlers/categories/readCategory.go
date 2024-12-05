package categories

import (
	"blogAPI/internal/config"
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ReadCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID категорії
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting category id!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	cfg := config.New()

	// Перевірка введеного значення для мови
	lang := r.URL.Query().Get("lang")

	err = checkLanguage(lang)
	if lang == "" {
		lang = cfg.Config.DefaultLang
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	exists, err := checkCategoryExists(id)
	if err != nil {
		http.Error(w, "Error checking category existance.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Category not found.", http.StatusNotFound)
		return
	}

	// Функція отримання категорії
	category, err := getCategoryWithTranslation(id, lang)
	if err != nil {
		http.Error(w, "Error processing translation!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategoryWithTranslation(id int, language string) (models.Category, error) {
	var category models.Category
	// Заповнення структури даними з бази даних
	err := database.DBGorm.First(&category, id).Error
	if err != nil {
		log.Println(err)
		return category, err
	}

	// Отримання перекладів
	fields := []string{"title", "slug", "seo_title", "seo_description"}
	for _, field := range fields {
		var content string

		result := database.DBGorm.Model(&models.Translations{}).
			Select("content").
			Where("type = ? AND object_id = ? AND language = ? AND field = ?", "category", id, language, field).
			Scan(&content)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				content = ""
			} else {
				log.Println("Error getting translation:", result.Error)
			}
		}

		// Заповнення структури отриманими перекладами
		switch field {
		case "title":
			category.Title = content
		case "slug":
			category.Slug = content
		case "seo_title":
			category.SeoTitle = content
		case "seo_description":
			category.SeoDescription = content
		}
	}

	// Повернення отриманої структури для виводу користувачу
	return category, nil
}
