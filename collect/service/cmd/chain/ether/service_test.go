package ether

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
	cfg := config.LoadConfig("./../../../../../cmd/collect/config.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("17658423", x.WithFields(logrus.Fields{}), false)
	log.Printf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("0xb7959ab3ed0f8b424f897ba8ec28168358079be0030a9337d20252f7d4c18cfd", x.WithFields(logrus.Fields{}))
	log.Println(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("0x2eadb796722102575e8a1fe75241cba96f605659c160a65b3ad9473027b74c21", x.WithFields(logrus.Fields{}))
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", r)
}

func TestService_GetReceiptByBlock(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceiptByBlock("", "17658423", x.WithFields(logrus.Fields{}))
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", r)
}
