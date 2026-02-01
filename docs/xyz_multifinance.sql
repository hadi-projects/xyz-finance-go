-- MySQL dump 10.13  Distrib 9.5.0, for macos26.1 (arm64)
--
-- Host: localhost    Database: xyz_multifinance
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
SET @MYSQLDUMP_TEMP_LOG_BIN = @@SESSION.SQL_LOG_BIN;
SET @@SESSION.SQL_LOG_BIN= 0;

--
-- GTID state at the beginning of the backup 
--

SET @@GLOBAL.GTID_PURGED=/*!80000 '+'*/ 'b8094f08-cade-11f0-9c83-87697a275430:1-28270';

--
-- Table structure for table `consumers`
--

DROP TABLE IF EXISTS `consumers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `consumers` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `nik` varchar(16) NOT NULL,
  `full_name` varchar(100) NOT NULL,
  `legal_name` varchar(100) NOT NULL,
  `place_of_birth` varchar(50) DEFAULT NULL,
  `date_of_birth` date DEFAULT NULL,
  `salary` decimal(15,2) DEFAULT NULL,
  `ktp_image` varchar(255) DEFAULT NULL,
  `selfie_image` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_consumers_nik` (`nik`),
  KEY `fk_users_consumer` (`user_id`),
  CONSTRAINT `fk_users_consumer` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `consumers`
--

LOCK TABLES `consumers` WRITE;
/*!40000 ALTER TABLE `consumers` DISABLE KEYS */;
INSERT INTO `consumers` VALUES (1,2,'1234567890123456','Budi Santoso','Budi Santoso','Jakarta','1990-01-01',10000000.00,'budi.webp','budi.jpeg','2026-02-01 22:08:20.280','2026-02-01 22:08:20.280'),(2,3,'6543210987654321','Annisa Putri','Annisa Putri','Bandung','1992-05-15',15000000.00,'annisa.jpeg','annisa.jpeg','2026-02-01 22:08:20.281','2026-02-01 22:08:20.281');
/*!40000 ALTER TABLE `consumers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `limit_mutations`
--

DROP TABLE IF EXISTS `limit_mutations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `limit_mutations` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL,
  `tenor_limit_id` bigint unsigned DEFAULT NULL,
  `old_amount` double DEFAULT NULL,
  `new_amount` double DEFAULT NULL,
  `reason` longtext,
  `action` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_limit_mutations_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `limit_mutations`
--

LOCK TABLES `limit_mutations` WRITE;
/*!40000 ALTER TABLE `limit_mutations` DISABLE KEYS */;
/*!40000 ALTER TABLE `limit_mutations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `permissions`
--

DROP TABLE IF EXISTS `permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `permissions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_permissions_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `permissions`
--

LOCK TABLES `permissions` WRITE;
/*!40000 ALTER TABLE `permissions` DISABLE KEYS */;
INSERT INTO `permissions` VALUES (1,'create-limit','2026-02-01 22:08:20.115','2026-02-01 22:08:20.115'),(2,'delete-limit','2026-02-01 22:08:20.119','2026-02-01 22:08:20.119'),(3,'edit-limit','2026-02-01 22:08:20.120','2026-02-01 22:08:20.120'),(4,'get-audit-log','2026-02-01 22:08:20.120','2026-02-01 22:08:20.120'),(5,'get-auth-log','2026-02-01 22:08:20.121','2026-02-01 22:08:20.121'),(6,'get-limit','2026-02-01 22:08:20.123','2026-02-01 22:08:20.123'),(7,'create-transaction','2026-02-01 22:08:20.125','2026-02-01 22:08:20.125'),(8,'get-transactions','2026-02-01 22:08:20.126','2026-02-01 22:08:20.126');
/*!40000 ALTER TABLE `permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `refresh_tokens`
--

DROP TABLE IF EXISTS `refresh_tokens`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `refresh_tokens` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `token` varchar(512) NOT NULL,
  `expires_at` datetime(3) NOT NULL,
  `revoked` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_refresh_tokens_token` (`token`),
  KEY `idx_refresh_tokens_user_id` (`user_id`),
  CONSTRAINT `fk_refresh_tokens_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `refresh_tokens`
--

LOCK TABLES `refresh_tokens` WRITE;
/*!40000 ALTER TABLE `refresh_tokens` DISABLE KEYS */;
/*!40000 ALTER TABLE `refresh_tokens` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `role_has_permissions`
--

DROP TABLE IF EXISTS `role_has_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `role_has_permissions` (
  `role_id` bigint unsigned NOT NULL,
  `permission_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`role_id`,`permission_id`),
  KEY `fk_role_has_permissions_permission` (`permission_id`),
  CONSTRAINT `fk_role_has_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `permissions` (`id`),
  CONSTRAINT `fk_role_has_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `role_has_permissions`
--

LOCK TABLES `role_has_permissions` WRITE;
/*!40000 ALTER TABLE `role_has_permissions` DISABLE KEYS */;
INSERT INTO `role_has_permissions` VALUES (1,1),(1,2),(1,3),(1,4),(1,5),(2,6),(2,7),(2,8);
/*!40000 ALTER TABLE `role_has_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `roles`
--

LOCK TABLES `roles` WRITE;
/*!40000 ALTER TABLE `roles` DISABLE KEYS */;
INSERT INTO `roles` VALUES (1,'admin','2026-02-01 22:08:20.122','2026-02-01 22:08:20.122'),(2,'user','2026-02-01 22:08:20.127','2026-02-01 22:08:20.127');
/*!40000 ALTER TABLE `roles` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tenor_limits`
--

DROP TABLE IF EXISTS `tenor_limits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tenor_limits` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `tenor_month` bigint NOT NULL COMMENT '''1, 2, 3, or 6''',
  `limit_amount` decimal(15,2) DEFAULT '0.00',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tenor_limits`
--

LOCK TABLES `tenor_limits` WRITE;
/*!40000 ALTER TABLE `tenor_limits` DISABLE KEYS */;
INSERT INTO `tenor_limits` VALUES (1,1,100000.00,'2026-02-01 22:08:20.223','2026-02-01 22:08:20.223'),(2,2,200000.00,'2026-02-01 22:08:20.238','2026-02-01 22:08:20.238'),(3,3,500000.00,'2026-02-01 22:08:20.249','2026-02-01 22:08:20.249'),(4,6,700000.00,'2026-02-01 22:08:20.264','2026-02-01 22:08:20.264'),(5,1,1000000.00,'2026-02-01 22:08:20.267','2026-02-01 22:08:20.267'),(6,2,1200000.00,'2026-02-01 22:08:20.269','2026-02-01 22:08:20.269'),(7,3,1500000.00,'2026-02-01 22:08:20.271','2026-02-01 22:08:20.271'),(8,6,2000000.00,'2026-02-01 22:08:20.275','2026-02-01 22:08:20.275');
/*!40000 ALTER TABLE `tenor_limits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `transactions`
--

DROP TABLE IF EXISTS `transactions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `transactions` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned NOT NULL,
  `contract_number` varchar(50) NOT NULL,
  `otr` decimal(15,2) NOT NULL,
  `admin_fee` decimal(15,2) NOT NULL,
  `installment_amount` decimal(15,2) NOT NULL,
  `interest_amount` decimal(15,2) NOT NULL,
  `asset_name` varchar(255) NOT NULL,
  `status` varchar(20) DEFAULT 'pending',
  `tenor` bigint NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_transactions_contract_number` (`contract_number`),
  KEY `idx_transactions_user_id` (`user_id`),
  KEY `idx_transactions_user_created` (`user_id`,`created_at`),
  CONSTRAINT `fk_transactions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `transactions`
--

LOCK TABLES `transactions` WRITE;
/*!40000 ALTER TABLE `transactions` DISABLE KEYS */;
/*!40000 ALTER TABLE `transactions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `user_has_tenor_limit`
--

DROP TABLE IF EXISTS `user_has_tenor_limit`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_has_tenor_limit` (
  `user_id` bigint unsigned NOT NULL,
  `tenor_limit_id` bigint unsigned NOT NULL,
  PRIMARY KEY (`user_id`,`tenor_limit_id`),
  KEY `fk_user_has_tenor_limit_tenor_limit` (`tenor_limit_id`),
  CONSTRAINT `fk_user_has_tenor_limit_tenor_limit` FOREIGN KEY (`tenor_limit_id`) REFERENCES `tenor_limits` (`id`),
  CONSTRAINT `fk_user_has_tenor_limit_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `user_has_tenor_limit`
--

LOCK TABLES `user_has_tenor_limit` WRITE;
/*!40000 ALTER TABLE `user_has_tenor_limit` DISABLE KEYS */;
INSERT INTO `user_has_tenor_limit` VALUES (2,1),(2,2),(2,3),(2,4),(3,5),(3,6),(3,7),(3,8);
/*!40000 ALTER TABLE `user_has_tenor_limit` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role_id` bigint unsigned NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email` (`email`),
  KEY `fk_users_role` (`role_id`),
  CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
INSERT INTO `users` VALUES (1,'admin@mail.com','$2a$08$oevqWjCWtoR6wAmPr0tct.iypiQcE15MvrA2x2YHmwcYDOny/b5wW',1,'2026-02-01 22:08:20.216','2026-02-01 22:08:20.216'),(2,'budi@mail.com','$2a$08$1GW/XN5bJewnrwVtFC0N/ux1kXrtwuPt017pkGXclz0uKgqNewQxa',2,'2026-02-01 22:08:20.217','2026-02-01 22:08:20.265'),(3,'annisa@mail.com','$2a$08$8AUKj.Y7qr6/.1M5R5pbReBlmfJjEsuhE3bUHl8mfcsR.0c6PbkhK',2,'2026-02-01 22:08:20.219','2026-02-01 22:08:20.276');
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
SET @@SESSION.SQL_LOG_BIN = @MYSQLDUMP_TEMP_LOG_BIN;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-02-01 22:11:05
