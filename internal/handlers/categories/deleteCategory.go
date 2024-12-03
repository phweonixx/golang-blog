package categories

import (
	"blogAPI/internal/database"
	"blogAPI/internal/translations"
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

	exists, err := checkCategoryExists(id)
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

	err = database.DBGorm.Model(&Category{}).
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
	err = database.DBGorm.
		Unscoped().
		Where("id = ?", id).
		Delete(&Category{}).Error
	if err != nil {
		http.Error(w, "Error deleting the category!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Видалення перекладів
	err = database.DBGorm.
		Unscoped().
		Where("object_id = ? AND type = ?", id, "category").
		Delete(&translations.Translations{}).Error
	if err != nil {
		http.Error(w, "Error deleting translations for the category!", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	// Інформування про успішне видалення категорії
	w.WriteHeader(http.StatusNoContent)
}

func checkCategoryExists(id int) (bool, error) {
	var count int64
	err := database.DBGorm.Model(&Category{}).
		Where("id = ?", id).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
