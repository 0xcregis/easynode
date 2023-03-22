package service

import "time"

const (
	BLOCK_NUMBER_TABLE = "block_number"

	DayFormat  = "20060102"
	TimeFormat = "2006-01-02 15:04:05"
)

/**
  CREATE TABLE `block_number` (
  `id` int NOT NULL AUTO_INCREMENT,
  `chain_code` int NOT NULL DEFAULT '100' COMMENT '公链代码',
  `recent_number` bigint DEFAULT '0' COMMENT '区块链的最新高度',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`chain_code`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='公链最新区块高度表';
*/

type BlockNumber struct {
	Id           int64     `json:"id" gorm:"column:id"`
	ChainCode    int64     `json:"chain_code" gorm:"column:chain_code"`       // 公链代码
	RecentNumber int64     `json:"recent_number" gorm:"column:recent_number"` // 已处理的高度
	LatestNumber int64     `json:"latest_number" gorm:"column:latest_number"` // 区块链的最新高度
	CreateTime   time.Time `json:"create_time" gorm:"column:create_time"`
	LogTime      time.Time `json:"log_time" gorm:"column:log_time"`
}

/**

CREATE TABLE `node_source` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `block_chain` int NOT NULL COMMENT '公链code',
  `tx_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `block_number` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `source_type` tinyint DEFAULT '0' COMMENT '任务类型\n1: 交易\n2:区块\n3.收据',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1383 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='任务来源表\n\n1. 自产生 block number\n2.重试任务: Tx, block\n3. 收据任务';

*/

type NodeSource struct {
	Id          int64     `json:"id" gorm:"column:id"`
	BlockChain  int64     `json:"block_chain" gorm:"column:block_chain"` // 公链code
	TxHash      string    `json:"tx_hash" gorm:"column:tx_hash"`
	BlockHash   string    `json:"block_hash" gorm:"column:block_hash"`
	BlockNumber string    `json:"block_number" gorm:"column:block_number"`
	SourceType  int8      `json:"source_type" gorm:"column:source_type"` // 任务类型 1: 交易 2:区块 3.收据
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
}

/**
CREATE TABLE `node_info` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `ip` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'ip',
  `host` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT 'host',
  `node_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点的唯一标识',
  `info` json NOT NULL COMMENT '''采集节点的配置信息 ，一个json 对象 如：\n[\n  {\n   "blockchain":"eth", \n   "block_code":100,\n   "config":{\n              "tx": 0,\n               "block":1\n          }\n  },\n  {}\n]'';',
  `create_time` datetime DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `node_id` (`node_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点信息表';

*/

type NodeInfo struct {
	Id         int64     `json:"id" gorm:"column:id"`
	Ip         string    `json:"ip" gorm:"column:ip"`           // ip
	Host       string    `json:"host" gorm:"column:host"`       // host
	NodeId     string    `json:"node_id" gorm:"column:node_id"` // 节点的唯一标识
	Info       string    `json:"info" gorm:"column:info"`       // 采集节点的配置信息 ，一个json 对象 如：
	CreateTime time.Time `json:"create_time" gorm:"column:create_time"`
}

/**
CREATE TABLE `node_task` (
  `node_id` varchar(255) COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点的唯一标识',
  `block_number` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块高度',
  `block_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块hash',
  `tx_hash` varchar(255) COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '交易hash',
  `task_type` tinyint NOT NULL DEFAULT '0' COMMENT ' 0:保留 1:同步Tx 2:同步Block 3:同步收据',
  `block_chain` int NOT NULL DEFAULT '100' COMMENT '公链code, 默认：100 (etc)',
  `task_status` int DEFAULT '0' COMMENT '0: 初始 1: 成功. 2: 失败.  3: 执行中 其他：重试次数',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `id` bigint NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点任务表';

*/

type NodeTask struct {
	NodeId      string    `json:"node_id" gorm:"column:node_id"`           // 节点的唯一标识
	BlockNumber string    `json:"block_number" gorm:"column:block_number"` // 区块高度
	BlockHash   string    `json:"block_hash" gorm:"column:block_hash"`     // 区块hash
	TxHash      string    `json:"tx_hash" gorm:"column:tx_hash"`           // 交易hash
	TaskType    int8      `json:"task_type" gorm:"column:task_type"`       //  0:保留 1:同步Tx 2:同步Block 3:同步收据
	BlockChain  int64     `json:"block_chain" gorm:"column:block_chain"`   // 公链code, 默认：100 (etc)
	TaskStatus  int       `json:"task_status" gorm:"column:task_status"`   // 0: 初始 1: 成功. 2: 失败.  3: 执行中 其他：重试次数
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	LogTime     time.Time `json:"log_time" gorm:"column:log_time"`
	Id          int64     `json:"id" gorm:"column:id"`
}
