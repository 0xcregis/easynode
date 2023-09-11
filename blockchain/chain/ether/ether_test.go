package ether

import (
	"testing"
)

func TestEth_GetToken(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken20("https://ethereum-mainnet-rpc.allthatnode.com", "", "0xdac17f958d2ee523a2206206994597c13d831ec7", "0xdac17f958d2ee523a2206206994597c13d831ec7")
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
