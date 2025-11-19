package models

import "time"

type Message struct {
	ID          int       `json:"id"`
	SenderID    int       `json:"sender_id"`
	SenderName  string    `json:"sender_name"`
	RecipientID int       `json:"recipient_id"`
	Message     string    `json:"message"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

type SendMessageRequest struct {
	RecipientID int    `json:"recipient_id"`
	Message     string `json:"message"`
}
