package mongo

import (
	"context"
	"database/sql"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"DB-Presentation/models"
)

// Connect connects to MongoDB using the provided URI.
func Connect(uri string) (*mongodriver.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongodriver.Connect(ctx, clientOpts)
	if err != nil {
		return nil, err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	if err := client.Ping(ctx2, nil); err != nil {
		return nil, err
	}

	return client, nil
}

// GetMessages returns messages between two users ordered by created_at ascending.
func GetMessages(ctx context.Context, client *mongodriver.Client, userID, friendID int) ([]models.Message, error) {
	coll := client.Database("chat").Collection("messages")
	filter := bson.M{"$or": []interface{}{
		bson.M{"sender_id": userID, "recipient_id": friendID},
		bson.M{"sender_id": friendID, "recipient_id": userID},
	}}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}).SetLimit(100)

	cur, err := coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var messages []models.Message
	for cur.Next(ctx) {
		var doc bson.M
		if err := cur.Decode(&doc); err != nil {
			continue
		}

		var msg models.Message
		// map fields robustly
		if v, ok := doc["sender_id"]; ok {
			switch t := v.(type) {
			case int32:
				msg.SenderID = int(t)
			case int64:
				msg.SenderID = int(t)
			case int:
				msg.SenderID = t
			case float64:
				msg.SenderID = int(t)
			}
		}
		if v, ok := doc["recipient_id"]; ok {
			switch t := v.(type) {
			case int32:
				msg.RecipientID = int(t)
			case int64:
				msg.RecipientID = int(t)
			case int:
				msg.RecipientID = t
			case float64:
				msg.RecipientID = int(t)
			}
		}
		if v, ok := doc["sender_name"]; ok {
			if s, ok := v.(string); ok {
				msg.SenderName = s
			}
		}
		if v, ok := doc["message"]; ok {
			if s, ok := v.(string); ok {
				msg.Message = s
			}
		}
		if v, ok := doc["is_read"]; ok {
			if b, ok := v.(bool); ok {
				msg.IsRead = b
			}
		}
		if v, ok := doc["created_at"]; ok {
			switch t := v.(type) {
			case primitive.DateTime:
				msg.CreatedAt = t.Time().UTC()
			case time.Time:
				msg.CreatedAt = t.UTC()
			}
		}

		messages = append(messages, msg)
	}

	return messages, nil
}

// InsertMessage inserts a message document into Mongo.
func InsertMessage(ctx context.Context, client *mongodriver.Client, msg models.Message) error {
	coll := client.Database("chat").Collection("messages")
	// store created_at as UTC DateTime
	created := primitive.NewDateTimeFromTime(msg.CreatedAt.UTC())
	_, err := coll.InsertOne(ctx, bson.M{
		"sender_id":    msg.SenderID,
		"sender_name":  msg.SenderName,
		"recipient_id": msg.RecipientID,
		"message":      msg.Message,
		"is_read":      msg.IsRead,
		"created_at":   created,
	})
	return err
}

// MarkMessagesRead marks messages sent by senderID to recipientID as read.
func MarkMessagesRead(ctx context.Context, client *mongodriver.Client, senderID, recipientID int) error {
	coll := client.Database("chat").Collection("messages")
	_, err := coll.UpdateMany(ctx, bson.M{"sender_id": senderID, "recipient_id": recipientID}, bson.M{"$set": bson.M{"is_read": true}})
	return err
}

// CountUnread returns number of unread messages for a recipient.
func CountUnread(ctx context.Context, client *mongodriver.Client, recipientID int) (int64, error) {
	coll := client.Database("chat").Collection("messages")
	cnt, err := coll.CountDocuments(ctx, bson.M{"recipient_id": recipientID, "is_read": false})
	return cnt, err
}

// MigrateFromSQLite copies messages from SQLite into Mongo. It does not delete SQLite rows.
func MigrateFromSQLite(ctx context.Context, client *mongodriver.Client, sqlDB *sql.DB) error {
	rows, err := sqlDB.Query(`
        SELECT m.id, m.sender_id, u.username, m.recipient_id, m.message, m.is_read, m.created_at
        FROM messages m
        JOIN users u ON m.sender_id = u.id
        ORDER BY m.created_at ASC
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	coll := client.Database("chat").Collection("messages")

	for rows.Next() {
		var msg models.Message
		if err := rows.Scan(&msg.ID, &msg.SenderID, &msg.SenderName, &msg.RecipientID, &msg.Message, &msg.IsRead, &msg.CreatedAt); err != nil {
			continue
		}
		// Insert into Mongo; ignore errors to be resilient
		created := primitive.NewDateTimeFromTime(msg.CreatedAt.UTC())
		_, _ = coll.InsertOne(ctx, bson.M{
			"sender_id":    msg.SenderID,
			"sender_name":  msg.SenderName,
			"recipient_id": msg.RecipientID,
			"message":      msg.Message,
			"is_read":      msg.IsRead,
			"created_at":   created,
		})
	}

	return nil
}
