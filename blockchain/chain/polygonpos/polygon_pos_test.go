package ether

import (
	"log"
	"testing"
)

func TestPolygon_GetToken(t *testing.T) {
	c := NewChainClient()
	log.Println(c.GetTokenBalance("https://nd-422-757-666.p2pify.com", "0a9d79d93fb2f4a4b1e04695da2b77a7", "0x05fe069626543842439ef90d9fa1633640c50cf1", "0xef743eab534bdeba2de7ab30c1117e8b8206d85e"))
}
