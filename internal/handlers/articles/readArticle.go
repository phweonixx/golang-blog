package articles

import (
	"blogAPI/internal/config"
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func ReadArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID статті
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting article id!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	cfg := config.New()

	// Перевірка введеного значення для мови
	lang := r.URL.Query().Get("lang")
	err = helpers.CheckLanguage(lang)
	if lang == "" {
		lang = cfg.Config.DefaultLang
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Перевірка на існування вказаної статті
	exists, err := helpers.CheckExists(id, "article")
	if err != nil {
		http.Error(w, "Error checking article existance.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Article not found.", http.StatusNotFound)
		return
	}

	// Функція отримання статті
	article, err := getArticleWithTranslation(id, lang)
	if err != nil {
		http.Error(w, "Error processing translation!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	helpers.SendJSONResponse(w, http.StatusOK, "Article found", article)
}

func getArticleWithTranslation(id int, language string) (models.Article, error) {
	var article models.Article
	// Заповнення структури даними з бази даних
	err := db.DBGorm.First(&article, id).Error
	if err != nil {
		log.Println(err)
		return article, err
	}

	// Отримання перекладів
	fields := []string{"title", "slug", "description", "seo_title", "seo_description"}
	for _, field := range fields {
		var content string

		result := db.DBGorm.Model(&models.Translations{}).
			Select("content").
			Where("type = ? AND object_id = ? AND language = ? AND field = ?", "article", id, language, field).
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
			article.Title = content
		case "description":
			article.Description = content
		case "slug":
			article.Slug = content
		case "seo_title":
			article.SeoTitle = content
		case "seo_description":
			article.SeoDescription = content
		}
	}

	var relatedArticles []int
	err = db.DBGorm.Model(&models.RelatedArticles{}).
		Select("related_article_id").
		Where("parent_article_id = ?", id).
		Pluck("related_article_id", &relatedArticles).Error

	if err != nil {
		log.Println("Error getting related articles.", err)
		return article, err
	}
	article.RelatedArticlesID = relatedArticles

	// Повернення отриманої структури для виводу користувачу
	return article, nil
}
