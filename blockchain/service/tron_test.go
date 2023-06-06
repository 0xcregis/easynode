package service

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	"log"
	"testing"
)

func Init2() API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_tron.json")
	return NewTron(cfg.Cluster[200], xlog.NewXLogger())
}

func TestTron_GetBlockByNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockByNumber(205, "49477110", true))
}

func TestTron_GetBlockReceiptByBlockNumber(t *testing.T) {
	c := Init2()
	log.Println(c.GetBlockReceiptByBlockNumber(205, "34222872"))
}
