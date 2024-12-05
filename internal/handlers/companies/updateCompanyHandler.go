package companies

import (
	"blogAPI/internal/database"
	"blogAPI/internal/models"
	"blogAPI/pkg/middleware"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func UpdateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Отримання ID категорії
	UUID := vars["uuid"]

	// Перевірка на існування вказаної компанії
	exists, err := checkCompanyExistsByUUID(UUID)
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

	err = database.DBGorm.Model(&models.Company{}).
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

	// Отримання тіла запросу в JSON
	var company models.Company
	err = json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Функція оновлення компанії
	err = UpdateCompany(&company, UUID)
	if err != nil {
		http.Error(w, "Error updating company", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Company updated successfully",
	})
}
