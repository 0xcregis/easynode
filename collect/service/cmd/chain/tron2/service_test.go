package tron2

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_tron.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), cfg, x
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("76f1ff8be6b3cf041f29b67c3a5d025f232d2a48a6d0810f0f234fc73c16adcc", x.WithFields(logrus.Fields{}))
	log.Printf("%+v\n", tx)
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, t1 := s.GetBlockByNumber("52642923", x.WithFields(logrus.Fields{}), false)
	log.Println(b)
	log.Println(len(t1))
}

func TestService_GetReceiptByBlock(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceiptByBlock("", "0xF9CC56", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(r)
	}
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("0x72fd440ff0542c2c28db762b4268f126c57f0fdf6daf69258cb9a306e26723e8", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(r)
	}
}
