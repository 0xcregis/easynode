package btc

import (
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_btc.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("808825", x.WithFields(logrus.Fields{}), true)
	t.Logf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("66acad1ac15fe6ff08e100a7ff02aa06e425d06f36a769be6065565aa0d52bea", x.WithFields(logrus.Fields{}))
	t.Log(tx)
}
