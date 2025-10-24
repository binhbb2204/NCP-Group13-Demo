package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// WebSocket clients
var clients = make(map[int]*websocket.Conn) // user_id -> connection
var clientsMux sync.Mutex

// Models
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Friend struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type FriendRequest struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	Username   string    `json:"username"`
	FriendID   int       `json:"friend_id"`
	FriendName string    `json:"friend_name"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Message struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	RecipientID int       `json:"recipient_id"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type WSMessage struct {
	Type        string      `json:"type"` // "message", "friend_request", "friend_accepted"
	Data        interface{} `json:"data"`
	RecipientID int         `json:"recipient_id,omitempty"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SendMessageRequest struct {
	RecipientID int    `json:"recipient_id"`
	Message     string `json:"message"`
}

type FriendRequestInput struct {
	Username string `json:"username"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func main() {
	// Database connection settings
	username := "root"
	password := "12342204"
	hostname := "127.0.0.1:3306"
	dbname := "chat_sys"

	// Data Source Name (DSN)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", username, password, hostname, dbname)

	// Open the connection
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("✅ Successfully connected to MySQL!")

	// Setup routes
	router := mux.NewRouter()

	// Auth routes
	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	router.HandleFunc("/api/login", loginHandler).Methods("POST")

	// Friend routes
	router.HandleFunc("/api/friends", getFriendsHandler).Methods("GET")
	router.HandleFunc("/api/friends/search", searchUsersHandler).Methods("GET")
	router.HandleFunc("/api/friends/request", sendFriendRequestHandler).Methods("POST")
	router.HandleFunc("/api/friends/requests", getFriendRequestsHandler).Methods("GET")
	router.HandleFunc("/api/friends/accept/{id}", acceptFriendRequestHandler).Methods("POST")
	router.HandleFunc("/api/friends/reject/{id}", rejectFriendRequestHandler).Methods("POST")

	// Message routes
	router.HandleFunc("/api/messages/{friendId}", getMessagesHandler).Methods("GET")
	router.HandleFunc("/api/messages", sendMessageHandler).Methods("POST")
	router.HandleFunc("/api/messages/unread", getUnreadCountHandler).Methods("GET")

	// WebSocket
	router.HandleFunc("/ws/{userId}", handleWebSocket)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	// Start server
	fmt.Println("🚀 Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", enableCORS(router)))
}

// Enable CORS
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Register handler
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		sendJSON(w, Response{Success: false, Message: "Username and password are required"}, http.StatusBadRequest)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error processing password"}, http.StatusInternalServerError)
		return
	}

	// Insert user
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, string(hashedPassword))
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Username already exists"}, http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()
	sendJSON(w, Response{
		Success: true,
		Message: "User registered successfully",
		Data:    map[string]interface{}{"user_id": userID, "username": req.Username},
	}, http.StatusCreated)
}

// Login handler
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	var user User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).
		Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid username or password"}, http.StatusUnauthorized)
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid username or password"}, http.StatusUnauthorized)
		return
	}

	sendJSON(w, Response{
		Success: true,
		Message: "Login successful",
		Data:    map[string]interface{}{"user_id": user.ID, "username": user.Username},
	}, http.StatusOK)
}

// Search users handler
func searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	userID := r.URL.Query().Get("user_id")

	if query == "" {
		sendJSON(w, Response{Success: false, Message: "Search query is required"}, http.StatusBadRequest)
		return
	}

	rows, err := db.Query(`
		SELECT id, username 
		FROM users 
		WHERE username LIKE ? AND id != ?
		LIMIT 10
	`, "%"+query+"%", userID)

	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error searching users"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			continue
		}
		users = append(users, map[string]interface{}{"id": id, "username": username})
	}

	sendJSON(w, Response{Success: true, Data: users}, http.StatusOK)
}

// Send friend request handler
func sendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	// Get friend ID by username
	var friendID int
	err := db.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&friendID)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "User not found"}, http.StatusNotFound)
		return
	}

	if friendID == req.UserID {
		sendJSON(w, Response{Success: false, Message: "Cannot add yourself as friend"}, http.StatusBadRequest)
		return
	}

	// Check if friendship already exists
	var exists int
	db.QueryRow(`
		SELECT COUNT(*) FROM friendships 
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)
	`, req.UserID, friendID, friendID, req.UserID).Scan(&exists)

	if exists > 0 {
		sendJSON(w, Response{Success: false, Message: "Friend request already exists"}, http.StatusConflict)
		return
	}

	// Create friend request
	_, err = db.Exec("INSERT INTO friendships (user_id, friend_id, status) VALUES (?, ?, 'pending')", req.UserID, friendID)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error sending friend request"}, http.StatusInternalServerError)
		return
	}

	// Notify via WebSocket
	notifyUser(friendID, WSMessage{
		Type: "friend_request",
		Data: map[string]interface{}{
			"user_id":  req.UserID,
			"username": req.Username,
		},
	})

	sendJSON(w, Response{Success: true, Message: "Friend request sent"}, http.StatusOK)
}

// Get friend requests handler
func getFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := db.Query(`
		SELECT f.id, f.user_id, u.username, f.created_at
		FROM friendships f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`, userID)

	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error fetching friend requests"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var username string
		var createdAt time.Time
		if err := rows.Scan(&id, &userID, &username, &createdAt); err != nil {
			continue
		}
		requests = append(requests, map[string]interface{}{
			"id":         id,
			"user_id":    userID,
			"username":   username,
			"created_at": createdAt,
		})
	}

	sendJSON(w, Response{Success: true, Data: requests}, http.StatusOK)
}

// Accept friend request handler
func acceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var userID, friendID int
	err := db.QueryRow("SELECT user_id, friend_id FROM friendships WHERE id = ?", requestID).Scan(&userID, &friendID)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Friend request not found"}, http.StatusNotFound)
		return
	}

	// Update status
	_, err = db.Exec("UPDATE friendships SET status = 'accepted' WHERE id = ?", requestID)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error accepting friend request"}, http.StatusInternalServerError)
		return
	}

	// Notify the requester
	notifyUser(userID, WSMessage{
		Type: "friend_accepted",
		Data: map[string]interface{}{
			"friend_id": friendID,
		},
	})

	sendJSON(w, Response{Success: true, Message: "Friend request accepted"}, http.StatusOK)
}

// Reject friend request handler
func rejectFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	_, err := db.Exec("UPDATE friendships SET status = 'rejected' WHERE id = ?", requestID)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error rejecting friend request"}, http.StatusInternalServerError)
		return
	}

	sendJSON(w, Response{Success: true, Message: "Friend request rejected"}, http.StatusOK)
}

// Get friends handler
func getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := db.Query(`
		SELECT DISTINCT u.id, u.username,
			(SELECT COUNT(*) FROM messages 
			 WHERE sender_id = u.id AND recipient_id = ? AND is_read = 0) as unread_count
		FROM users u
		INNER JOIN friendships f ON 
			(f.user_id = ? AND f.friend_id = u.id) OR 
			(f.friend_id = ? AND f.user_id = u.id)
		WHERE f.status = 'accepted' AND u.id != ?
		ORDER BY u.username
	`, userID, userID, userID, userID)

	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error fetching friends"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var friends []map[string]interface{}
	for rows.Next() {
		var id int
		var username string
		var unreadCount int
		if err := rows.Scan(&id, &username, &unreadCount); err != nil {
			continue
		}
		friends = append(friends, map[string]interface{}{
			"id":           id,
			"username":     username,
			"unread_count": unreadCount,
		})
	}

	sendJSON(w, Response{Success: true, Data: friends}, http.StatusOK)
}

// Get messages between two users
func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	friendID := vars["friendId"]
	userID := r.URL.Query().Get("user_id")

	rows, err := db.Query(`
		SELECT m.id, m.sender_id, u.username, m.recipient_id, m.message, m.is_read, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE (m.sender_id = ? AND m.recipient_id = ?) 
		   OR (m.sender_id = ? AND m.recipient_id = ?)
		ORDER BY m.created_at ASC
		LIMIT 100
	`, userID, friendID, friendID, userID)

	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error fetching messages"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	// Mark messages as read
	db.Exec("UPDATE messages SET is_read = 1 WHERE sender_id = ? AND recipient_id = ?", friendID, userID)

	sendJSON(w, Response{Success: true, Data: messages}, http.StatusOK)
}

// Send message handler
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SenderID    int    `json:"sender_id"`
		RecipientID int    `json:"recipient_id"`
		Message     string `json:"message"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSON(w, Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	if req.SenderID == 0 || req.RecipientID == 0 || req.Message == "" {
		sendJSON(w, Response{Success: false, Message: "All fields are required"}, http.StatusBadRequest)
		return
	}

	// Insert message
	result, err := db.Exec("INSERT INTO messages (sender_id, recipient_id, message) VALUES (?, ?, ?)",
		req.SenderID, req.RecipientID, req.Message)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error sending message"}, http.StatusInternalServerError)
		return
	}

	messageID, _ := result.LastInsertId()

	// Get complete message
	var msg Message
	err = db.QueryRow(`
		SELECT m.id, m.sender_id, u.username, m.recipient_id, m.message, m.is_read, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?
	`, messageID).Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt)

	if err == nil {
		// Send to recipient via WebSocket
		notifyUser(req.RecipientID, WSMessage{
			Type: "message",
			Data: msg,
		})
	}

	sendJSON(w, Response{Success: true, Data: msg}, http.StatusCreated)
}

// Get unread count
func getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM messages WHERE recipient_id = ? AND is_read = 0", userID).Scan(&count)
	if err != nil {
		sendJSON(w, Response{Success: false, Message: "Error fetching unread count"}, http.StatusInternalServerError)
		return
	}

	sendJSON(w, Response{Success: true, Data: map[string]interface{}{"count": count}}, http.StatusOK)
}

// WebSocket handler
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
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

	// Keep connection alive and handle pings
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

// Notify user via WebSocket
func notifyUser(userID int, msg WSMessage) {
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

// Helper function to send JSON response
func sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
