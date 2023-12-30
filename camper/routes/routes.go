// routes/routes.go

package routes

import (
	"encoding/json"
	"fmt"
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
	http.HandleFunc("/ws", handleWebSocket)
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

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			return
		}

		// Print received message from client
		fmt.Printf("Received from client: %s\n", p)

		// Write "hello" back to the client
		if err := conn.WriteMessage(messageType, []byte("hello from server")); err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}
	}
}

// It worked for ws://localhost:8080/ws
