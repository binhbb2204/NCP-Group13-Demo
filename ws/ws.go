package ws

import (
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"DB-Presentation/models"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[int]*websocket.Conn)
var clientsMux sync.Mutex

// HandleWebSocket upgrades and manages a user's websocket connection.
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDStr := vars["userId"]
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	clientsMux.Lock()
	clients[userID] = conn
	clientsMux.Unlock()

	log.Printf("User %d connected. Total clients: %d", userID, len(clients))

	defer func() {
		clientsMux.Lock()
		delete(clients, userID)
		clientsMux.Unlock()
		conn.Close()
		log.Printf("User %d disconnected. Total clients: %d", userID, len(clients))
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// NotifyUser sends a WSMessage to a connected user (if present)
func NotifyUser(userID int, msg models.WSMessage) {
	clientsMux.Lock()
	defer clientsMux.Unlock()

	if conn, ok := clients[userID]; ok {
		err := conn.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending to user %d: %v", userID, err)
			conn.Close()
			delete(clients, userID)
		}
	}
}
