/*
 Navicat Premium Data Transfer

 Source Server Type    : MySQL
 Source Server Version : 100703
 Source Host           : fitenne.com:3306
 Source Schema         : dousheng

 Target Server Type    : MySQL
 Target Server Version : 100703
 File Encoding         : 65001

 Date: 21/05/2022 18:28:14
*/

DROP TABLE IF EXISTS `videos`;
CREATE TABLE `videos` (
                          `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '短视频id',
                          `author_id` bigint(20) NOT NULL COMMENT '作者id',
                          `play_url` varchar(50) NOT NULL COMMENT '短视频url',
                          `cover_url` varchar(50) NOT NULL COMMENT '封面url',
                          `favorite_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '点赞数',
                          `comment_count` bigint(20) NOT NULL DEFAULT 0 COMMENT '评论数',
                          `created_at` bigint(20) NOT NULL COMMENT '投递时间',
                          `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除标记位',
                          `title`	VARCHAR(50) NOT NULL COMMENT '标题',
                          PRIMARY KEY (`id`) USING BTREE,
                          UNIQUE KEY `create_time_index` (`created_at`,`deleted_at`,`id`,`author_id`,`play_url`,`cover_url`,`favorite_count`,`comment_count`,`title`) USING BTREE,
                          KEY `idx_videos_deleted_at` (`deleted_at`) USING BTREE,
                          KEY `fk_videos_author` (`author_id`) USING BTREE,
                          CONSTRAINT `fk_videos_author` FOREIGN KEY (`author_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 ROW_FORMAT=DYNAMIC;


