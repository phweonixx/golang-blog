package articles

import (
	"blogAPI/internal/config"
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func SearchArticles(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	searchValue := r.URL.Query().Get("value")

	cfg := config.New()

	// Перевірка введеного значення для мови
	err := checkLanguage(lang)
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
	var results []struct {
		ObjectID int
		Score    float64
	}
	err = database.DBGorm.Raw(`
		SELECT object_id, MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE) AS score
		FROM translations
		WHERE MATCH(content) AGAINST(? IN NATURAL LANGUAGE MODE)
		HAVING score > ?;
	`, searchValue, searchValue, cfg.Config.Score).Scan(&results).Error
	if err != nil {
		http.Error(w, "Error searching a value!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	articlesID := make([]int, len(results))
	for i, res := range results {
		articlesID[i] = res.ObjectID
	}

	if len(articlesID) == 0 {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "No articles found."})
		return
	}

	articles, err := getArticlesById(articlesID, lang)
	if err != nil {
		http.Error(w, "Error getting articles", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var total int64
	err = database.DBGorm.Model(&models.Article{}).Count(&total).Error
	if err != nil {
		http.Error(w, "Error counting articles", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := Response{
		Articles: articles,
		Total:    total,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функція отримання статті
func getArticlesById(ID []int, language string) ([]models.Article, error) {
	var articles []models.Article

	err := database.DBGorm.Where("id IN ?", ID).Find(&articles).Error
	if err != nil {
		return nil, err
	}

	for i := range articles {
		article := &articles[i]

		var translations []struct {
			Field   string
			Content string
		}
		err := database.DBGorm.
			Table("translations").
			Select("field, content").
			Where("type = ? AND object_id = ? AND language = ?", "article", article.ID, language).
			Find(&translations).Error
		if err != nil {
			log.Println("Error fetching translations:", err)
			continue
		}

		for _, translation := range translations {
			switch translation.Field {
			case "title":
				article.Title = translation.Content
			case "description":
				article.Description = translation.Content
			case "slug":
				article.Slug = translation.Content
			case "seo_title":
				article.SeoTitle = translation.Content
			case "seo_description":
				article.SeoDescription = translation.Content
			}
		}

		var relatedArticles []int
		err = database.DBGorm.
			Table("related_articles").
			Select("related_article_id").
			Where("parent_article_id = ?", article.ID).
			Scan(&relatedArticles).Error
		if err != nil {
			log.Println("Error fetching related articles:", err)
			continue
		}
		article.RelatedArticlesID = relatedArticles
	}

	return articles, nil
}
