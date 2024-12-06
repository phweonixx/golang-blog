package companies

import (
	"blogAPI/internal/helpers"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func SoftDeleteCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання UUID компанії
	UUID := vars["uuid"]

	// Перевірка на існування вказаної компанії
	exists, err := helpers.CheckExists(UUID, "company")
	if err != nil {
		http.Error(w, "Error checking company existence:", http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "Company with this UUID not found!", http.StatusNotFound)
		return
	}

	// Перевірка чи є користувач автором компанії
	var companyAuthorUUID string

	err = db.DBGorm.Model(&models.Company{}).
		Select("owner_uuid").
		Where("uuid = ?", UUID).
		Limit(1).
		Scan(&companyAuthorUUID).Error
	if err != nil {
		http.Error(w, "Error checking company", http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if companyAuthorUUID != middleware.User_UUID {
		http.Error(w, "You do not have the right to delete this category! You are not its author!", http.StatusForbidden)
		return
	}

	// Функція soft-видалення компанії
	err = SoftDeleteCompany(UUID)
	if err != nil {
		http.Error(w, "Error deleting company!", http.StatusNotFound)
		log.Println(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
