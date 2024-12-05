package articles

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// Функція створення статті
func CreateArticle(w http.ResponseWriter, r *http.Request) {
	// Отримання тіла запросу в JSON
	var article models.Article
	err := json.NewDecoder(r.Body).Decode(&article)
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

	// Перевірка чи має користувач компанію
	var companyUUID string
	err = database.DBGorm.Model(&models.Company{}).
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

	if article.Title == "" || article.Slug == "" {
		http.Error(w, "Title and Slug are required fields!", http.StatusBadRequest)
		return
	}

	// Задання часу створення статті
	timeNow := time.Now()
	article.CreatedAt = timeNow
	article.UpdatedAt = timeNow
	article.CompanyUUID = companyUUID
	article.UserUUID = middleware.User_UUID

	// Запрос для створення статті
	result := database.DBGorm.Create(&article)
	if result.Error != nil {
		http.Error(w, "Error creating article!", http.StatusInternalServerError)
		log.Println(result.Error)
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
				ObjectID: article.ID,
				Field:    field,
				Language: article.Language,
				Content:  content,
			}

			err := database.DBGorm.Create(&translation).Error
			if err != nil {
				log.Printf("Error creating translation for field %s: %v", field, err)
			}
		}
	}

	if len(article.RelatedArticlesID) != 0 {
		for _, relatedID := range article.RelatedArticlesID {
			relatedArticles := models.RelatedArticles{
				ParentArticleID:  article.ID,
				RelatedArticleID: relatedID,
			}

			err := database.DBGorm.Create(&relatedArticles).Error
			if err != nil {
				log.Println("Error creating relations for articles! The selected article may not exist.\n",
					err,
					"\nProblem with article with id", relatedID)
			}
		}
	}

	// Інформування про успішне створення статті
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Article created successfully",
		"article": article,
	})
}

func checkLanguage(lang string) error {
	validLanguages := map[string]bool{"en": true, "uk": true}
	if !validLanguages[lang] {
		return errors.New("invalid language: valid values are 'en' or 'uk'")
	}
	return nil
}
