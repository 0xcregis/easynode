package service

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	"log"
	"testing"
)

func Init() *Tron2 {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_tron.json")
	return NewTron2(cfg.Cluster, xlog.NewXLogger())
}

func TestTron2_GetTransactionReceiptByHash(t *testing.T) {
	c := Init()
	log.Println(c.GetTransactionReceiptByHash(205, "568a42d70abf652610083266b3044b8d753c8610746c82e69109839f19ed63b0"))
}

func TestTron2_GetBlockByNumber(t *testing.T) {
	c := Init()
	log.Println(c.GetBlockByNumber(205, "49477110"))
}
