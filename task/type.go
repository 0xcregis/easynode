package task

import "time"

const (
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
	LogTime      time.Time `json:"log_time" gorm:"column:log_time"`
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
	Id          int64     `json:"id"  gorm:"column:id"`
	NodeId      string    `json:"nodeId" gorm:"column:node_id"`
	BlockNumber string    `json:"blockNumber" gorm:"column:block_number"`
	BlockHash   string    `json:"blockHash" gorm:"column:block_hash"`
	TxHash      string    `json:"txHash" gorm:"column:tx_hash"`
	TaskType    int       `json:"taskType" gorm:"column:task_type"` // 0:保留 1:同步Tx. 2:同步Block 3:同步Receipt
	BlockChain  int64     `json:"blockChain" gorm:"column:block_chain"`
	TaskStatus  int       `json:"taskStatus" gorm:"column:task_status"` //0: 初始 1: 成功. 2: 失败.  3: 执行中 其他：重试次数
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time"`
	LogTime     time.Time `json:"logTime" gorm:"column:log_time"`
}
