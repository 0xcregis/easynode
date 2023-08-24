package xrp

import (
	"testing"
)

func TestEther_SendRequestToChain(t *testing.T) {
	c := NewChainClient()
	query := `
	{
		"method": "ledger_closed",
		"params": [
			{}
		]
	}
	`
	resp, err := c.SendRequestToChain("https://xrplcluster.com", "", query)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
