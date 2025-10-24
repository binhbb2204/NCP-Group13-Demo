# Private Chat System - Go Backend + Frontend

A real-time private messaging application with friend system, built with Go (backend) and vanilla JavaScript (frontend).

## ✨ Features

- ✅ User Registration & Login
- ✅ Password Hashing (bcrypt)
- ✅ Friend System (Send/Accept/Reject friend requests)
- ✅ Private one-on-one messaging
- ✅ Real-time messaging with WebSocket
- ✅ Unread message indicators
- ✅ Message history
- ✅ Modern responsive UI
- ✅ MySQL database integration

## 📊 Database Schema

```sql
use chat_sys;

-- Users table
CREATE TABLE `users` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
);

-- Friendships table
CREATE TABLE `friendships` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `friend_id` int(11) NOT NULL,
  `status` enum('pending','accepted','rejected') NOT NULL DEFAULT 'pending',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  KEY `friend_id` (`friend_id`),
  UNIQUE KEY `unique_friendship` (`user_id`, `friend_id`),
  CONSTRAINT `friendships_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `friendships_ibfk_2` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);

-- Messages table (private messaging)
CREATE TABLE `messages` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sender_id` int(11) NOT NULL,
  `recipient_id` int(11) NOT NULL,
  `message` text NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT 0,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `sender_id` (`sender_id`),
  KEY `recipient_id` (`recipient_id`),
  KEY `conversation` (`sender_id`, `recipient_id`),
  CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`recipient_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
);
```

## 🚀 Setup Instructions

### 1. Install Dependencies

```bash
go mod download
```

### 2. Configure Database

Make sure MySQL is running and update the database credentials in `main.go`:

```go
username := "root"
password := "12342204"
hostname := "127.0.0.1:3306"
dbname := "chat_sys"
```

### 3. Create Database

Run the SQL setup script:

```bash
mysql -u root -p < setup_database.sql
```

Or manually run the SQL commands from `setup_database.sql`

### 4. Run the Application

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## 📡 API Endpoints

### Authentication
- `POST /api/register` - Register a new user
- `POST /api/login` - Login user

### Friends
- `GET /api/friends?user_id={id}` - Get user's friends list
- `GET /api/friends/search?q={query}&user_id={id}` - Search users
- `POST /api/friends/request` - Send friend request
- `GET /api/friends/requests?user_id={id}` - Get pending friend requests
- `POST /api/friends/accept/{id}` - Accept friend request
- `POST /api/friends/reject/{id}` - Reject friend request

### Messages
- `GET /api/messages/{friendId}?user_id={id}` - Get conversation with a friend
- `POST /api/messages` - Send a message
- `GET /api/messages/unread?user_id={id}` - Get unread message count

### WebSocket
- `GET /ws/{userId}` - WebSocket connection for real-time updates

## 📁 Project Structure

```
.
├── main.go                # Backend server
├── go.mod                 # Go dependencies
├── setup_database.sql     # Database setup script
├── static/
│   ├── index.html        # Frontend HTML
│   ├── style.css         # Styles
│   └── app.js            # Frontend JavaScript
├── README.md             # This file
└── QUICKSTART.md         # Quick start guide
```

## 🛠️ Technologies Used

### Backend
- Go (Golang)
- Gorilla Mux (HTTP router)
- Gorilla WebSocket
- MySQL Driver
- bcrypt (password hashing)

### Frontend
- HTML5
- CSS3
- Vanilla JavaScript
- WebSocket API

## 🎮 Usage

1. Open `http://localhost:8080` in your browser
2. **Register** a new account
3. **Login** with your credentials
4. **Add friends** by clicking the ➕ button and searching for usernames
5. **Accept friend requests** from the notification badge
6. **Select a friend** from your friends list to start chatting
7. **Send messages** in real-time!
8. Open multiple browser windows with different accounts to test the chat system

## 🔒 Security Features

- Password hashing with bcrypt
- SQL injection prevention with prepared statements
- XSS prevention with HTML escaping
- CORS enabled for API access
- Secure friend system (users can only message friends)

## 🌟 Key Features Explained

### Friend System
- Users can search for other users by username
- Send friend requests to other users
- Accept or reject incoming friend requests
- Only friends can message each other

### Private Messaging
- One-on-one conversations
- Real-time message delivery via WebSocket
- Unread message indicators
- Message history persistence

### Real-time Updates
- Instant message delivery
- Friend request notifications
- Unread count updates
- Connection status tracking

## 📝 License

MIT License

## 🙏 Credits

Built with ❤️ using Go and vanilla JavaScript
