package articles

import (
	"blogAPI/internal/config"
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

func ReadArticlesMy(w http.ResponseWriter, r *http.Request) {
	// Отримання параметрів пошуку
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	lang := r.URL.Query().Get("lang")
	categoryID := r.URL.Query().Get("category_id")

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

	cfg := config.New()

	// Перевірка введеного значення для мови
	err := checkLanguage(lang)
	if lang == "" {
		lang = cfg.Config.DefaultLang
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Вираховування параметру offset
	offset := (pageDefault - 1) * limitDefault

	// Функція отримання статей, що підходять по параметрам пошуку
	articles, err := getArticlesMyWithTranslation(limitDefault, offset, lang, categoryID)
	if err != nil {
		http.Error(w, "Error showing articles", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var total int64
	err = database.DBGorm.Model(&models.Article{}).
		Count(&total).Error
	if err != nil {
		http.Error(w, "Error counting articles", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := Response{
		Articles: articles,
		Total:    total,
		Page:     pageDefault,
		Limit:    limitDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getArticlesMyWithTranslation(limit, offset int, language, category_id string) ([]models.Article, error) {
	var articles []models.Article
	// Запрос для пошуку статей по введеним значенням
	query := database.DBGorm.Model(&models.Article{}).
		Select("id, category_id, company_uuid, language, created_at, updated_at, user_uuid")

	if category_id != "" {
		query = query.Where("category_id = ?", category_id)
	}
	query = query.
		Where("user_uuid = ?", middleware.User_UUID).
		Limit(limit).
		Offset(offset)

	err := query.Find(&articles).Error
	if err != nil {
		return nil, err
	}

	// Отримання перекладів для знайденої статті
	for idx, article := range articles {
		fields := []string{"title", "slug", "description", "seo_title", "seo_description"}
		for _, field := range fields {
			var content string
			err := database.DBGorm.Model(&models.Translations{}).
				Select("content").
				Where("type = ? AND object_id = ? AND language = ? AND field = ?", "article", article.ID, language, field).
				Scan(&content).Error
			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Println("Error getting translation:", err)
			}
			switch field {
			case "title":
				articles[idx].Title = content
			case "description":
				articles[idx].Description = content
			case "slug":
				articles[idx].Slug = content
			case "seo_title":
				articles[idx].SeoTitle = content
			case "seo_description":
				articles[idx].SeoDescription = content
			}
		}
		var relatedArticles []int
		err := database.DBGorm.Model(&models.RelatedArticles{}).
			Where("parent_article_id = ?", article.ID).
			Pluck("related_article_id", &relatedArticles).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("Error getting related articles:", err)
		}
		articles[idx].RelatedArticlesID = relatedArticles
	}

	// Повернення отриманого списку структур для виводу користувачу
	return articles, nil
}
