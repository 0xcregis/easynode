package ether

import (
	"log"
	"testing"
)

func TestEth_GetBlockNumber(t *testing.T) {
	log.Println(Eth_GetBlockNumber("https://eth-mainnet.g.alchemy.com/v2", "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop"))
}
