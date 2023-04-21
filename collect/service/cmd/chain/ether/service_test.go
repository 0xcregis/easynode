package ether

import (
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"log"
	"testing"
)

func Init() (service.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config.json")
	x := xlog.NewXLogger()
	return NewService(cfg.Chains[0], x), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, cfg, x := Init()
	b, _ := s.GetBlockByNumber("16939002", cfg.Chains[0].BlockTask, x.WithFields(logrus.Fields{}))
	log.Printf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, cfg, x := Init()
	tx := s.GetTx("0xf31a17a18d5360e60a7b37c00de3286a3e0ebee43b9d9ba2009ad44f6536c323", cfg.Chains[0].TxTask, x.WithFields(logrus.Fields{}))
	log.Println(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, cfg, x := Init()
	r := s.GetReceipt("0xf31a17a18d5360e60a7b37c00de3286a3e0ebee43b9d9ba2009ad44f6536c323", cfg.Chains[0].ReceiptTask, x.WithFields(logrus.Fields{}))
	log.Printf("%+v", r)
}

func TestService_BalanceCluster(t *testing.T) {
	s, cfg, _ := Init()
	c, err := s.BalanceCluster("", cfg.Chains[0].BlockTask.FromCluster)
	if err != nil {
		panic(err)
	}
	log.Println(c)
}

func TestService_Monitor(t *testing.T) {
	_, cfg, _ := Init()

	c1 := cfg.Chains[0]
	c2 := c1.CopyChain()

	c1.BlockTask.FromCluster[0].ErrorCount = 11
	c2.BlockTask.FromCluster[0].ErrorCount = 12

	c1.TxTask.FromCluster[0].ErrorCount = 13
	c2.TxTask.FromCluster[0].ErrorCount = 14

	c1.ReceiptTask.FromCluster[0].ErrorCount = 15
	c2.ReceiptTask.FromCluster[0].ErrorCount = 16

	c1.BlockChainCode = 101
	c2.BlockChainCode = 102

	log.Println(c2)
	log.Println(c1)
}
