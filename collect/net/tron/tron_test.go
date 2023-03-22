package tron

import (
	"github.com/sunjiangjun/xlog"
	"testing"
)

func TestEth_GetBlockByNumber(t *testing.T) {
	log := xlog.NewXLogger()
	log.Println(Eth_GetBlockByNumber("https://api.trongrid.io/jsonrpc", "244f918d-56b5-4a16-9665-9637598b1223", "0xF9CC56", log))
}
