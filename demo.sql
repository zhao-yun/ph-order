/*
 Navicat MySQL Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 50744
 Source Host           : localhost:3306
 Source Schema         : demo

 Target Server Type    : MySQL
 Target Server Version : 50744
 File Encoding         : 65001

 Date: 23/03/2025 12:58:08
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for test_data
-- ----------------------------
DROP TABLE IF EXISTS `test_data`;
CREATE TABLE `test_data` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `num` int(11) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of test_data
-- ----------------------------
BEGIN;
INSERT INTO `test_data` (`id`, `name`, `num`) VALUES (1, '小明', 10);
INSERT INTO `test_data` (`id`, `name`, `num`) VALUES (2, '小王', 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
