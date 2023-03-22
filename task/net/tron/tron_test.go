package tron

import (
	"log"
	"testing"
)

func TestEth_GetBlockNumber(t *testing.T) {
	log.Println(Eth_GetBlockNumber("https://api.trongrid.io/jsonrpc", "244f918d-56b5-4a16-9665-9637598b1223"))
}
