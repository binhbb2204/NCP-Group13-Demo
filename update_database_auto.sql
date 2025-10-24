-- AUTOMATIC UPDATE: Detects what you need and adds it
-- Works with any existing users table structure!

USE chat_sys;

-- Show current structure
SELECT 'Your current users table structure:' as Info;
DESCRIBE users;

-- Add id column if missing (should exist)
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'id');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN id INT(11) NOT NULL AUTO_INCREMENT PRIMARY KEY FIRST', 'SELECT "id exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add username column if missing (should exist)
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'username');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN username VARCHAR(255) NOT NULL UNIQUE', 'SELECT "username exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add password column if missing (should exist)
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'password');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN password VARCHAR(255) NOT NULL', 'SELECT "password exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add email column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'email');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN email VARCHAR(255) DEFAULT NULL', 'SELECT "email exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add bio column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'bio');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN bio TEXT DEFAULT NULL', 'SELECT "bio exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add avatar_color column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'avatar_color');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN avatar_color VARCHAR(7) DEFAULT "#8774e1"', 'SELECT "avatar_color exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add status column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'status');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN status ENUM("online","offline","away") DEFAULT "offline"', 'SELECT "status exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add created_at column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'created_at');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP', 'SELECT "created_at exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Add updated_at column if missing
SET @col_exists = (SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA = 'chat_sys' AND TABLE_NAME = 'users' AND COLUMN_NAME = 'updated_at');
SET @query = IF(@col_exists = 0, 'ALTER TABLE users ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP', 'SELECT "updated_at exists" as Note');
PREPARE stmt FROM @query; EXECUTE stmt; DEALLOCATE PREPARE stmt;

-- Show updated structure
SELECT '✅ Update complete! New structure:' as Status;
DESCRIBE users;

-- Show all users (your data is safe!)
SELECT 'All existing users (data preserved):' as Info;
SELECT id, username, COALESCE(email, 'not set') as email, avatar_color FROM users;

SELECT CONCAT('✅ Total users in database: ', COUNT(*)) as Summary FROM users;
