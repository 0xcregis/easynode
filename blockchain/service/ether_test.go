package service

import (
	"log"
	"testing"
	"time"

	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init3() API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config.json")
	return NewEth(cfg.Cluster[200], xlog.NewXLogger())
}

func TestEther_GetBlockByNumber(t *testing.T) {
	s := Init3()
	log.Println(s.GetBlockByNumber(200, "0xF3F088", false))
}

func TestEther_GetBlockByHash(t *testing.T) {
	s := Init3()
	log.Println(s.GetBlockByHash(200, "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f", false))

	time.Sleep(10 * time.Second)
}
