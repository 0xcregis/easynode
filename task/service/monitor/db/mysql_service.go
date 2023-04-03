package db

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/common/sql"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"gorm.io/gorm"
	"strings"
	"time"
)

type MonitorDB struct {
	config        *config.Config
	taskDb        *gorm.DB
	blockNumberDb *gorm.DB
	log           *xlog.XLog
}

func (m *MonitorDB) CreateNodeTaskTable() {
	for true {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 23, 50, 0, 0, now.Location())
		<-time.After(next.Sub(now))

		log := m.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "createNodeTaskTable",
		})
		//new next table
		createSql := "CREATE TABLE if NOT EXISTS `%v` (\n  `node_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '节点的唯一标识',\n  `block_number` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块高度',\n  `block_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '区块hash',\n  `tx_hash` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL COMMENT '交易hash',\n  `task_type` tinyint NOT NULL DEFAULT '0' COMMENT ' 0:保留 1:同步Tx 2:同步Block 3:同步收据',\n  `block_chain` int NOT NULL DEFAULT '100' COMMENT '公链code, 默认：100 (etc)',\n  `task_status` int DEFAULT '0' COMMENT '0: 初始 1: 成功. 2: 失败.  3: 执行中. 4:kafka 写入中 5:重试 其他：重试次数(5以上)',\n  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,\n  `log_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,\n  `id` bigint NOT NULL AUTO_INCREMENT,\n  PRIMARY KEY (`id`),\n  KEY `type` (`task_type`) USING BTREE,\n  KEY `status` (`task_status`) USING BTREE,\n  KEY `tx_hash` (`tx_hash`) USING BTREE,\n  KEY `block_number` (`block_number`) USING BTREE,\n  KEY `block_hash` (`block_hash`) USING BTREE, \n KEY `block_chain` (`block_chain`) USING BTREE,\n  KEY `node_id` (`node_id`) USING BTREE\n) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='节点任务表';"
		day := next.Add(5 * time.Hour).Format(service.DayFormat)
		pre := next.Add(-5 * time.Hour).Format(service.DayFormat)

		dayTable := fmt.Sprintf("%v_%v", m.config.NodeTaskDb.Table, day)
		preTable := fmt.Sprintf("%v_%v", m.config.NodeTaskDb.Table, pre)

		createSql = fmt.Sprintf(createSql, dayTable)
		err := m.taskDb.Exec(createSql).Error
		if err != nil {
			log.Errorf("task.exec|sql=%v,error=%v", createSql, err)
		}

		next = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		<-time.After(next.Sub(time.Now()))

		//cp data from pre table to current table
		cpSql := `INSERT IGNORE INTO %v (id, node_id, block_number, block_hash, tx_hash, task_type, block_chain, task_status )
				SELECT
				id,
				node_id,
				block_number,
				block_hash,
				tx_hash,
				task_type,
				block_chain,
				task_status 
				FROM %v where task_status in (0,2,3,4)`
		cpSql = fmt.Sprintf(cpSql, dayTable, preTable)

		err = m.taskDb.Exec(cpSql).Error
		if err != nil {
			log.Errorf("taskDb.Exec|sql=%v,error=%v", cpSql, err)
		}

		//delete pre table
		dropSql := "drop table %v"
		dropSql = fmt.Sprintf(dropSql, preTable)
		err = m.taskDb.Exec(dropSql).Error
		if err != nil {
			log.Printf("taskDb.Exec|sql=%v,error=%v", dropSql, err)
		}
		//delete binlog
		m.taskDb.Exec("RESET MASTER")
	}
}

func (m *MonitorDB) RetryTaskForFail() {
	log := m.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "RetryTaskForFail",
	})

	var ids []int64
	err := m.taskDb.Table(m.getNodeTaskTable()).Select("id").Where("task_status=?", 2).Pluck("id", &ids).Error
	if err != nil {
		log.Errorf("taskDb|err=%v", err)
		return
	}

	if len(ids) > 0 {
		str2 := `UPDATE %v SET task_status=5 WHERE task_status=2 and id in (?)`
		str2 = fmt.Sprintf(str2, m.getNodeTaskTable())
		err = m.taskDb.Exec(str2, ids).Error
		if err != nil {
			log.Errorf("taskDb|sql=%v,error=%v", str2, err)
		}
	}
}

func (m *MonitorDB) HandlerDeadTask() {
	log := m.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "HandlerDeadTask",
	})
	//长时间"正在执行"的任务，改成失败状态
	err := m.taskDb.Table(m.getNodeTaskTable()).Where("task_status in (?) and create_time<?", []int{3, 4}, time.Now().Add(-30*time.Minute).UTC().Format("2006-01-02 15:04:05")).UpdateColumn("task_status", 2).Error
	if err != nil {
		log.Errorf("taskDb|update|err=%v", err.Error())
	}
}

func (m *MonitorDB) CheckTable() {
	//node_task
	tableName := fmt.Sprintf("%v_%v", m.config.NodeTaskDb.Table, time.Now().Format(service.DayFormat))
	createSql := fmt.Sprintf(NodeTaskTable, m.config.NodeTaskDb.DbName, m.config.NodeTaskDb.DbName, tableName)
	sqlList := strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := m.taskDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//block_number
	createSql = fmt.Sprintf(BlockNumberTable, m.config.BlockNumberDb.DbName, m.config.BlockNumberDb.DbName, m.config.BlockNumberDb.Table)
	sqlList = strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := m.blockNumberDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//NodeTaskTable check
	var TaskNum int64
	err := m.taskDb.Raw("SELECT count(1) as task_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", m.config.NodeTaskDb.DbName, tableName).Pluck("task_num", &TaskNum).Error
	if err != nil || TaskNum < 1 {
		panic("not found NodeTaskTable")
	}

	//blockNumberTable check
	var blockNum int64
	err = m.blockNumberDb.Raw("SELECT count(1) as block_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", m.config.BlockNumberDb.DbName, m.config.BlockNumberDb.Table).Pluck("block_num", &blockNum).Error
	if err != nil || blockNum < 1 {
		panic("not found BlockNumberTable")
	}

}

func (m *MonitorDB) getNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", m.config.NodeTaskDb.Table, time.Now().Format(service.DayFormat))
	return table
}

func (m *MonitorDB) HandlerManyFailTask() {
	log := m.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "HandlerManyFailTask",
	})

	//如果任务多次重试，仍然失败，则放弃
	str := `SELECT block_chain, block_number,block_hash,tx_hash,task_type,count(1) as c FROM %v WHERE task_status in (2,5) GROUP BY block_chain, block_number,block_hash,tx_hash,task_type HAVING c>?`
	str = fmt.Sprintf(str, m.getNodeTaskTable())
	var list []*service.NodeTask
	err := m.taskDb.Raw(str, 5).Scan(&list).Error
	if err != nil {
		log.Printf("taskDb|raw|sql=%v,err=%v", str, err)
		return
	}

	for _, v := range list {
		err := m.taskDb.Table(m.getNodeTaskTable()).Where("block_chain=? and block_number=? and block_hash=? and tx_hash=? and task_type=?", v.BlockChain, v.BlockNumber, v.BlockHash, v.TxHash, v.TaskType).UpdateColumn("task_status", 5).Error
		if err != nil {
			log.Errorf("taskDb|update|err=%v", err)
			continue
		}
	}
}

func NewMySQLMonitorService(config *config.Config, log *xlog.XLog) service.DbTaskMonitorInterface {
	task, err := sql.Open(config.NodeTaskDb.User, config.NodeTaskDb.Password, config.NodeTaskDb.Addr, config.NodeTaskDb.DbName, config.NodeTaskDb.Port, log)
	if err != nil {
		panic(err)
	}

	blockNumber, err := sql.Open(config.BlockNumberDb.User, config.BlockNumberDb.Password, config.BlockNumberDb.Addr, config.BlockNumberDb.DbName, config.BlockNumberDb.Port, log)
	if err != nil {
		panic(err)
	}

	return &MonitorDB{
		config:        config,
		taskDb:        task,
		blockNumberDb: blockNumber,
		log:           log,
	}
}
