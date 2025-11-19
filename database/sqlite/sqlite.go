package sqlite

import (
	"database/sql"

	"DB-Presentation/models"
)

// GetMessagesSQLite fetches messages between two users from SQLite.
func GetMessagesSQLite(db *sql.DB, userID, friendID int) ([]models.Message, error) {
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
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message
	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
			continue
		}
		messages = append(messages, msg)
	}

	return messages, nil
}

// InsertMessageSQLite inserts a message into SQLite and returns the created message (with created_at filled).
func InsertMessageSQLite(db *sql.DB, senderID, recipientID int, message string) (models.Message, error) {
	result, err := db.Exec("INSERT INTO messages (sender_id, recipient_id, message) VALUES (?, ?, ?)", senderID, recipientID, message)
	if err != nil {
		return models.Message{}, err
	}

	messageID, _ := result.LastInsertId()

	var msg models.Message
	err = db.QueryRow(`
        SELECT m.id, m.sender_id, u.username, m.recipient_id, m.message, m.is_read, m.created_at
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        WHERE m.id = ?
    `, messageID).Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt)

	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}

// MarkMessagesReadSQLite marks messages as read in SQLite.
func MarkMessagesReadSQLite(db *sql.DB, senderID, recipientID int) error {
	_, err := db.Exec("UPDATE messages SET is_read = 1 WHERE sender_id = ? AND recipient_id = ?", senderID, recipientID)
	return err
}

// CountUnreadSQLite returns unread count for a recipient.
func CountUnreadSQLite(db *sql.DB, recipientID int) (int, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM messages WHERE recipient_id = ? AND is_read = 0", recipientID).Scan(&count)
	return count, err
}
