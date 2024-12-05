package articles

import (
	"blogAPI/internal/database"
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

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID статті
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting article id!", http.StatusBadRequest)
		log.Println("Invalid article ID:", err)
		return
	}

	// Перевірка на існування вказаної статті
	exists, err := checkArticleExists(id)
	if err != nil {
		http.Error(w, "Error checking article existance.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Article not found.", http.StatusNotFound)
		return
	}

	// Перевірка чи є користувач автором статті
	var articleAuthorUUID string

	err = database.DBGorm.Model(&models.Article{}).
		Select("user_uuid").
		Where("id = ?", id).
		Limit(1).
		Scan(&articleAuthorUUID).Error
	if err != nil {
		http.Error(w, "Error checking article", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if articleAuthorUUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to delete this article! You are not its author!", http.StatusForbidden)
		return
	}

	// Отримання тіла запросу в JSON
	var article models.Article
	err = json.NewDecoder(r.Body).Decode(&article)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}
	// Перевірка введеного значення для мови
	err = checkLanguage(article.Language)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Задання часу оновлення статті
	article.UpdatedAt = time.Now()

	// Запрос для оновлення статті
	err = database.DBGorm.Model(&models.Article{}).
		Where("id = ?", id).
		Updates(models.Article{
			CategoryID: article.CategoryID,
			Language:   article.Language,
			UpdatedAt:  article.UpdatedAt,
		}).Error
	if err != nil {
		http.Error(w, "Error updating article!", http.StatusInternalServerError)
		log.Println("Error updating article:", err)
		return
	}

	// Оновлення перекладів
	translationType := "article"

	fields := []string{"title", "slug", "description", "seo_title", "seo_description"}
	for _, field := range fields {
		var content string
		switch field {
		case "title":
			content = article.Title
		case "description":
			content = article.Description
		case "slug":
			content = article.Slug
		case "seo_title":
			content = article.SeoTitle
		case "seo_description":
			content = article.SeoDescription
		}

		if content != "" {
			translation := models.Translations{
				Type:     translationType,
				ObjectID: id,
				Field:    field,
				Language: article.Language,
				Content:  content,
			}

			var existingTranslation models.Translations
			err = database.DBGorm.Model(&models.Translations{}).
				Where("type = ? AND object_id = ? AND language = ? AND field = ?", translation.Type, translation.ObjectID, translation.Language, translation.Field).
				First(&existingTranslation).Error

			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = database.DBGorm.Create(&translation).Error
				if err != nil {
					http.Error(w, "Error creating translation!", http.StatusInternalServerError)
					log.Println(err)
					return
				}
			} else if err == nil {
				err = database.DBGorm.Model(&models.Translations{}).
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

	if len(article.RelatedArticlesID) != 0 {
		for _, relatedID := range article.RelatedArticlesID {
			var existing models.RelatedArticles
			err := database.DBGorm.Where("parent_article_id = ? AND related_article_id = ?", id, relatedID).First(&existing).Error

			if err != nil && err != gorm.ErrRecordNotFound {
				http.Error(w, "Error checking relation existence!", http.StatusInternalServerError)
				log.Println(err)
				return
			}

			if err == gorm.ErrRecordNotFound {
				err = database.DBGorm.Create(&models.RelatedArticles{
					ParentArticleID:  id,
					RelatedArticleID: relatedID,
				}).Error

				if err != nil {
					http.Error(w, "Error creating relation for article!", http.StatusInternalServerError)
					log.Println(err)
					return
				}
			} else {
				log.Println("Relation already exists, skipping creation for parent_article_id:", id, "and related_article_id:", relatedID)
			}
		}
	}

	// Інформування про успішне оновлення статті
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Article updated successfully",
	})
}
