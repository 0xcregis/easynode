package monitor

import (
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/common/util"
	"path"
	"time"
)

type Service struct {
	logConfig config.LogConfig
	log       *xlog.XLog
}

func (s *Service) Start() {

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

func NewService(logConfig *config.LogConfig, xg *xlog.XLog) *Service {
	return &Service{
		logConfig: *logConfig,
		log:       xg,
	}
}
