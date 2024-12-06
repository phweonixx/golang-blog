package categories

import (
	"blogAPI/internal/config"
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func SearchCategories(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	searchValue := r.URL.Query().Get("value")

	cfg := config.New()

	// Перевірка введеного значення для мови
	err := helpers.CheckLanguage(lang)
	if lang == "" {
		lang = cfg.Config.DefaultLang
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Перевірка поля value
	if searchValue == "" {
		http.Error(w, "Enter a valid value for search!", http.StatusBadRequest)
		return
	}

	// Запрос для пошуку
	var translationMatches []struct {
		ObjectID int     `gorm:"column:object_id"`
		Score    float64 `gorm:"column:score"`
	}

	err = db.DBGorm.Table("translations").
		Select("object_id, MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) AS score", searchValue).
		Where("MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)", searchValue).
		Having("score > ?", cfg.Config.Score).
		Scan(&translationMatches).Error

	if err != nil {
		http.Error(w, "Error searching a value!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var categoriesID []int
	for _, match := range translationMatches {
		categoriesID = append(categoriesID, match.ObjectID)
	}

	if len(categoriesID) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No categories found."})
		return
	}

	categories, err := getCategoriesById(categoriesID, lang)
	if err != nil {
		http.Error(w, "Error getting categories", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var total int64
	err = db.DBGorm.Model(&models.Category{}).
		Count(&total).Error
	if err != nil {
		http.Error(w, "Error counting categories", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := Response{
		Categories: categories,
		Total:      total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функція отримання категорії
func getCategoriesById(ID []int, language string) ([]models.Category, error) {
	var categories []models.Category

	// Получение категорий по ID
	err := db.DBGorm.Where("id IN ?", ID).Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// Получение переводов для категорий
	for i := range categories {
		var translations []struct {
			Field   string
			Content string
		}

		err := db.DBGorm.Table("translations").
			Select("field, content").
			Where("type = ? AND object_id = ? AND language = ?", "category", categories[i].ID, language).
			Find(&translations).Error
		if err != nil {
			log.Println("Error getting translations:", err)
			continue
		}

		// Заполнение динамических полей
		for _, translation := range translations {
			switch translation.Field {
			case "title":
				categories[i].Title = translation.Content
			case "slug":
				categories[i].Slug = translation.Content
			case "seo_title":
				categories[i].SeoTitle = translation.Content
			case "seo_description":
				categories[i].SeoDescription = translation.Content
			}
		}
	}

	return categories, nil
}
