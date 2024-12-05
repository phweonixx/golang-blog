package companies

import (
	"blogAPI/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

func CreateCompanyHandler(w http.ResponseWriter, r *http.Request) {
	// Отримання тіла запросу в JSON
	var company models.Company
	err := json.NewDecoder(r.Body).Decode(&company)
	if err != nil {
		http.Error(w, "Invalid Input!", http.StatusBadRequest)
		log.Println("Error decoding request body:", err)
		return
	}

	// Функція створення компанії
	err = CreateCompany(&company)
	if err != nil {
		switch err {
		case ErrTitleTooShort:
			http.Error(w, "Title must be more than 3 characters!", http.StatusBadRequest)
		case ErrCompanyExist:
			http.Error(w, "You already have a company!", http.StatusConflict)
		default:
			http.Error(w, "Error creating company!", http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	// Інформування про успішне створення компанії
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Company created successfully",
		"company": company,
	})
}
