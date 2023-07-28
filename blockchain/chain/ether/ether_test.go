package ether

import (
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetTokenBalance("https://ethereum-mainnet-rpc.allthatnode.com", "", "0xdAC17F958D2ee523a2206206994597C13D831ec7", "0xd7Aa9ba6cAAC7b0436c91396f22ca5a7F31664fC")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_SendRequestToChain(t *testing.T) {
	c := NewChainClient()
	query := `
	{ "jsonrpc":"2.0", "method":"eth_blockNumber","params":[],"id":1}
	`
	resp, err := c.SendRequestToChain("https://ethereum-mainnet-rpc.allthatnode.com", "", query)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
