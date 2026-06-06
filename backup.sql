-- MySQL dump 10.13  Distrib 9.4.0, for Linux (x86_64)
--
-- Host: localhost    Database: KyzStickersBot_db
-- ------------------------------------------------------
-- Server version       9.4.0

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
-- Table structure for table `line`
--

DROP TABLE IF EXISTS `line`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `line` (
  `line_id` varchar(128) DEFAULT NULL,
  `tg_id` varchar(128) DEFAULT NULL,
  `tg_title` varchar(255) DEFAULT NULL,
  `line_link` varchar(512) DEFAULT NULL,
  `auto_emoji` tinyint(1) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `line`
--

LOCK TABLES `line` WRITE;
/*!40000 ALTER TABLE `line` DISABLE KEYS */;
INSERT INTO `line` VALUES ('32914211','line3291421145ea_by_KyzStickersBot','動く♡ゆるねこ女子（エンジェル） @KyzStickersBot','https://store.line.me/stickershop/product/32914211/en',1);
/*!40000 ALTER TABLE `line` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `properties`
--

DROP TABLE IF EXISTS `properties`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `properties` (
  `name` varchar(128) NOT NULL,
  `value` varchar(128) DEFAULT NULL,
  PRIMARY KEY (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `properties`
--

LOCK TABLES `properties` WRITE;
/*!40000 ALTER TABLE `properties` DISABLE KEYS */;
INSERT INTO `properties` VALUES ('DB_VER','2'),('last_line_dedup_index','-1');
/*!40000 ALTER TABLE `properties` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `stickers`
--

DROP TABLE IF EXISTS `stickers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `stickers` (
  `user_id` bigint DEFAULT NULL,
  `tg_id` varchar(128) DEFAULT NULL,
  `tg_title` varchar(255) DEFAULT NULL,
  `timestamp` bigint DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `stickers`
--

LOCK TABLES `stickers` WRITE;
/*!40000 ALTER TABLE `stickers` DISABLE KEYS */;
INSERT INTO `stickers` VALUES (7939703484,'sticker_d348afe1_by_KyzStickersBot','GIFS',1778677244),(7939703484,'line3291421145ea_by_KyzStickersBot','動く♡ゆるねこ女子（エンジェル） @KyzStickersBot',1778802309),(7939703484,'spacetd_by_KyzStickersBot','𝘦𝘴𝘱𝘢𝘤𝘪𝘰 𝘴𝘪𝘥𝘦𝘳𝘢𝘭 🪐͙͘͡★ @monochromeabril',1779274938),(1291618041,'SuzuZyk_HubStrk_by_KyzStickersBot','SuzuHub Stickers  [@SuzuZyk_Hub]',1779395613),(1291618041,'Zzul_Miku_by_KyzStickersBot','𝗭𝘇𝘂𝗹 / 𝗠𝗶𝗸𝘂 @suzuzyk_hub',1779825582);
/*!40000 ALTER TABLE `stickers` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-06 13:12:58