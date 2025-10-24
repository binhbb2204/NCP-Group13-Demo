# 🎨 Account Settings Feature

## ✨ Features Added

### 1. **Profile Management**
- ✅ Update email address
- ✅ Add/edit bio
- ✅ Customize avatar color with color picker
- ✅ View account creation date

### 2. **Security**
- ✅ Change password securely
- ✅ Current password verification
- ✅ Password strength validation

---

## 📋 Setup Instructions

### Step 1: Update Your Database

**Option A - Fresh Install:**
```sql
-- Drop and recreate (WARNING: Deletes all data!)
DROP DATABASE IF EXISTS chat_sys;
```
Then run:
```bash
mysql -u root -p < setup_database.sql
```

**Option B - Update Existing Database:**
```bash
mysql -u root -p < update_database.sql
```

### Step 2: Run the Application
```powershell
go run main.go
```

---

## 🎯 How to Use

### Accessing Account Settings:
1. Login to your account
2. Click the **⚙️ Settings icon** in the top-right corner
3. Two tabs will appear:
   - **Profile**: Edit email, bio, and avatar color
   - **Security**: Change your password

### Profile Tab:
- **Username**: Read-only (cannot be changed)
- **Email**: Optional contact email
- **Bio**: Tell others about yourself
- **Avatar Color**: Click the color picker to choose your favorite color

### Security Tab:
- Enter your **current password**
- Enter your **new password** (min 4 characters)
- **Confirm** the new password
- Click "Change Password"

---

## 🎨 New Database Fields

```sql
users table:
- email          VARCHAR(255)   -- User's email (optional, unique)
- bio            TEXT           -- User biography
- avatar_color   VARCHAR(7)     -- Hex color code for avatar
- status         ENUM           -- online/offline/away
- updated_at     TIMESTAMP      -- Last profile update time
```

---

## 🔒 Security Features

- ✅ Passwords are hashed with bcrypt
- ✅ Old password verification before changing
- ✅ Password strength validation (min 4 chars)
- ✅ Email uniqueness enforced
- ✅ XSS protection with input escaping

---

## 🎨 UI Features

### Settings Modal:
- Modern tabbed interface
- Smooth animations
- Dark/Light mode compatible
- Color picker with live preview
- Success/Error message feedback

### Avatar Colors:
- Gradient effect applied automatically
- Updates across all UI elements
- Saved to database

---

## 🚀 API Endpoints

### Profile Management:
```
GET  /api/user/profile?user_id={id}  - Get user profile
PUT  /api/user/profile?user_id={id}  - Update profile
PUT  /api/user/password?user_id={id} - Change password
```

### Request/Response Examples:

**Update Profile:**
```json
PUT /api/user/profile?user_id=1
{
  "email": "user@example.com",
  "bio": "I love coding!",
  "avatar_color": "#8774e1"
}
```

**Change Password:**
```json
PUT /api/user/password?user_id=1
{
  "old_password": "oldpass123",
  "new_password": "newpass456"
}
```

---

## 🐛 Troubleshooting

### "User ID is required"
- Make sure you're logged in
- Check that currentUser object exists

### "Current password is incorrect"
- Double-check your current password
- Passwords are case-sensitive

### Database errors:
```bash
# Check if columns exist
mysql -u root -p
USE chat_sys;
DESCRIBE users;
```

---

## 💡 Tips

1. **Choose a unique avatar color** to stand out in chats!
2. **Add a bio** to let friends know about you
3. **Change your password regularly** for better security
4. **Email is optional** - only add it if you want

---

Enjoy your new account settings! 🎉
