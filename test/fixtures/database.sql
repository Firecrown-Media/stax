-- Test database dump for stax unit tests
-- MySQL dump 10.13  Distrib 8.0.33, for macos13 (arm64)

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS=0;

-- Table structure for table `wp_posts`
DROP TABLE IF EXISTS `wp_posts`;
CREATE TABLE `wp_posts` (
  `ID` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `post_author` bigint(20) unsigned NOT NULL DEFAULT '0',
  `post_date` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `post_content` longtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `post_title` text COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `post_status` varchar(20) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT 'publish',
  `post_name` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `post_type` varchar(20) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT 'post',
  `guid` varchar(255) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`ID`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_posts` VALUES
(1, 1, '2024-01-01 12:00:00', 'Welcome to WordPress. This is your first post.', 'Hello World', 'publish', 'hello-world', 'post', 'https://example.wpengine.com/?p=1');

-- Table structure for table `wp_options`
DROP TABLE IF EXISTS `wp_options`;
CREATE TABLE `wp_options` (
  `option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `option_name` varchar(191) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `option_value` longtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `autoload` varchar(20) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT 'yes',
  PRIMARY KEY (`option_id`),
  UNIQUE KEY `option_name` (`option_name`)
) ENGINE=InnoDB AUTO_INCREMENT=100 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_options` VALUES
(1, 'siteurl', 'https://example.wpengine.com', 'yes'),
(2, 'home', 'https://example.wpengine.com', 'yes'),
(3, 'blogname', 'Test Network', 'yes'),
(4, 'blogdescription', 'Just another WordPress site', 'yes'),
(5, 'admin_email', 'admin@example.wpengine.com', 'yes'),
(6, 'active_plugins', 'a:0:{}', 'yes'),
(7, 'template', 'twentytwentyfour', 'yes'),
(8, 'stylesheet', 'twentytwentyfour', 'yes');

-- Table structure for table `wp_blogs`
DROP TABLE IF EXISTS `wp_blogs`;
CREATE TABLE `wp_blogs` (
  `blog_id` bigint(20) NOT NULL AUTO_INCREMENT,
  `site_id` bigint(20) NOT NULL DEFAULT '0',
  `domain` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `path` varchar(100) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `registered` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `last_updated` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',
  `public` tinyint(2) NOT NULL DEFAULT '1',
  `archived` tinyint(2) NOT NULL DEFAULT '0',
  `mature` tinyint(2) NOT NULL DEFAULT '0',
  `spam` tinyint(2) NOT NULL DEFAULT '0',
  `deleted` tinyint(2) NOT NULL DEFAULT '0',
  `lang_id` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`blog_id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_blogs` VALUES
(1, 1, 'example.wpengine.com', '/', '2024-01-01 12:00:00', '2024-01-01 12:00:00', 1, 0, 0, 0, 0, 0),
(2, 1, 'site1.wpengine.com', '/', '2024-01-01 12:00:00', '2024-01-01 12:00:00', 1, 0, 0, 0, 0, 0),
(3, 1, 'site2.wpengine.com', '/', '2024-01-01 12:00:00', '2024-01-01 12:00:00', 1, 0, 0, 0, 0, 0);

-- Table structure for table `wp_site`
DROP TABLE IF EXISTS `wp_site`;
CREATE TABLE `wp_site` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `domain` varchar(200) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `path` varchar(100) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_site` VALUES
(1, 'example.wpengine.com', '/');

-- Table structure for table `wp_2_options`
DROP TABLE IF EXISTS `wp_2_options`;
CREATE TABLE `wp_2_options` (
  `option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `option_name` varchar(191) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `option_value` longtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `autoload` varchar(20) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT 'yes',
  PRIMARY KEY (`option_id`),
  UNIQUE KEY `option_name` (`option_name`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_2_options` VALUES
(1, 'siteurl', 'https://site1.wpengine.com', 'yes'),
(2, 'home', 'https://site1.wpengine.com', 'yes'),
(3, 'blogname', 'Site 1', 'yes');

-- Table structure for table `wp_3_options`
DROP TABLE IF EXISTS `wp_3_options`;
CREATE TABLE `wp_3_options` (
  `option_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `option_name` varchar(191) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT '',
  `option_value` longtext COLLATE utf8mb4_unicode_520_ci NOT NULL,
  `autoload` varchar(20) COLLATE utf8mb4_unicode_520_ci NOT NULL DEFAULT 'yes',
  PRIMARY KEY (`option_id`),
  UNIQUE KEY `option_name` (`option_name`)
) ENGINE=InnoDB AUTO_INCREMENT=50 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci;

INSERT INTO `wp_3_options` VALUES
(1, 'siteurl', 'https://site2.wpengine.com', 'yes'),
(2, 'home', 'https://site2.wpengine.com', 'yes'),
(3, 'blogname', 'Site 2', 'yes');

SET FOREIGN_KEY_CHECKS=1;
