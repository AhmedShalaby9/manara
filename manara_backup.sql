-- MySQL dump 10.13  Distrib 9.5.0, for macos15 (arm64)
--
-- Host: localhost    Database: manara
-- ------------------------------------------------------
-- Server version	9.5.0

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
-- Table structure for table `academicYears`
--

DROP TABLE IF EXISTS `academicYears`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `academicYears` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `academicYears`
--

LOCK TABLES `academicYears` WRITE;
/*!40000 ALTER TABLE `academicYears` DISABLE KEYS */;
INSERT INTO `academicYears` VALUES (3,'Second',1,'2025-12-16 19:40:28.000','2025-12-16 19:40:28.000'),(4,'First',1,'2025-12-16 19:43:46.000','2025-12-16 19:43:46.000'),(5,'First',1,'2025-12-16 21:04:36.000','2025-12-16 21:04:36.000'),(6,'test',1,'2025-12-16 22:34:57.000','2025-12-16 22:34:57.000');
/*!40000 ALTER TABLE `academicYears` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `chapters`
--

DROP TABLE IF EXISTS `chapters`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `chapters` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `course_id` bigint unsigned NOT NULL,
  `name` varchar(100) NOT NULL,
  `order` bigint NOT NULL DEFAULT '1',
  `description` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `chapters`
--

LOCK TABLES `chapters` WRITE;
/*!40000 ALTER TABLE `chapters` DISABLE KEYS */;
INSERT INTO `chapters` VALUES (1,1,'Chapter 1: Introduction to Chemistry',1,'Basic concepts and introduction','2025-12-17 17:57:35.824','2025-12-17 17:57:35.824'),(2,1,'Chapter 2: Atoms and Molecules',2,'Understanding atomic structure','2025-12-17 17:57:59.944','2025-12-17 17:57:59.944');
/*!40000 ALTER TABLE `chapters` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `courses`
--

DROP TABLE IF EXISTS `courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `courses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `description` text,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `image_url` varchar(500) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `courses`
--

LOCK TABLES `courses` WRITE;
/*!40000 ALTER TABLE `courses` DISABLE KEYS */;
INSERT INTO `courses` VALUES (1,'Physics','Physics course covering mechanics, thermodynamics, and more','2025-12-17 17:54:23.949','2025-12-17 17:54:23.949',NULL),(2,'Arabic Language','Arabic language course for all levels','2025-12-17 17:54:42.141','2025-12-17 17:54:42.141',NULL);
/*!40000 ALTER TABLE `courses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `lessons`
--

DROP TABLE IF EXISTS `lessons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `lessons` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `chapter_id` bigint unsigned NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` text,
  `order` bigint NOT NULL DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `lessons`
--

LOCK TABLES `lessons` WRITE;
/*!40000 ALTER TABLE `lessons` DISABLE KEYS */;
INSERT INTO `lessons` VALUES (1,1,1,'Lesson 1: Basic Concepts','Introduction to basic chemistry concepts',1,'2025-12-17 18:00:34.688','2025-12-17 18:00:34.688'),(2,1,1,'Lesson 2: Atomic Structure','Understanding the structure of atoms',2,'2025-12-17 18:00:46.723','2025-12-17 18:00:46.723'),(3,1,1,'Lesson 3: Chemical Bonding','Types of chemical bonds',3,'2025-12-17 18:00:49.860','2025-12-17 18:00:49.860');
/*!40000 ALTER TABLE `lessons` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_name` varchar(100) NOT NULL,
  `role_value` varchar(100) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_role_value` (`role_value`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'Admin','admin','2025-12-07 19:10:05.000','2025-12-07 19:10:05.000'),(2,'Super Admin','super_admin','2025-12-07 19:10:05.000','2025-12-07 19:10:05.000'),(3,'Teacher','teacher','2025-12-07 19:10:05.000','2025-12-07 19:10:05.000'),(4,'Student','student','2025-12-07 19:10:05.000','2025-12-07 19:10:05.000');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `schema_migrations`
--

LOCK TABLES `schema_migrations` WRITE;
/*!40000 ALTER TABLE `schema_migrations` DISABLE KEYS */;
INSERT INTO `schema_migrations` VALUES (1,0);
/*!40000 ALTER TABLE `schema_migrations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `students`
--

DROP TABLE IF EXISTS `students`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `students` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `teacher_id` bigint unsigned NOT NULL,
  `academic_year_id` bigint unsigned NOT NULL DEFAULT '3',
  `parent_phone` varchar(20) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_students_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `students`
--

LOCK TABLES `students` WRITE;
/*!40000 ALTER TABLE `students` DISABLE KEYS */;
/*!40000 ALTER TABLE `students` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teacher_courses`
--

DROP TABLE IF EXISTS `teacher_courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teacher_courses` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `teacher_id` bigint unsigned NOT NULL,
  `course_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teacher_courses`
--

LOCK TABLES `teacher_courses` WRITE;
/*!40000 ALTER TABLE `teacher_courses` DISABLE KEYS */;
/*!40000 ALTER TABLE `teacher_courses` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `teachers`
--

DROP TABLE IF EXISTS `teachers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `teachers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `bio` text,
  `specialization` varchar(255) DEFAULT NULL,
  `experience_years` bigint DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_teachers_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `teachers`
--

LOCK TABLES `teachers` WRITE;
/*!40000 ALTER TABLE `teachers` DISABLE KEYS */;
INSERT INTO `teachers` VALUES (1,15,'Experienced mathematics teacher with a passion for making complex concepts simple','Mathematics & Physics',5,'2025-12-17 17:59:23.711','2025-12-17 17:59:23.711');
/*!40000 ALTER TABLE `teachers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `email` varchar(255) NOT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `user_name` varchar(100) NOT NULL,
  `password_hash` varchar(255) DEFAULT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`),
  UNIQUE KEY `uni_users_user_name` (`user_name`),
  KEY `fk_roles_users` (`role_id`),
  CONSTRAINT `fk_roles_users` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `users_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (2,'Omar','Student',4,'2025-12-07 19:12:03.000','2025-12-07 20:21:11.000','omar.student@example.com','0100000002','omarstudent','hashedpass123',1),(4,'Admin','Root',2,'2025-12-07 19:12:03.000','2025-12-07 20:21:12.000','root@example.com','0100000004','superadmin','hashedpass123',1),(5,'Ahmed','Gamal',2,'2025-12-07 23:21:42.000','2025-12-07 23:29:10.000','ahmed@example.com','01012345678','ahmedgamal','$2a$10$2T7maKHEnngzstVqwyE9TO7dcKcLqdj1ui2nNFeJigssEIxtnpfZW',1),(6,'Mohamed','Hassan',3,'2025-12-14 17:51:05.000','2025-12-14 17:51:05.000','mohamed.teacher@manara.com','01011111111','mohamedteacher','$2a$10$2WzF5Aj0poTRaOcAjQiud.jQt8OuB/pNEQTldTmaqnmAFpyHyWHMy',1),(7,'Mohamed','Hassan',4,'2025-12-15 19:11:57.000','2025-12-15 19:11:57.000','mohamed.st@manara.com','01022222222','mohamedstudent','$2a$10$fGsJ77yRMoMn7kQwDDg03eBr3J90L/yDnFLXUXWRaQteWm8FDgk5a',1),(10,'Ahmed','Ganmal',4,'2025-12-16 20:25:11.551','2025-12-16 20:25:11.551','a.gamal@manara.com','212121212','AhmedGamal2020','$2a$10$O2l97Kq8xqLe1vZR1J1Do.usuViMwmTtF6mGI60QUw2gxoeVLTWzq',1),(11,'Ahmed','Ganmal',4,'2025-12-16 20:31:42.526','2025-12-16 20:31:42.526','a.gamaal@manara.com','212121212','AhmedGamal20202','$2a$10$VUQMZpyOk55Vn/uFXSgmkuDiE37peKGlOeqJFmQCVTn1sX3NLxSSO',1),(12,'Ahmed','Ganmal',4,'2025-12-17 00:37:51.508','2025-12-17 00:37:51.508','a.gamaal11@manara.com','212121312','ahmed2121','$2a$10$IMWTsHxuSBi0AC45G9zenevClTsyI3fUjvcCTEFm.jS513N1h6AM.',1),(13,'Ahmed','Ganmal',4,'2025-12-17 00:39:42.493','2025-12-17 00:39:42.493','a.gamaal112@manara.com','212121312','ahmed2121q','$2a$10$DJOC/JegbpkAXvmeiWu0keeFY8O623PzBoQYOzgDZLdl8wZkvJPPe',1),(14,'Ahmed','Ganmal',4,'2025-12-17 00:57:00.203','2025-12-17 00:57:00.203','a.gamaal12212@manara.com','212121012','ahmed2121q1','$2a$10$4XSfmAZLyWGHbzBf0jZqauS4inPS5UHQuVyOwDppMkME3u1WO35Ei',1),(15,'Mohamed','Hassan',3,'2025-12-17 17:59:23.707','2025-12-17 17:59:23.707','mohamed.teacher22@manara.com','01011111111','mohamedteacher22','$2a$10$0Qux3dD8Hk3f0FfcU1F6xuf.8a/QZxHNMCqp0jJlqI7FJSfD/NXM6',1);
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

-- Dump completed on 2025-12-17 22:19:31
