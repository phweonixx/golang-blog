package categories

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID категорії
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Error getting category id!", http.StatusBadRequest)
		log.Println(err)
		return
	}

	exists, err := helpers.CheckExists(id, "category")
	if err != nil {
		http.Error(w, "Error checking category existance.", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if !exists {
		http.Error(w, "Category not found.", http.StatusNotFound)
		return
	}

	// Перевірка чи є користувач автором категорії
	var categoryAuthorUUID string

	err = db.DBGorm.Model(&models.Category{}).
		Select("user_uuid").
		Where("id = ?", id).
		Limit(1).
		Scan(&categoryAuthorUUID).Error
	if err != nil {
		http.Error(w, "Error checking category", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if categoryAuthorUUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to delete this category! You are not its author!", http.StatusForbidden)
		return
	}

	// Видалення категорії
	err = db.DBGorm.
		Unscoped().
		Where("id = ?", id).
		Delete(&models.Category{}).Error
	if err != nil {
		http.Error(w, "Error deleting the category!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Видалення перекладів
	err = db.DBGorm.
		Unscoped().
		Where("object_id = ? AND type = ?", id, "category").
		Delete(&models.Translations{}).Error
	if err != nil {
		http.Error(w, "Error deleting translations for the category!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Інформування про успішне видалення категорії
	w.WriteHeader(http.StatusNoContent)
}
