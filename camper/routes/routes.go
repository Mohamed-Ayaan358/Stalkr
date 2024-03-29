// routes/routes.go

package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

type ApiResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
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

	response := ApiResponse{
		Message: "Success",
		Data:    websites,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response as JSON", http.StatusInternalServerError)
	}
}
func AddWebsite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

	// Respond with a JSON message
	response := ApiResponse{
		Message: "Website added successfully",
		Data:    website, // You can include additional data if needed
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
}

func DeleteWebsite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	var websiteName struct {
		Name string `json:"websiteName"`
	}

	err := json.NewDecoder(r.Body).Decode(&websiteName)
	if err != nil {
		http.Error(w, "Invalid website ID", http.StatusBadRequest)
		return
	}

	err = database.DeleteWebsite(websiteName.Name)
	if err != nil {
		http.Error(w, "Failed to delete website from the database", http.StatusInternalServerError)
		return
	}
	response := ApiResponse{
		Message: "Website added successfully",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to serialize JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse)
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
	fmt.Println("Hi entered sendmessagetoclient")

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
