package bnb

import (
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_bnb.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef", "0xc3d58168c5ae7397731d063d5bbf3d657854427343f4c083240f7aacaa2d0f62"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("32504450", x.WithFields(logrus.Fields{}), false)
	t.Logf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("0x4f4ee0761b29670f2261d5e7632f1215c9807f0e7fbec6a151e1201b97546383", x.WithFields(logrus.Fields{}))
	t.Log(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("0x4f4ee0761b29670f2261d5e7632f1215c9807f0e7fbec6a151e1201b97546383", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}

//func TestService_GetReceiptByBlock(t *testing.T) {
//	s, _, x := Init()
//	r, err := s.GetReceiptByBlock("", "17658423", x.WithFields(logrus.Fields{}))
//	if err != nil {
//		t.Error(err)
//	} else {
//		t.Logf("%+v", r)
//	}
//}
