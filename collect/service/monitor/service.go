package monitor

import (
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/common/pg"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/util"
	"gorm.io/gorm"
	"path"
	"strings"
	"time"
)

type Service struct {
	config       *config.Config
	logConfig    config.LogConfig
	nodeSourceDb *gorm.DB
	taskDb       *gorm.DB
	nodeInfoDb   *gorm.DB
}

func (s *Service) CheckTable() {

	//node_task
	tableName := fmt.Sprintf("%v_%v", s.config.TaskDb.Table, time.Now().Format(service.DateFormat))
	createSql := fmt.Sprintf(NodeTaskTable, s.config.TaskDb.DbName, s.config.TaskDb.DbName, tableName)
	sqlList := strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := s.taskDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//node_info
	createSql = fmt.Sprintf(NodeInfoTable, s.config.NodeInfoDb.DbName, s.config.NodeInfoDb.DbName, s.config.NodeInfoDb.Table)
	sqlList = strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := s.taskDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//node_source
	createSql = fmt.Sprintf(NodeSourceTable, s.config.SourceDb.DbName, s.config.SourceDb.DbName, s.config.SourceDb.Table)
	sqlList = strings.Split(createSql, ";")
	for _, sql := range sqlList {
		err := s.taskDb.Exec(sql).Error
		if err != nil {
			panic(err)
		}
	}

	//NodeTaskTable check
	var TaskNum int64
	err := s.taskDb.Raw("SELECT count(1) as task_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", s.config.TaskDb.DbName, tableName).Pluck("task_num", &TaskNum).Error
	if err != nil || TaskNum < 1 {
		panic("not found NodeTaskTable")
	}

	//NodeSourceTable check
	var SourceNum int64
	err = s.nodeSourceDb.Raw("SELECT count(1) as source_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", s.config.SourceDb.DbName, s.config.SourceDb.Table).Pluck("source_num", &SourceNum).Error
	if err != nil || SourceNum < 1 {
		panic("not found NodeSourceTable")
	}

	//NodeInfoTable check
	var InfoNum int64
	err = s.nodeInfoDb.Raw("SELECT count(1) as info_num FROM information_schema.`TABLES` WHERE TABLE_SCHEMA=? and TABLE_NAME=?", s.config.NodeInfoDb.DbName, s.config.NodeInfoDb.Table).Pluck("info_num", &InfoNum).Error
	if err != nil || InfoNum < 1 {
		panic("not found NodeInfoTable")
	}
}

func (s *Service) Start() {

	//检查表
	s.CheckTable()

	//监控服务
	go func() {

		for true {
			<-time.After(7 * time.Hour)
			p := s.logConfig.Path
			d := s.logConfig.Delay

			h := time.Duration(d*24) * time.Hour
			t := time.Now().Add(-h)

			for i := 0; i < 5; i++ {
				datePath := t.Format(service.DateFormat)

				datePath = fmt.Sprintf("%v%v", datePath, "0000")
				cmdLog := fmt.Sprintf("%v_%v", "cmd_log", datePath)
				_ = util.DeleteFile(path.Join(p, cmdLog))

				//
				nodeInfoLog := fmt.Sprintf("%v_%v", "node_info_log", datePath)
				_ = util.DeleteFile(path.Join(p, nodeInfoLog))

				chainInfoLog := fmt.Sprintf("%v_%v", "chain_info_log", datePath)
				_ = util.DeleteFile(path.Join(p, chainInfoLog))

				t = t.Add(-24 * time.Hour)
			}
		}

	}()

}

func (s *Service) Stop() {
	panic("implement me")
}

func NewService(config *config.Config, logConfig *config.LogConfig, xg *xlog.XLog) *Service {

	s, err := pg.Open(config.SourceDb.User, config.SourceDb.Password, config.SourceDb.Addr, config.SourceDb.DbName, config.SourceDb.Port, xg)
	if err != nil {
		panic(err)
	}

	info, err := pg.Open(config.NodeInfoDb.User, config.NodeInfoDb.Password, config.NodeInfoDb.Addr, config.NodeInfoDb.DbName, config.NodeInfoDb.Port, xg)
	if err != nil {
		panic(err)
	}

	task, err := pg.Open(config.TaskDb.User, config.TaskDb.Password, config.TaskDb.Addr, config.TaskDb.DbName, config.TaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &Service{
		config:       config,
		logConfig:    *logConfig,
		nodeSourceDb: s,
		nodeInfoDb:   info,
		taskDb:       task,
	}
}
