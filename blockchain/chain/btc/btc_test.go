package btc

import (
	"testing"
)

func TestEther_SendRequestToChain(t *testing.T) {
	c := NewChainClient()
	query := `
	{"method":"getblockcount"}
	`
	resp, err := c.SendRequestToChain("https://docs-demo.btc.quiknode.pro", "", query)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

	resp, err = c.SendRequestToChain("https://docs-demo.btc.quiknode.pro", "", query)

}
