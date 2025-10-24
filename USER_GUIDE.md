# 📖 How to Use the Private Chat System

## 🚀 Getting Started

### Step 1: Create Accounts
1. Open http://localhost:8080
2. Click **Register** and create your first account (e.g., username: "alice")
3. Open a **new browser window** (or incognito/private window)
4. Register a second account (e.g., username: "bob")

### Step 2: Send Friend Request
1. Login as **Alice**
2. Click the **➕ button** in the top-right of the sidebar
3. Search for "bob" in the search box
4. Click **Add Friend**
5. You'll see: "✅ Friend request sent to bob! They need to accept your request before you can chat."

### Step 3: Accept Friend Request
1. Switch to **Bob's window**
2. You should see a **yellow notification badge** at the top: "📬 You have 1 friend request(s) - Click to view!"
3. **Click on the yellow notification badge**
4. A modal will open showing the friend request from Alice
5. Click **Accept** button

### Step 4: Start Chatting!
1. Now **both Alice and Bob** will see each other in their friends list
2. Click on your friend's name in the sidebar
3. Type a message and press **Enter** or click **Send**
4. Messages appear in real-time! 🎉

## 📋 Important Notes

### Friend System Flow
```
User A                          User B
  |                               |
  |------ Send Friend Request --->|
  |                               |
  |                         [Notification Badge]
  |                               |
  |                         [Click Badge]
  |                               |
  |                         [Accept Request]
  |                               |
  |<----- Both Now Friends ------>|
  |                               |
  |<----- Can Chat Now ---------->|
```

### Why Don't I See My Friend?
- ❌ **Just sent request** → Friend request is pending, wait for acceptance
- ✅ **Request accepted** → Friend will appear in friends list
- ✅ **Both users can now see each other** and start chatting

### Tips
- 💡 The **yellow notification badge** pulses to get your attention
- 💡 **Unread messages** show a red badge with the count
- 💡 You can have **multiple friends** and switch between chats
- 💡 **Messages persist** - they're saved in the database
- 💡 Use **multiple browser windows** to test with different accounts

## 🐛 Troubleshooting

### "No friends yet" message
**This is normal!** It means:
1. No one has sent you a friend request yet, OR
2. You sent requests but they haven't been accepted yet

**Solution:** 
- Check if you have pending friend requests (yellow badge)
- Wait for others to accept your requests
- Have your friend accept your request in their account

### Friend request not showing
**Make sure:**
1. The other user is logged in and refreshed their page
2. Check the browser console for errors (F12)
3. Verify WebSocket is connected (check console logs)

### Can't send messages
**This means:**
- You haven't selected a friend to chat with
- Click on a friend's name in the sidebar first
- Make sure the friend request was accepted

## 🎯 Quick Test (2 Users)

```bash
# Window 1 - Alice
1. Register as "alice"
2. Login as "alice"  
3. Click ➕ → Search "bob" → Add Friend
4. Wait for Bob to accept...

# Window 2 - Bob
1. Register as "bob"
2. Login as "bob"
3. See yellow notification badge
4. Click badge → Accept Alice's request
5. Click "alice" in friends list
6. Send message: "Hi Alice!"

# Window 1 - Alice
1. See "bob" appear in friends list
2. Click "bob"
3. See Bob's message in real-time!
4. Reply: "Hi Bob!"
```

## ✨ Features

- 👥 **Friend System** - Send, accept, reject requests
- 💬 **Private Chats** - One-on-one conversations
- ⚡ **Real-time** - Instant message delivery
- 🔔 **Notifications** - Unread counts and request alerts
- 💾 **Persistent** - All messages saved to database
- 🔐 **Secure** - Password hashing and private conversations

Enjoy chatting! 🎊
