package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	"DB-Presentation/models"
	"DB-Presentation/utils"
	"DB-Presentation/ws"

	dbmongo "DB-Presentation/database/mongo"
	dbsqlite "DB-Presentation/database/sqlite"
)

var dbase *sql.DB
var mClient *mongodriver.Client

// RegisterRoutes registers all HTTP routes with the provided router and DB handle.
func RegisterRoutes(router *mux.Router, db *sql.DB, mc *mongodriver.Client) {
	dbase = db
	mClient = mc

	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	router.HandleFunc("/api/login", loginHandler).Methods("POST")

	router.HandleFunc("/api/friends", getFriendsHandler).Methods("GET")
	router.HandleFunc("/api/friends/search", searchUsersHandler).Methods("GET")
	router.HandleFunc("/api/friends/request", sendFriendRequestHandler).Methods("POST")
	router.HandleFunc("/api/friends/requests", getFriendRequestsHandler).Methods("GET")
	router.HandleFunc("/api/friends/accept/{id}", acceptFriendRequestHandler).Methods("POST")
	router.HandleFunc("/api/friends/reject/{id}", rejectFriendRequestHandler).Methods("POST")

	router.HandleFunc("/api/messages/{friendId}", getMessagesHandler).Methods("GET")
	router.HandleFunc("/api/messages", sendMessageHandler).Methods("POST")
	router.HandleFunc("/api/messages/unread", getUnreadCountHandler).Methods("GET")
}

// registerHandler registers a new user
func registerHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		utils.SendJSON(w, models.Response{Success: false, Message: "Username and password are required"}, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error processing password"}, http.StatusInternalServerError)
		return
	}

	result, err := dbase.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, string(hashedPassword))
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Username already exists"}, http.StatusConflict)
		return
	}

	userID, _ := result.LastInsertId()
	utils.SendJSON(w, models.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    map[string]interface{}{"user_id": userID, "username": req.Username},
	}, http.StatusCreated)
}

// loginHandler authenticates a user
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	var id int
	var username, password string
	err := dbase.QueryRow("SELECT id, username, password FROM users WHERE username = ?", req.Username).Scan(&id, &username, &password)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid username or password"}, http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid username or password"}, http.StatusUnauthorized)
		return
	}

	utils.SendJSON(w, models.Response{Success: true, Message: "Login successful", Data: map[string]interface{}{"user_id": id, "username": username}}, http.StatusOK)
}

// searchUsersHandler searches users by query param 'q' and excludes user_id
func searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	userID := r.URL.Query().Get("user_id")

	if query == "" {
		utils.SendJSON(w, models.Response{Success: false, Message: "Search query is required"}, http.StatusBadRequest)
		return
	}

	rows, err := dbase.Query(`
		SELECT id, username 
		FROM users 
		WHERE username LIKE ? AND id != ?
		LIMIT 10
	`, "%"+query+"%", userID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error searching users"}, http.StatusInternalServerError)
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

	utils.SendJSON(w, models.Response{Success: true, Data: users}, http.StatusOK)
}

// sendFriendRequestHandler creates a pending friendship
func sendFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	var friendID int
	err := dbase.QueryRow("SELECT id FROM users WHERE username = ?", req.Username).Scan(&friendID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "User not found"}, http.StatusNotFound)
		return
	}

	if friendID == req.UserID {
		utils.SendJSON(w, models.Response{Success: false, Message: "Cannot add yourself as friend"}, http.StatusBadRequest)
		return
	}

	var exists int
	dbase.QueryRow(`
		SELECT COUNT(*) FROM friendships 
		WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)
	`, req.UserID, friendID, friendID, req.UserID).Scan(&exists)

	if exists > 0 {
		utils.SendJSON(w, models.Response{Success: false, Message: "Friend request already exists"}, http.StatusConflict)
		return
	}

	_, err = dbase.Exec("INSERT INTO friendships (user_id, friend_id, status) VALUES (?, ?, 'pending')", req.UserID, friendID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error sending friend request"}, http.StatusInternalServerError)
		return
	}

	ws.NotifyUser(friendID, models.WSMessage{Type: "friend_request", Data: map[string]interface{}{"user_id": req.UserID, "username": req.Username}})

	utils.SendJSON(w, models.Response{Success: true, Message: "Friend request sent"}, http.StatusOK)
}

// getFriendRequestsHandler returns pending requests for a user
func getFriendRequestsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := dbase.Query(`
		SELECT f.id, f.user_id, u.username, f.created_at
		FROM friendships f
		JOIN users u ON f.user_id = u.id
		WHERE f.friend_id = ? AND f.status = 'pending'
		ORDER BY f.created_at DESC
	`, userID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error fetching friend requests"}, http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var requests []map[string]interface{}
	for rows.Next() {
		var id, userIDint int
		var username string
		var createdAt time.Time
		if err := rows.Scan(&id, &userIDint, &username, &createdAt); err != nil {
			continue
		}
		requests = append(requests, map[string]interface{}{"id": id, "user_id": userIDint, "username": username, "created_at": createdAt})
	}

	utils.SendJSON(w, models.Response{Success: true, Data: requests}, http.StatusOK)
}

// acceptFriendRequestHandler accepts a pending friendship
func acceptFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	var userID, friendID int
	err := dbase.QueryRow("SELECT user_id, friend_id FROM friendships WHERE id = ?", requestID).Scan(&userID, &friendID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Friend request not found"}, http.StatusNotFound)
		return
	}

	_, err = dbase.Exec("UPDATE friendships SET status = 'accepted' WHERE id = ?", requestID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error accepting friend request"}, http.StatusInternalServerError)
		return
	}

	ws.NotifyUser(userID, models.WSMessage{Type: "friend_accepted", Data: map[string]interface{}{"friend_id": friendID}})
	utils.SendJSON(w, models.Response{Success: true, Message: "Friend request accepted"}, http.StatusOK)
}

// rejectFriendRequestHandler rejects a pending friendship
func rejectFriendRequestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestID := vars["id"]

	_, err := dbase.Exec("UPDATE friendships SET status = 'rejected' WHERE id = ?", requestID)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error rejecting friend request"}, http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, models.Response{Success: true, Message: "Friend request rejected"}, http.StatusOK)
}

// getFriendsHandler returns accepted friends with unread counts
func getFriendsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	rows, err := dbase.Query(`
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
		utils.SendJSON(w, models.Response{Success: false, Message: "Error fetching friends"}, http.StatusInternalServerError)
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
		friends = append(friends, map[string]interface{}{"id": id, "username": username, "unread_count": unreadCount})
	}

	utils.SendJSON(w, models.Response{Success: true, Data: friends}, http.StatusOK)
}

// getMessagesHandler fetches messages between two users and marks as read
func getMessagesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	friendID := vars["friendId"]
	userID := r.URL.Query().Get("user_id")
	// Primary: use Mongo adapter if client available
	if mClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		msgs, err := dbmongo.GetMessages(ctx, mClient, toInt(userID), toInt(friendID))
		if err == nil {
			_ = dbmongo.MarkMessagesRead(ctx, mClient, toInt(friendID), toInt(userID))
			utils.SendJSON(w, models.Response{Success: true, Data: msgs}, http.StatusOK)
			return
		}
	}

	// Fallback to SQLite
	msgs, err := dbsqlite.GetMessagesSQLite(dbase, toInt(userID), toInt(friendID))
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error fetching messages"}, http.StatusInternalServerError)
		return
	}

	_ = dbsqlite.MarkMessagesReadSQLite(dbase, toInt(friendID), toInt(userID))
	utils.SendJSON(w, models.Response{Success: true, Data: msgs}, http.StatusOK)
}

// sendMessageHandler inserts a message and notifies recipient via WS
func sendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		SenderID    int    `json:"sender_id"`
		RecipientID int    `json:"recipient_id"`
		Message     string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Invalid request"}, http.StatusBadRequest)
		return
	}

	if req.SenderID == 0 || req.RecipientID == 0 || req.Message == "" {
		utils.SendJSON(w, models.Response{Success: false, Message: "All fields are required"}, http.StatusBadRequest)
		return
	}

	result, err := dbase.Exec("INSERT INTO messages (sender_id, recipient_id, message) VALUES (?, ?, ?)", req.SenderID, req.RecipientID, req.Message)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error sending message"}, http.StatusInternalServerError)
		return
	}

	messageID, _ := result.LastInsertId()

	var msg models.Message
	err = dbase.QueryRow(`
		SELECT m.id, m.sender_id, u.username, m.recipient_id, m.message, m.is_read, m.created_at
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.id = ?
	`, messageID).Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt)

	if err == nil {
		// store in Mongo if available (Mongo is primary for messages)
		if mClient != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			// Normalize CreatedAt to UTC before inserting and notifying so clients
			// always receive an ISO timestamp with timezone information.
			msg.CreatedAt = msg.CreatedAt.UTC()
			_ = dbmongo.InsertMessage(ctx, mClient, msg)
		}

		ws.NotifyUser(req.RecipientID, models.WSMessage{Type: "message", Data: msg})
	}

	utils.SendJSON(w, models.Response{Success: true, Data: msg}, http.StatusCreated)
}

// getUnreadCountHandler returns total unread messages for a user
func getUnreadCountHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")

	// Prefer Mongo if available
	if mClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		cnt, err := dbmongo.CountUnread(ctx, mClient, toInt(userID))
		if err != nil {
			utils.SendJSON(w, models.Response{Success: false, Message: "Error fetching unread count"}, http.StatusInternalServerError)
			return
		}
		utils.SendJSON(w, models.Response{Success: true, Data: map[string]interface{}{"count": cnt}}, http.StatusOK)
		return
	}

	var count int
	err := dbase.QueryRow("SELECT COUNT(*) FROM messages WHERE recipient_id = ? AND is_read = 0", userID).Scan(&count)
	if err != nil {
		utils.SendJSON(w, models.Response{Success: false, Message: "Error fetching unread count"}, http.StatusInternalServerError)
		return
	}

	utils.SendJSON(w, models.Response{Success: true, Data: map[string]interface{}{"count": count}}, http.StatusOK)
}

// toInt converts a numeric string to int, returns 0 on error.
func toInt(s string) int {
	var i int
	_, err := fmt.Sscan(s, &i)
	if err != nil {
		return 0
	}
	return i
}

//
