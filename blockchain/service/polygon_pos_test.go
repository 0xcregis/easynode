package service

import (
	"log"
	"testing"
	"time"

	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init4() API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_polygon.json")
	return NewPolygonPos(cfg.Cluster[201], xlog.NewXLogger())
}

func TestPolygonPos_GetBlockByNumber(t *testing.T) {
	s := Init4()
	log.Println(s.GetBlockByNumber(201, "0xF3F088", false))
}

func TestPolygonPos_GetBlockByHash(t *testing.T) {
	s := Init4()
	log.Println(s.GetBlockByHash(201, "0x51fb993ede72e8a3428f4782435724d07ae1df5d3961e278753980ff7e8b9846", false))

	time.Sleep(10 * time.Second)
}
