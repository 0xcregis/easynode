package ether

import (
	"log"
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalance("https://eth-mainnet.g.alchemy.com/v2", "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop", "0xdac17f958d2ee523a2206206994597c13d831ec7", "0x11cc6083a3f2608e7c6b2862185e4874ea7b2b56"))
}
