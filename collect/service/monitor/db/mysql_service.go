package db

import (
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/common/pg"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"gorm.io/gorm"
	"strings"
	"time"
)

type MonitorDb struct {
	taskDb *gorm.DB
	config *config.Config
}

func (m MonitorDb) CheckTable() {
	//node_task
	tableName := fmt.Sprintf("%v_%v", m.config.TaskDb.Table, time.Now().Format(service.DateFormat))
	createSql := fmt.Sprintf(NodeTaskTable, m.config.TaskDb.DbName, m.config.TaskDb.DbName, tableName)
	sqlList := strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := m.taskDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//NodeTaskTable check
	var TaskNum int64
	err := m.taskDb.Raw("SELECT count(1) as task_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", m.config.TaskDb.DbName, tableName).Pluck("task_num", &TaskNum).Error
	if err != nil || TaskNum < 1 {
		panic("not found NodeTaskTable")
	}
}

func NewMonitorDbService(config *config.Config, xg *xlog.XLog) service.MonitorDbInterface {
	taskDb, err := pg.Open(config.TaskDb.User, config.TaskDb.Password, config.TaskDb.Addr, config.TaskDb.DbName, config.TaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &MonitorDb{
		taskDb: taskDb,
		config: config,
	}
}
