package polygonpos

import (
	"testing"
)

func TestPolygonPos_GetToken(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetTokenBalance("https://polygon.rpc.blxrbdn.com", "", "0xc2132d05d31c914a87c6611c10748aeb04b58e8f", "0xa006b7ba6fb6fd1df91a6c0478bc126702299a47")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_SendRequestToChain(t *testing.T) {
	c := NewChainClient()
	query := `
	{ "jsonrpc":"2.0", "method":"eth_blockNumber","params":[],"id":1}
	`
	resp, err := c.SendRequestToChain("https://polygon.rpc.blxrbdn.com", "", query)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
