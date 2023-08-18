package tron2

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/db"
)

func Init() (service.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_tron.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f"), cfg, x
}

func TestService_GetTx(t *testing.T) {

	s, _, x := Init()

	tx := s.GetTx("447a8ebba389d27c25079843779690d2f695b4d7f28188515979c9a80b337f4e", x.WithFields(logrus.Fields{}))

	log.Printf("%+v\n", tx)

	bs, _ := json.Marshal(tx)
	log.Println(string(bs))
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, t1 := s.GetBlockByNumber("49469984", x.WithFields(logrus.Fields{}), true)
	log.Println(b)
	log.Println(t1[0])
}

func TestService_GetReceiptByBlock(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceiptByBlock("", "0xF9CC56", x.WithFields(logrus.Fields{}))
	log.Println(r, err)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("50e6dd05c37b8666cf4a689fe6c0d52053b76b53d8649b256e6b9dca8c9df098", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	}
	log.Printf("%+v", r)
}
