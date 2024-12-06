package articles

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
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

	// Перевірка чи є користувач автором статті
	var articleAuthorUUID string

	err = db.DBGorm.Model(&models.Article{}).
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
	err = db.DBGorm.Where("id = ?", id).Delete(&models.Article{}).Error
	if err != nil {
		http.Error(w, "Error deleting the article!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Видалення перекладів
	err = db.DBGorm.Where("object_id = ? AND type = ?", id, "article").Delete(&models.Translations{}).Error
	if err != nil {
		http.Error(w, "Error deleting translations for the article!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Інформування про успішне видалення статті
	w.WriteHeader(http.StatusNoContent)
}
