SET NAMES utf8mb4;
CREATE DATABASE IF NOT EXISTS todolist COLLATE utf8mb4_unicode_ci;

USE todolist;

CREATE TABLE IF NOT EXISTS `user` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT 'Primary Key ID',
  `name` varchar(128) NOT NULL DEFAULT '' COMMENT 'User Nickname',
  `password` varchar(128) NOT NULL DEFAULT '' COMMENT 'Password (Encrypted)',
  `icon_uri` varchar(512) NOT NULL DEFAULT '' COMMENT 'Avatar URI',
  `created_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Creation Time (Milliseconds)',
  `updated_at` bigint unsigned NOT NULL DEFAULT 0 COMMENT 'Update Time (Milliseconds)',
  `deleted_at` bigint unsigned NULL COMMENT 'Deletion Time (Milliseconds)',
  PRIMARY KEY (`id`),
  UNIQUE INDEX `uniq_unique_name` (`name`)
) ENGINE=InnoDB CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT 'User Table';

CREATE TABLE IF NOT EXISTS `task` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'Task ID',
  `user_id` bigint NOT NULL COMMENT 'Task OwnerID',
  `title` varchar(255) NOT NULL COMMENT 'Task Title',
  `content` text NOT NULL COMMENT 'Task Content',
  `status` tinyint NOT NULL COMMENT 'Task Status',
  `created_at` bigint NOT NULL COMMENT 'Creation Time (Milliseconds)',
  `updated_at` bigint NOT NULL COMMENT 'Update Time (Milliseconds)',
  PRIMARY KEY (`id`),
  INDEX idx_user_status_utime (`user_id`, `status`, `updated_at`)
) ENGINE=InnoDB CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT 'Task Table';