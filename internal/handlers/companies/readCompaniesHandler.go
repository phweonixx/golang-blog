package companies

import (
	"blogAPI/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	Companies []Company `json:"categories"`
	Total     int64     `json:"total"`
	Page      int       `json:"page"`
	Limit     int       `json:"limit"`
}

func ReadCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	limit := r.URL.Query().Get("limit")
	page := r.URL.Query().Get("page")
	UUID := r.URL.Query().Get("uuid")
	ownerUUID := r.URL.Query().Get("owner_uuid")

	limitDefault := 10
	pageDefault := 1

	// Перевірка валідності введених значень для ліміту
	if limit != "" {
		limitInt, err := strconv.Atoi(limit)
		if err != nil || limitInt <= 0 {
			http.Error(w, "Enter a correct value for limit!", http.StatusBadRequest)
			log.Println(err)
			return
		}
		limitDefault = limitInt
	}

	// Перевірка валідності введених значень для сторінок
	if page != "" {
		pageInt, err := strconv.Atoi(page)
		if err != nil || pageInt <= 0 {
			http.Error(w, "Enter a correct value for page!", http.StatusBadRequest)
			log.Println(err)
			return
		}
		pageDefault = pageInt
	}

	// Вираховування параметру offset
	offset := (pageDefault - 1) * limitDefault

	// Функція отримання компаній, що підходять по параметрам пошуку
	companies, err := GetCompanies(limitDefault, offset, UUID, ownerUUID)
	if err != nil {
		http.Error(w, "Error showing companies", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	var total int64
	err = database.DBGorm.Model(&Company{}).
		Count(&total).Error
	if err != nil {
		http.Error(w, "Error counting companies", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	response := Response{
		Companies: companies,
		Total:     total,
		Page:      pageDefault,
		Limit:     limitDefault,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
