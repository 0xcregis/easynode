package tron

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init2() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_tron.json")
	return NewTron(cfg.Cluster[205], 205, xlog.NewXLogger())
}

func TestTron_Balance(t *testing.T) {
	s := Init2()
	resp, err := s.Balance(205, "TXeZAknJe2gbqSJyZYXbNMVvQgsKQbSoxX", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetLatestBlock(t *testing.T) {
	s := Init2()
	resp, err := s.LatestBlock(205)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetTxByHash(t *testing.T) {
	s := Init2()
	resp, err := s.GetTxByHash(205, "d0ff91487dd11ab6bd2cffa4af97bb472ede4f1713786fa2b15bf32011d0b681")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetBlockByNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockByNumber(205, "45611899", false))
}

func TestTron_GetBlockByHash(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockByHash(205, "0000000002f2f5f62d94d85ec1abf2c0dfc26d72da4f5e5d5a2624d51e231425", false))
}

func TestTron_GetBlockReceiptByBlockNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockReceiptByBlockNumber(205, "45611899"))
}
