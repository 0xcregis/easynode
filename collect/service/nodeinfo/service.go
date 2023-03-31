package nodeinfo

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/common/pg"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Service struct {
	log    *xlog.XLog
	orm    *gorm.DB
	cfg    *config.NodeInfoDb
	chains []*config.Chain
}

func (s *Service) Stop() {
	//TODO implement me
	panic("implement me")
}

func (s *Service) Start() {

	go func() {
		for true {
			log := s.log.WithFields(logrus.Fields{
				"id":    time.Now().UnixMilli(),
				"model": "Start",
			})
			cs, _ := json.Marshal(s.chains)
			node := &service.NodeInfo{
				NodeId:     util.GetLocalNodeId(),
				Info:       string(cs),
				Ip:         util.GetIPs()[0],
				Host:       util.GetMacAddrs()[0],
				CreateTime: time.Now(),
			}
			s.PushNodeInfo(node, log)
			time.Sleep(10 * time.Second)
		}
	}()

}

func (s *Service) PushNodeInfo(node *service.NodeInfo, log *logrus.Entry) {
	err := s.orm.Table(s.cfg.Table).Omit("id").Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"ip", "host", "info", "create_time"})}).Create(node).Error
	if err != nil {
		log.Error("PushNodeInfo|error", err.Error())
	}
}

func NewService(db *config.NodeInfoDb, chains []*config.Chain, logConfig *config.LogConfig, x *xlog.XLog) *Service {
	//x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/node_info", logConfig.Path), 24*time.Hour)
	g, err := pg.Open(db.User, db.Password, db.Addr, db.DbName, db.Port, x)
	if err != nil {
		panic(err)
	}
	return &Service{
		log:    x,
		orm:    g,
		cfg:    db,
		chains: chains,
	}
}
