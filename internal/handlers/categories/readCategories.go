package categories

import (
	"blogAPI/internal/config"
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Response struct {
	Categories []models.Category `json:"categories"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	Limit      int               `json:"limit"`
}

func ReadCategories(w http.ResponseWriter, r *http.Request) {
	// Отримання параметрів пошуку
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	lang := r.URL.Query().Get("lang")
	companyUUID := r.URL.Query().Get("company_uuid")
	userUUID := r.URL.Query().Get("user_uuid")

	limitDefault := 10
	pageDefault := 1

	// Перевірка валідності введених значень для ліміту
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt <= 0 {
			http.Error(w, "Enter a correct value for limit!", http.StatusBadRequest)
			log.Println(err)
			return
		}
		limitDefault = limitInt
	}

	// Перевірка валідності введених значень для сторінок
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			http.Error(w, "Enter a correct value for page!", http.StatusBadRequest)
			log.Println(err)
			return
		}
		pageDefault = pageInt
	}

	if userUUID != "" {
		if !helpers.IsValidUUID(userUUID) {
			http.Error(w, "Enter a correct value for user UUID!", http.StatusBadRequest)
			return
		}
	}
	if companyUUID != "" {
		if !helpers.IsValidUUID(companyUUID) {
			http.Error(w, "Enter a correct value for company UUID!", http.StatusBadRequest)
			return
		}
	}

	cfg := config.New()

	// Перевірка введеного значення для мови
	err := helpers.CheckLanguage(lang)
	if lang == "" {
		lang = cfg.Config.DefaultLang
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вираховування параметру offset
	offset := (pageDefault - 1) * limitDefault

	// Функція отримання категорій, що підходять по параметрам пошуку
	categories, err := getCategoriesWithTranslation(limitDefault, offset, lang, companyUUID, userUUID)
	if err != nil {
		http.Error(w, "Error showing categories", http.StatusInternalServerError)
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
		Page:       pageDefault,
		Limit:      limitDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getCategoriesWithTranslation(limit, offset int, language, company_uuid, user_uuid string) ([]models.Category, error) {
	var categories []models.Category
	// Запрос для пошуку статей по введеним значенням

	query := db.DBGorm.Model(&models.Category{}).
		Select("id, company_uuid, language, created_at, updated_at, user_uuid, parent_id")

	if company_uuid != "" {
		query = query.Where("company_uuid = ?", company_uuid)
	}
	if user_uuid != "" {
		query = query.Where("user_uuid = ?", user_uuid)
	}
	query = query.Limit(limit).Offset(offset)

	err := query.Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// Отримання перекладів для знайденої категорії
	for idx, category := range categories {
		fields := []string{"title", "slug", "seo_title", "seo_description"}
		for _, field := range fields {
			var content string
			err := db.DBGorm.Model(&models.Translations{}).
				Select("content").
				Where("type = ? AND object_id = ? AND language = ? AND field = ?", "category", category.ID, language, field).
				Scan(&content).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Error getting translation:", err)
			}
			switch field {
			case "title":
				categories[idx].Title = content
			case "slug":
				categories[idx].Slug = content
			case "seo_title":
				categories[idx].SeoTitle = content
			case "seo_description":
				categories[idx].SeoDescription = content
			}
		}
	}

	// Повернення отриманого списку структур для виводу користувачу
	return categories, nil
}
