package models

import "time"

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

type FriendRequestInput struct {
	Username string `json:"username"`
}
