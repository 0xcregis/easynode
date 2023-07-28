package polygonpos

import (
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_polygon.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("45611899", x.WithFields(logrus.Fields{}), false)
	t.Logf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("0x9f656ad21cad7853f58aa05191ec4c11bd0459f40bec1a259f089fce4c80232f", x.WithFields(logrus.Fields{}))
	t.Log(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("0x9f656ad21cad7853f58aa05191ec4c11bd0459f40bec1a259f089fce4c80232f", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}

//func TestService_GetReceiptByBlock(t *testing.T) {
//	s, _, x := Init()
//	r, err := s.GetReceiptByBlock("", "45611899", x.WithFields(logrus.Fields{}))
//	if err != nil {
//		t.Error(err)
//	} else {
//		t.Logf("%+v", r)
//	}
//}
