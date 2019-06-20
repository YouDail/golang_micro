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

 Date: 29/05/2019 14:09:46
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hackathon_class
-- ----------------------------
DROP TABLE IF EXISTS `hackathon_class`;
CREATE TABLE `hackathon_class`  (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL DEFAULT '' COMMENT '班级名，（1）班',
  `gradeId` int(11) NOT NULL DEFAULT 0 COMMENT '年级',
  `isDel` tinyint(2) NOT NULL DEFAULT 0 COMMENT '是否删除 1：已删除； 0未删除；',
  `createTime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) COMMENT '创建时间',
  `updateTime` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP(0) ON UPDATE CURRENT_TIMESTAMP(0) COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `index_hackathon_class_gradeId`(`gradeId`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 50 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci COMMENT = '骇客马拉松—班级基础表' ROW_FORMAT = Compact;

-- ----------------------------
-- Records of hackathon_class
-- ----------------------------
INSERT INTO `hackathon_class` VALUES (1, '(1)班', 1, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (2, '(2)班', 1, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (3, '(3)班', 1, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (6, '(6)班', 1, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (7, '(7)班', 1, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (8, '(1)班', 2, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (10, '(3)班', 2, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (11, '(4)班', 2, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (12, '(5)班', 2, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (14, '(7)班', 2, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (22, '(1)班', 4, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (23, '(2)班', 4, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (25, '(4)班', 4, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (26, '(5)班', 4, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (28, '(7)班', 4, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (29, '(1)班', 5, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (30, '(2)班', 5, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (31, '(3)班', 5, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (32, '(4)班', 5, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (33, '(5)班', 5, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (36, '(1)班', 6, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (37, '(2)班', 6, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (40, '(5)班', 6, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (41, '(6)班', 6, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (42, '(7)班', 6, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (43, '(1)班', 7, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (44, '(2)班', 7, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (45, '(3)班', 7, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (46, '(4)班', 7, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');
INSERT INTO `hackathon_class` VALUES (49, '(7)班', 7, 0, '2019-05-29 11:38:52', '2019-05-29 11:38:52');

SET FOREIGN_KEY_CHECKS = 1;
