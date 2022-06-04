-- MariaDB dump 10.19  Distrib 10.7.3-MariaDB, for Linux (x86_64)
--
-- Database: dousheng
-- ------------------------------------------------------
-- Server version	10.7.3-MariaDB-1:10.7.3+maria~focal

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `users` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `user_name` varchar(32) NOT NULL,
  `salt` tinyblob NOT NULL,
  `password` tinyblob NOT NULL,
  `follow_count` bigint(20) DEFAULT 0,
  `follower_count` bigint(20) DEFAULT 0,
  `created_at` datetime(3) DEFAULT current_timestamp(3),
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_user_user_name` (`user_name`),
  FULLTEXT KEY `username` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=201 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;


-- ----------------------------
-- Table structure for videos
-- ----------------------------
DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos`  (
   `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '短视频id',
   `author_id` bigint NOT NULL COMMENT '作者id',
   `play_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '短视频url',
   `cover_url` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '封面url',
   `favorite_count` bigint NOT NULL DEFAULT 0 COMMENT '点赞数',
   `comment_count` bigint NOT NULL DEFAULT 0 COMMENT '评论数',
   `created_at` bigint NOT NULL COMMENT '投递时间',
   `deleted_at` datetime(3) default NULL COMMENT '删除标记位',
   PRIMARY KEY (`id`) USING BTREE,
   #created_at聚簇索引
   UNIQUE INDEX `create_time_index`(`created_at`, `deleted_at`, `id`, `author_id`, `play_url`, `cover_url`, `favorite_count`, `comment_count`) USING BTREE,
   #deleted_at普通索引
   INDEX `idx_videos_deleted_at`(`deleted_at`) USING BTREE,
   #author_id普通索引
   INDEX `fk_videos_author`(`author_id`) USING BTREE,
   #users.id->videos.author_id外键约束
   CONSTRAINT `fk_videos_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`) ON DELETE RESTRICT ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;