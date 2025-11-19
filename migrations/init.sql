-- Converted MySQL dump to SQLite-compatible schema and data
PRAGMA foreign_keys = OFF;
BEGIN TRANSACTION;

-- Users table
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT NOT NULL UNIQUE,
  password TEXT NOT NULL,
  email TEXT,
  bio TEXT,
  avatar_color TEXT DEFAULT '#8774e1',
  status TEXT DEFAULT 'offline',
  created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP)
);

-- Friendships table
CREATE TABLE IF NOT EXISTS friendships (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  friend_id INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending',
  created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  updated_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  UNIQUE(user_id, friend_id),
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(friend_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Messages table
CREATE TABLE IF NOT EXISTS messages (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  sender_id INTEGER NOT NULL,
  recipient_id INTEGER NOT NULL,
  message TEXT NOT NULL,
  is_read INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP),
  FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Insert users
INSERT INTO users (id, username, password, email, bio, avatar_color, status, created_at, updated_at) VALUES
(2,'binhbb','$2a$10$gEeXIT47K/Cw.B7CPoRuoetcX8n9Gdoey9MLOcO.xSa0afNDiJMWe',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),
(3,'testing','$2a$10$F2ZS71NffP0cgwW8jKVgT.ht5xYBIzfoIV.sAdGpn8rcq5ILw.7y6',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),
(4,'thuannm','$2a$10$IjJrClXWzskTBBuoGWkJGeRRKgr9bXr3sHybk8fGPmhUmAgNBfgvG',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),
(5,'kittyo','$2a$10$ktHO7LOpGNR/rpfq9cscPeAWP2eWvkM8gTJmXCwpl7GD2WGjv5gpO',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),
(6,'risa','$2a$10$ekYhzV3e/oAFJlLDQ/BAd.0eZp2iSWhcSj/Sp7c5ha/8O7s/Sm4EC',NULL,NULL,'#8774e1','offline','2025-10-15 14:48:51','2025-10-15 14:48:51');

-- Insert friendships
INSERT INTO friendships (id, user_id, friend_id, status, created_at, updated_at) VALUES
(1,3,2,'accepted','2025-10-14 17:06:46','2025-10-14 17:06:52'),
(2,4,2,'accepted','2025-10-14 17:54:30','2025-10-14 17:54:44'),
(3,6,2,'accepted','2025-10-15 14:49:17','2025-10-15 14:49:23');

-- Insert messages
INSERT INTO messages (id, sender_id, recipient_id, message, is_read, created_at) VALUES
(1,2,3,'hello',1,'2025-10-14 17:20:48'),
(2,3,2,'hiiii',1,'2025-10-14 17:21:07'),
(3,2,3,'How are you my nigga',1,'2025-10-14 17:21:50'),
(4,2,3,'you dont answer me?:(',1,'2025-10-14 17:23:07'),
(5,3,2,'sowwy i was busy :(',1,'2025-10-14 17:23:13'),
(6,2,3,'its fine',1,'2025-10-14 17:50:15'),
(7,4,2,'wassup',1,'2025-10-14 17:54:12'),
(8,2,4,'sup',1,'2025-10-14 17:54:50'),
(9,2,3,'hiii',1,'2025-10-15 08:02:57'),
(10,3,2,'hiii',1,'2025-10-15 08:03:14'),
(11,3,2,'how are you?',1,'2025-10-15 08:03:19'),
(12,2,3,'its great',1,'2025-10-15 08:03:32'),
(13,2,3,'sup',1,'2025-10-15 08:50:58'),
(14,2,3,'đi vũng tàu ko',1,'2025-10-15 08:54:09'),
(15,3,2,'đi :V',1,'2025-10-15 08:54:30'),
(16,2,3,'oke lun bạn ey',1,'2025-10-15 08:55:22'),
(17,3,2,'mấy giờ đi',1,'2025-10-15 08:55:27'),
(18,3,2,'mai đi',1,'2025-10-15 08:55:38'),
(19,2,3,'oke lun',1,'2025-10-15 08:55:42'),
(20,2,3,'mà mang gì đii á',1,'2025-10-15 08:55:52'),
(21,3,2,'mang áo bơi chứ mang gì',1,'2025-10-15 08:56:00'),
(22,3,2,'à thế à',1,'2025-10-15 09:15:17'),
(23,2,3,'ừa',1,'2025-10-15 09:15:22'),
(24,3,2,'oke mai t chở m',1,'2025-10-15 09:15:59'),
(25,2,3,'oke bạn ei',1,'2025-10-15 09:16:04'),
(26,2,3,'<3',1,'2025-10-15 09:16:06'),
(27,3,2,'hihi',1,'2025-10-15 09:16:09'),
(28,3,2,'mai mưa đó',1,'2025-10-15 09:16:17'),
(29,3,2,':V',1,'2025-10-15 09:16:19'),
(30,2,3,'mịa mày V:',1,'2025-10-15 09:16:25'),
(31,2,3,':v',1,'2025-10-15 09:17:52'),
(32,2,3,'bruh',1,'2025-10-15 09:18:17'),
(33,2,3,'hmmmm',1,'2025-10-15 09:19:21'),
(34,3,2,'hmm?',1,'2025-10-15 09:19:47'),
(35,2,3,'huh?',1,'2025-10-15 09:21:19'),
(36,3,2,'imma take a nap cya',1,'2025-10-15 09:23:54'),
(37,2,3,'oke bye bye',1,'2025-10-15 09:24:00'),
(38,2,6,'hewwoo babeeeeee :DDDDDDD',1,'2025-10-15 14:49:31'),
(39,2,6,'aaaaa babeeeee',1,'2025-10-15 14:52:53'),
(40,6,2,'Hello babyy',1,'2025-10-15 14:53:12'),
(41,2,6,'yayyyy it worksss babeee :DDDD',1,'2025-10-15 14:53:38'),
(42,3,2,'heil hitler',1,'2025-10-16 01:38:28'),
(43,2,4,'dit me may T',1,'2025-10-24 06:30:57'),
(44,4,2,'!!!!',1,'2025-10-24 06:32:23'),
(45,2,4,'heil hitler',1,'2025-10-24 06:32:26'),
(46,2,4,'T chan',1,'2025-10-24 06:33:14'),
(47,2,4,'suck my di',1,'2025-10-24 06:33:15'),
(48,4,2,'LOL',1,'2025-10-24 06:34:28'),
(49,2,4,'i love you nigga',0,'2025-10-24 06:35:28'),
(50,2,4,'bruh',0,'2025-11-19 05:06:41');

COMMIT;
PRAGMA foreign_keys = ON;
