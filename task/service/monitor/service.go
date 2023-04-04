package monitor

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/monitor/db"
	"time"
)

/**
  1. 判断任务长时间处于 task_status=3,则 直接改成2（失败）
  2. 如果一个任务 多次失败，则不在重试
*/

type Service struct {
	config    *config.Config
	dbMonitor service.DbTaskMonitorInterface
	log       *xlog.XLog
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/monitor_task", 24*time.Hour)

	db := db.NewMySQLMonitorService(config, xg)

	return &Service{
		config:    config,
		dbMonitor: db,
		log:       xg,
	}
}

func (s *Service) Start() {

	//检查数据表 是否完备
	s.dbMonitor.CheckTable()

	//每日分表
	go s.dbMonitor.CreateNodeTaskTable()

	//定时处理异常数据
	go func() {
		for true {
			<-time.After(30 * time.Minute)
			//处理僵死任务：即 长期处理进行中 status=3
			s.dbMonitor.HandlerDeadTask()

			//处理任务失败多次的
			s.dbMonitor.HandlerManyFailTask()

			//失败任务重试
			s.dbMonitor.RetryTaskForFail()
		}
	}()
}
