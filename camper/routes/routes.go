// routes/routes.go

package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	database "camper/db"
	"camper/models"
)

func SetupRoutes() {
	http.HandleFunc("/add", AddWebsite)
	http.HandleFunc("/delete", DeleteWebsite)
	http.HandleFunc("/websites", GetWebsites)
}

// Add one yourself to obtain all the websites in the database

func GetWebsites(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	websites, err := database.GetWebsites()
	if err != nil {
		http.Error(w, "Failed to get websites from the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(websites)
}

func AddWebsite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var website models.Website

	err := json.NewDecoder(r.Body).Decode(&website)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = database.AddWebsite(&website)
	if err != nil {
		http.Error(w, "Failed to add website to the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Website added successfully"))
}

func DeleteWebsite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	websiteID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid website ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteWebsite(websiteID)
	if err != nil {
		http.Error(w, "Failed to delete website from the database", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Website deleted successfully"))
}

// Put in a function that can be used to add or get methods  in ordee to derive a hash from a website
