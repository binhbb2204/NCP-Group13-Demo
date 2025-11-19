-- MySQL dump 10.13  Distrib 8.0.42, for Win64 (x86_64)
--
-- Host: localhost    Database: chat_sys
-- ------------------------------------------------------
-- Server version	8.0.42

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `friendships`
--

DROP TABLE IF EXISTS `friendships`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `friendships` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `friend_id` int NOT NULL,
  `status` enum('pending','accepted','rejected') NOT NULL DEFAULT 'pending',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_friendship` (`user_id`,`friend_id`),
  KEY `user_id` (`user_id`),
  KEY `friend_id` (`friend_id`),
  CONSTRAINT `friendships_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `friendships_ibfk_2` FOREIGN KEY (`friend_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `friendships`
--

LOCK TABLES `friendships` WRITE;
/*!40000 ALTER TABLE `friendships` DISABLE KEYS */;
INSERT INTO `friendships` VALUES (1,3,2,'accepted','2025-10-14 17:06:46','2025-10-14 17:06:52'),(2,4,2,'accepted','2025-10-14 17:54:30','2025-10-14 17:54:44'),(3,6,2,'accepted','2025-10-15 14:49:17','2025-10-15 14:49:23');
/*!40000 ALTER TABLE `friendships` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `messages`
--

DROP TABLE IF EXISTS `messages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `messages` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sender_id` int NOT NULL,
  `recipient_id` int NOT NULL,
  `message` text NOT NULL,
  `is_read` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `sender_id` (`sender_id`),
  KEY `recipient_id` (`recipient_id`),
  KEY `conversation` (`sender_id`,`recipient_id`),
  CONSTRAINT `messages_ibfk_1` FOREIGN KEY (`sender_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
  CONSTRAINT `messages_ibfk_2` FOREIGN KEY (`recipient_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `messages`
--

LOCK TABLES `messages` WRITE;
/*!40000 ALTER TABLE `messages` DISABLE KEYS */;
INSERT INTO `messages` VALUES (1,2,3,'hello',1,'2025-10-14 17:20:48'),(2,3,2,'hiiii',1,'2025-10-14 17:21:07'),(3,2,3,'How are you my nigga',1,'2025-10-14 17:21:50'),(4,2,3,'you dont answer me?:(',1,'2025-10-14 17:23:07'),(5,3,2,'sowwy i was busy :(',1,'2025-10-14 17:23:13'),(6,2,3,'its fine',1,'2025-10-14 17:50:15'),(7,4,2,'wassup',1,'2025-10-14 17:54:12'),(8,2,4,'sup',1,'2025-10-14 17:54:50'),(9,2,3,'hiii',1,'2025-10-15 08:02:57'),(10,3,2,'hiii',1,'2025-10-15 08:03:14'),(11,3,2,'how are you?',1,'2025-10-15 08:03:19'),(12,2,3,'its great',1,'2025-10-15 08:03:32'),(13,2,3,'sup',1,'2025-10-15 08:50:58'),(14,2,3,'đi vũng tàu ko',1,'2025-10-15 08:54:09'),(15,3,2,'đi :V',1,'2025-10-15 08:54:30'),(16,2,3,'oke lun bạn ey',1,'2025-10-15 08:55:22'),(17,3,2,'mấy giờ đi',1,'2025-10-15 08:55:27'),(18,3,2,'mai đi',1,'2025-10-15 08:55:38'),(19,2,3,'oke lun',1,'2025-10-15 08:55:42'),(20,2,3,'mà mang gì đii á',1,'2025-10-15 08:55:52'),(21,3,2,'mang áo bơi chứ mang gì',1,'2025-10-15 08:56:00'),(22,3,2,'à thế à',1,'2025-10-15 09:15:17'),(23,2,3,'ừa',1,'2025-10-15 09:15:22'),(24,3,2,'oke mai t chở m',1,'2025-10-15 09:15:59'),(25,2,3,'oke bạn ei',1,'2025-10-15 09:16:04'),(26,2,3,'<3',1,'2025-10-15 09:16:06'),(27,3,2,'hihi',1,'2025-10-15 09:16:09'),(28,3,2,'mai mưa đó',1,'2025-10-15 09:16:17'),(29,3,2,':V',1,'2025-10-15 09:16:19'),(30,2,3,'mịa mày V:',1,'2025-10-15 09:16:25'),(31,2,3,':v',1,'2025-10-15 09:17:52'),(32,2,3,'bruh',1,'2025-10-15 09:18:17'),(33,2,3,'hmmmm',1,'2025-10-15 09:19:21'),(34,3,2,'hmm?',1,'2025-10-15 09:19:47'),(35,2,3,'huh?',1,'2025-10-15 09:21:19'),(36,3,2,'imma take a nap cya',1,'2025-10-15 09:23:54'),(37,2,3,'oke bye bye',1,'2025-10-15 09:24:00'),(38,2,6,'hewwoo babeeeeee :DDDDDDD',1,'2025-10-15 14:49:31'),(39,2,6,'aaaaa babeeeee',1,'2025-10-15 14:52:53'),(40,6,2,'Hello babyy',1,'2025-10-15 14:53:12'),(41,2,6,'yayyyy it worksss babeee :DDDD',1,'2025-10-15 14:53:38'),(42,3,2,'heil hitler',1,'2025-10-16 01:38:28'),(43,2,4,'dit me may T',1,'2025-10-24 06:30:57'),(44,4,2,'!!!!',1,'2025-10-24 06:32:23'),(45,2,4,'heil hitler',1,'2025-10-24 06:32:26'),(46,2,4,'T chan',1,'2025-10-24 06:33:14'),(47,2,4,'suck my di',1,'2025-10-24 06:33:15'),(48,4,2,'LOL',1,'2025-10-24 06:34:28'),(49,2,4,'i love you nigga',0,'2025-10-24 06:35:28'),(50,2,4,'bruh',0,'2025-11-19 05:06:41');
/*!40000 ALTER TABLE `messages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) DEFAULT NULL,
  `bio` text,
  `avatar_color` varchar(7) DEFAULT '#8774e1',
  `status` enum('online','offline','away') DEFAULT 'offline',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'binhbb','$2a$10$gEeXIT47K/Cw.B7CPoRuoetcX8n9Gdoey9MLOcO.xSa0afNDiJMWe',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),(3,'testing','$2a$10$F2ZS71NffP0cgwW8jKVgT.ht5xYBIzfoIV.sAdGpn8rcq5ILw.7y6',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),(4,'thuannm','$2a$10$IjJrClXWzskTBBuoGWkJGeRRKgr9bXr3sHybk8fGPmhUmAgNBfgvG',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),(5,'kittyo','$2a$10$ktHO7LOpGNR/rpfq9cscPeAWP2eWvkM8gTJmXCwpl7GD2WGjv5gpO',NULL,NULL,'#8774e1','offline','2025-10-15 08:17:11','2025-10-15 08:17:11'),(6,'risa','$2a$10$ekYhzV3e/oAFJlLDQ/BAd.0eZp2iSWhcSj/Sp7c5ha/8O7s/Sm4EC',NULL,NULL,'#8774e1','offline','2025-10-15 14:48:51','2025-10-15 14:48:51');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-11-19 12:08:38
