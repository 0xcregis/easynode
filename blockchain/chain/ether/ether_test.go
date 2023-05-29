package ether

import (
	"log"
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalance("https://eth-mainnet.g.alchemy.com/v2", "RzxBjjh_c4y0LVHZ7GNm8zoXEZR3HYop", "0x3edbf7e027d280bcd8126a87f382941409364269", "0x3edbf7e027d280bcd8126a87f382941409364269"))
	log.Println(c.GetTokenBalance("https://nd-422-757-666.p2pify.com", "0a9d79d93fb2f4a4b1e04695da2b77a7", "0xdac17f958d2ee523a2206206994597c13d831ec7", "0x163118a92fb85400f36Ad418252F6D75403b86ad"))
	log.Println(c.GetTokenBalance("https://nd-422-757-666.p2pify.com", "0a9d79d93fb2f4a4b1e04695da2b77a7", "0x05fe069626543842439ef90d9fa1633640c50cf1", "0xef743eab534bdeba2de7ab30c1117e8b8206d85e"))
}
