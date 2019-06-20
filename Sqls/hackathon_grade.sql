/*
 Navicat Premium Data Transfer

 Source Server         : dev-hls
 Source Server Type    : MySQL
 Source Server Version : 50636
 Source Host           : mysql-a.hfjy.com:3306
 Source Schema         : hls_test2

 Target Server Type    : MySQL
 Target Server Version : 50636
 File Encoding         : 65001

 Date: 29/05/2019 14:09:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hackathon_grade
-- ----------------------------
DROP TABLE IF EXISTS `hackathon_grade`;
CREATE TABLE `hackathon_grade`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '年级名称，一年级、二年级、三年级...',
  `isDel` tinyint(2) NOT NULL DEFAULT 0 COMMENT '是否删除 1：已删除； 0未删除；',
  `createTime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updateTime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `index_hackathon_grade_name`(`name`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 8 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '骇客马拉松—年级基础表' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of hackathon_grade
-- ----------------------------
INSERT INTO `hackathon_grade` VALUES (1, '一年级', 0, '2019-05-29 11:17:41', '2019-05-29 11:17:41');
INSERT INTO `hackathon_grade` VALUES (2, '二年级', 0, '2019-05-29 11:17:49', '2019-05-29 11:17:49');
INSERT INTO `hackathon_grade` VALUES (4, '初一', 0, '2019-05-29 11:18:15', '2019-05-29 11:20:07');
INSERT INTO `hackathon_grade` VALUES (5, '初二', 0, '2019-05-29 11:18:20', '2019-05-29 11:20:22');
INSERT INTO `hackathon_grade` VALUES (6, '高一', 0, '2019-05-29 11:18:35', '2019-05-29 11:20:41');
INSERT INTO `hackathon_grade` VALUES (7, '高二', 0, '2019-05-29 11:20:17', '2019-05-29 11:20:43');

SET FOREIGN_KEY_CHECKS = 1;
