/*
 Navicat MySQL Data Transfer

 Source Server         : easy_node
 Source Server Type    : MySQL
 Source Server Version : 80031
 Source Host           : 192.168.2.11:3306
 Source Schema         : easy_node

 Target Server Type    : MySQL
 Target Server Version : 80031
 File Encoding         : 65001

 Date: 13/02/2023 18:23:53
*/

CREATE DATABASE easynode;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for block_number
-- ----------------------------
DROP TABLE IF EXISTS `block_number`;
CREATE TABLE `block_number` (
  `id` int NOT NULL AUTO_INCREMENT,
  `chain_code` int NOT NULL DEFAULT '100' COMMENT '公链代码',
  `latest_number` bigint DEFAULT '0' COMMENT '区块链最新高度',
  `recent_number` bigint DEFAULT '0' COMMENT '已经处理区块高度',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`chain_code`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=1588893 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公链最新区块高度表';

-- ----------------------------
-- Table structure for node_error
-- ----------------------------
DROP TABLE IF EXISTS `node_error`;
CREATE TABLE `node_error` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block_chain` int NOT NULL COMMENT '公链code',
  `tx_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `source_type` tinyint DEFAULT '0' COMMENT '任务类型 1: 交易 2:区块 3.收据',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`block_chain`,`tx_hash`,`block_hash`,`block_number`,`source_type`) USING BTREE,
  KEY `type` (`source_type`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='系统自检的异常数据表';

-- ----------------------------
-- Table structure for node_info
-- ----------------------------
DROP TABLE IF EXISTS `node_info`;
CREATE TABLE `node_info` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `ip` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'ip',
  `host` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'host',
  `node_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点的唯一标识',
  `info` json NOT NULL COMMENT '''采集节点的配置信息 ，一个json 对象 如：\n[\n  {\n   "blockchain":"eth", \n   "block_code":100,\n   "config":{\n              "tx": 0,\n               "block":1\n          }\n  },\n  {}\n]'';',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `node_id` (`node_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=3936448 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点信息表';

-- ----------------------------
-- Table structure for node_source
-- ----------------------------
DROP TABLE IF EXISTS `node_source`;
CREATE TABLE `node_source` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block_chain` int NOT NULL COMMENT '公链code',
  `tx_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `source_type` tinyint DEFAULT '0' COMMENT '任务类型 1: 交易 2:区块 3.收据',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id` (`block_chain`,`tx_hash`,`block_hash`,`block_number`,`source_type`) USING BTREE,
  KEY `type` (`source_type`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2750196 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务来源表\n\n1. 自产生 block number\n2.重试任务: Tx, block\n3. 收据任务';

-- ----------------------------
-- Table structure for node_task_20230213
-- ----------------------------
DROP TABLE IF EXISTS `node_task_20230213`;
CREATE TABLE `node_task_20230213` (
  `node_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点的唯一标识',
  `block_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块高度',
  `block_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块hash',
  `tx_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '交易hash',
  `task_type` tinyint NOT NULL DEFAULT '0' COMMENT ' 0:保留 1:同步Tx 2:同步Block 3:同步收据',
  `block_chain` int NOT NULL DEFAULT '100' COMMENT '公链code, 默认：100 (etc)',
  `task_status` int DEFAULT '0' COMMENT '0: 初始 1: 成功. 2: 失败.  3: 执行中. 4:kafka 写入中 5:重试 其他：重试次数(5以上)',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `id` bigint NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `type` (`task_type`) USING BTREE,
  KEY `status` (`task_status`) USING BTREE,
  KEY `tx_hash` (`tx_hash`) USING BTREE,
  KEY `block_number` (`block_number`) USING BTREE,
  KEY `block_hash` (`block_hash`) USING BTREE,
  KEY `block_chain` (`block_chain`) USING BTREE,
  KEY `node_id` (`node_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=173770503 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点任务表';

SET FOREIGN_KEY_CHECKS = 1;
