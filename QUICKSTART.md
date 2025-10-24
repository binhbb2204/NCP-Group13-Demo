# 🚀 Quick Start Guide - Private Chat System

## Step 1: Setup MySQL Database

Run the SQL setup script:

```bash
mysql -u root -p < setup_database.sql
```

Or manually in MySQL:
```sql
CREATE DATABASE IF NOT EXISTS chat_sys;
USE chat_sys;

-- Copy and run the table creation queries from setup_database.sql
```

## Step 2: Verify Database Connection

Make sure your MySQL credentials in `main.go` are correct:
- Username: root
- Password: 12342204
- Host: 127.0.0.1:3306
- Database: chat_sys

## Step 3: Run the Application

```bash
go run main.go
```

You should see:
```
✅ Successfully connected to MySQL!
🚀 Server starting on http://localhost:8080
```

## Step 4: Open in Browser

Navigate to: http://localhost:8080

## Test the Private Chat System

1. **Register** two or more accounts (use different browser windows/tabs or incognito mode)
2. **Login** with first account
3. **Add friend** by clicking the ➕ button and searching for the other username
4. **Switch to second account** and accept the friend request (check the notification badge)
5. **Select your friend** from the friends list
6. **Start chatting** - messages appear in real-time!
7. See **unread message indicators** when you receive new messages

## Troubleshooting

### Database Connection Error
- Check if MySQL is running
- Verify credentials in main.go
- Ensure chat_sys database exists

### Port Already in Use
- Change port in main.go: `http.ListenAndServe(":8080", ...)`
- Use a different port like `:8081`, `:3000`, etc.

### Dependencies Missing
```bash
go mod download
go mod tidy
```

## Architecture Overview

```
┌─────────────────┐
│   Browser(s)    │ ← Multiple users can connect
│  (Frontend)     │
└────────┬────────┘
         │ HTTP/WebSocket
         ↓
┌─────────────────┐
│   Go Server     │ ← main.go (Backend API + WebSocket)
│  (Backend)      │
└────────┬────────┘
         │ SQL
         ↓
┌─────────────────┐
│  MySQL Database │ ← chat_sys (users, friendships, messages)
│  (Storage)      │
└─────────────────┘
```

## Key Features

✅ **Friend System** - Add friends, send/accept requests
✅ **Private Messaging** - One-on-one conversations
✅ **Real-time Updates** - Instant message delivery via WebSocket
✅ **Unread Indicators** - See unread message counts
✅ **Secure Authentication** - Password hashing with bcrypt
✅ **Message History** - All messages persisted in database
✅ **Modern UI** - Beautiful, responsive design
✅ **Session Management** - Stay logged in with localStorage

Enjoy your chat system! 💬
