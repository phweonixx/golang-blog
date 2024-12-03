package articles

import (
	"blogAPI/internal/database"
	"blogAPI/internal/translations"
	"blogAPI/pkg/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID статті
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting article id!", http.StatusBadRequest)
		log.Println(err)
		return
	}

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

	err = database.DBGorm.Model(&Article{}).
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

	// Видалення статті
	err = database.DBGorm.Where("id = ?", id).Delete(&Article{}).Error
	if err != nil {
		http.Error(w, "Error deleting the article!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Видалення перекладів
	err = database.DBGorm.Where("object_id = ? AND type = ?", id, "article").Delete(&translations.Translations{}).Error
	if err != nil {
		http.Error(w, "Error deleting translations for the article!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Інформування про успішне видалення статті
	w.WriteHeader(http.StatusNoContent)
}

func checkArticleExists(id int) (bool, error) {
	var count int64
	err := database.DBGorm.Model(&Article{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
