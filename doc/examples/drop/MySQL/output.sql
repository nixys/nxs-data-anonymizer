-- MySQL dump 10.13  Distrib 8.4.0, for Linux (x86_64)
--
-- Host: localhost    Database: name_db
-- ------------------------------------------------------
-- Server version	8.4.0

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `authors`
--

DROP TABLE IF EXISTS `authors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `authors` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(50) COLLATE utf8mb3_unicode_ci NOT NULL,
  `last_name` varchar(50) COLLATE utf8mb3_unicode_ci NOT NULL,
  `email` varchar(100) COLLATE utf8mb3_unicode_ci NOT NULL,
  `birthdate` date NOT NULL,
  `added` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `authors`
--

LOCK TABLES `authors` WRITE;
/*!40000 ALTER TABLE `authors` DISABLE KEYS */;
INSERT INTO `authors` VALUES (2,'John','Smith','JohnSmith@example.com','1980-10-28 23:57:00','2000-01-01 12:00:00'),(3,'John','Smith','JohnSmith@example.com','1997-03-14 18:41:06','2000-01-01 12:00:00'),(4,'John','Smith','JohnSmith@example.com','2010-04-12 00:38:42','2000-01-01 12:00:00'),(5,'John','Smith','JohnSmith@example.com','1996-03-25 09:07:39','2000-01-01 12:00:00'),(6,'John','Smith','JohnSmith@example.com','1979-03-16 10:10:12','2000-01-01 12:00:00'),(7,'John','Smith','JohnSmith@example.com','2020-08-02 02:01:34','2000-01-01 12:00:00'),(8,'John','Smith','JohnSmith@example.com','1971-01-02 14:19:52','2000-01-01 12:00:00'),(9,'John','Smith','JohnSmith@example.com','2004-06-10 12:33:11','2000-01-01 12:00:00'),(10,'John','Smith','JohnSmith@example.com','2007-09-19 13:48:20','2000-01-01 12:00:00');
/*!40000 ALTER TABLE `authors` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `posts`
--

DROP TABLE IF EXISTS `posts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `posts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `author_id` int NOT NULL,
  `title` varchar(255) COLLATE utf8mb3_unicode_ci NOT NULL,
  `description` varchar(500) COLLATE utf8mb3_unicode_ci NOT NULL,
  `content` text COLLATE utf8mb3_unicode_ci NOT NULL,
  `date` date NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8mb3 COLLATE=utf8mb3_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `posts`
--

LOCK TABLES `posts` WRITE;
/*!40000 ALTER TABLE `posts` DISABLE KEYS */;
INSERT INTO `posts` VALUES (2,2,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(3,3,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(4,4,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(5,5,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(6,6,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(7,7,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(8,8,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(9,9,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(10,10,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(12,2,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(13,3,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(14,4,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(15,5,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(16,6,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(17,7,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(18,8,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(19,9,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(20,10,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(22,2,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(23,3,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(24,4,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(25,5,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(26,6,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(27,7,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(28,8,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(29,9,'anon_title','anon_description','anon_content','2001-01-01 12:00:00'),(30,10,'anon_title','anon_description','anon_content','2001-01-01 12:00:00');
/*!40000 ALTER TABLE `posts` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-07-19 10:47:53
