package service

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	"log"
	"testing"
)

func Init2() API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_tron.json")
	return NewTron(cfg.Cluster[205], xlog.NewXLogger())
}

func TestTron_GetBlockByNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockByNumber(205, "49477110", false))
}

func TestTron_GetBlockByHash(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockByHash(205, "0000000002f2f5f62d94d85ec1abf2c0dfc26d72da4f5e5d5a2624d51e231425", false))
}

func TestTron_GetBlockReceiptByBlockNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockReceiptByBlockNumber(205, "34222872"))
}
