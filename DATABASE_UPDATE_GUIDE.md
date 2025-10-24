# 📊 Database Update Guide - Keep Your Data Safe!

## ✅ Safe Update Process

Your existing users, messages, and friendships **WILL NOT BE DELETED**!
We're only **ADDING** new columns to the `users` table.

---

## 🚀 Option 1: Simple Update (Recommended)

If this is your **first time** updating:

```bash
mysql -u root -p < update_database.sql
```

**What happens:**
- ✅ Adds 5 new columns: `email`, `bio`, `avatar_color`, `status`, `updated_at`
- ✅ All existing users remain intact
- ✅ New columns get default values (NULL for optional fields, #8774e1 for avatar_color)

---

## 🛡️ Option 2: Extra Safe Update

If you want to be **extra careful** or might run it multiple times:

```bash
mysql -u root -p < update_database_safe.sql
```

**What happens:**
- ✅ Checks if columns exist before adding them
- ✅ Won't error if you run it twice
- ✅ Shows you the results
- ✅ 100% safe for existing data

---

## 📋 Step-by-Step Instructions

### 1. **Backup First (Optional but Recommended)**
```bash
mysqldump -u root -p chat_sys > backup_before_update.sql
```

### 2. **Run the Update**
```bash
# Enter your MySQL password when prompted
mysql -u root -p < update_database.sql
```

### 3. **Check the Results**
```bash
mysql -u root -p
```

Then in MySQL:
```sql
USE chat_sys;
DESCRIBE users;
SELECT * FROM users;
```

You should see:
```
+-------------+----------+------+-----+-------------------+
| Field       | Type     | Null | Key | Default           |
+-------------+----------+------+-----+-------------------+
| id          | int(11)  | NO   | PRI | NULL              |
| username    | varchar  | NO   | UNI | NULL              |
| password    | varchar  | NO   |     | NULL              |
| email       | varchar  | YES  | UNI | NULL              | ← NEW!
| bio         | text     | YES  |     | NULL              | ← NEW!
| avatar_color| varchar  | YES  |     | #8774e1           | ← NEW!
| status      | enum     | YES  |     | offline           | ← NEW!
| created_at  | timestamp| NO   |     | CURRENT_TIMESTAMP |
| updated_at  | timestamp| NO   |     | CURRENT_TIMESTAMP | ← NEW!
+-------------+----------+------+-----+-------------------+
```

---

## ❓ What If Something Goes Wrong?

### Error: "Duplicate column name"
**Meaning:** Column already exists (you already ran the update)
**Solution:** Everything is fine! The column is already there.

### Error: "Duplicate entry for key 'email'"
**Meaning:** Two users somehow have the same email
**Solution:** 
```sql
-- Find duplicates
SELECT email, COUNT(*) FROM users WHERE email IS NOT NULL GROUP BY email HAVING COUNT(*) > 1;

-- Set duplicate emails to NULL
UPDATE users SET email = NULL WHERE email = 'duplicate@email.com' AND id != 1;

-- Then run the update again
```

### Want to restore backup?
```bash
mysql -u root -p chat_sys < backup_before_update.sql
```

---

## 🎯 What Happens to Your Data?

### Before Update:
```
users table:
- id: 1, username: "alice", password: "$2a$..."
- id: 2, username: "bob", password: "$2a$..."
```

### After Update:
```
users table:
- id: 1, username: "alice", password: "$2a$...", email: NULL, bio: NULL, avatar_color: "#8774e1"
- id: 2, username: "bob", password: "$2a$...", email: NULL, bio: NULL, avatar_color: "#8774e1"
```

**All usernames, passwords, and IDs stay exactly the same!**

---

## ✨ After Update

Your users can now:
1. Login with their existing accounts ✅
2. Update their profile in settings ✅
3. Add email, bio, and choose avatar color ✅
4. Change their password ✅

---

## 🔍 Quick Verification

After running the update, verify everything:

```bash
mysql -u root -p -e "USE chat_sys; SELECT COUNT(*) as 'Total Users' FROM users;"
```

This should show the **same number** of users as before!

---

**Your data is safe! The ALTER TABLE command only adds columns, it doesn't delete rows.** 🎉
