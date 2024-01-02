// routes/routes.go

package routes

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	database "camper/db"
	"camper/models"

	"github.com/gorilla/websocket"
)

func SetupRoutes() {
	http.HandleFunc("/add", AddWebsite)
	http.HandleFunc("/delete", DeleteWebsite)
	http.HandleFunc("/websites", GetWebsites)
	http.HandleFunc("/ws", HandleWebSocket)
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

// ------------------------------------------

var (
	WsConn   *websocket.Conn
	WsMutex  sync.Mutex
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error upgrading to WebSocket: %v", err)
		return
	}
	defer conn.Close()

	WsMutex.Lock()
	WsConn = conn
	WsMutex.Unlock()

	for {
		messageType, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}
		if err := conn.WriteMessage(messageType, []byte("Server echo: ")); err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}
	}
}

func SendMessageToClient(message int) {
	WsMutex.Lock()
	defer WsMutex.Unlock()

	database.QueryInterval(message)
	websites, _ := database.GetWebsites()
	jsonified, _ := json.Marshal(websites)

	if WsConn != nil {
		err := WsConn.WriteMessage(websocket.TextMessage, jsonified)
		if err != nil {
			log.Printf("Error writing periodic message: %v", err)
		}
	}
}

// It worked for ws://localhost:8080/ws
