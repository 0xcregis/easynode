package tron

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"testing"
)

func TestEth_GetBlockByNumber(t *testing.T) {
	log := xlog.NewXLogger()
	client := chainService.NewTron(map[int64][]*config.NodeCluster{205: {{NodeUrl: "https://api.trongrid.io", NodeToken: "244f918d-56b5-4a16-9665-9637598b1223"}}}, log)
	log.Println(Eth_GetBlockByNumber(client, "0xF9CC56", log))
}

func TestEth_GetTransactionByHash(t *testing.T) {
	log := xlog.NewXLogger()
	client := chainService.NewTron(map[int64][]*config.NodeCluster{205: {{NodeUrl: "https://api.trongrid.io", NodeToken: "244f918d-56b5-4a16-9665-9637598b1223"}}}, log)
	log.Println(Eth_GetTransactionByHash(client, "0x72fd440ff0542c2c28db762b4268f126c57f0fdf6daf69258cb9a306e26723e8", log))
}

func TestEth_GetTransactionReceiptByHash(t *testing.T) {
	log := xlog.NewXLogger()
	client := chainService.NewTron(map[int64][]*config.NodeCluster{205: {{NodeUrl: "https://api.trongrid.io", NodeToken: "244f918d-56b5-4a16-9665-9637598b1223"}}}, log)
	log.Println(Eth_GetTransactionReceiptByHash(client, "0x72fd440ff0542c2c28db762b4268f126c57f0fdf6daf69258cb9a306e26723e8", log))
}

func TestEth_GetCode(t *testing.T) {
	log := xlog.NewXLogger()
	client := chainService.NewTron(map[int64][]*config.NodeCluster{205: {{NodeUrl: "https://api.trongrid.io", NodeToken: "244f918d-56b5-4a16-9665-9637598b1223"}}}, log)
	log.Println(Eth_GetCode(client, 205, "41a614f803b6fd780986a42c78ec9c7f77e6ded13c"))
}
