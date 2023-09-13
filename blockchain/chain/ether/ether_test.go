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

func TestEther_GetToken721(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken721("https://ethereum-mainnet-rpc.allthatnode.com", "", "0x4577fcfB0642afD21b5f2502753ED6D497B830E9", "0x1294332c03933c770a0d91adc7e0f1fccc7476b9")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetToken1155(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken1155("https://ethereum-mainnet-rpc.allthatnode.com", "", "0x0521fa0bf785ae9759c7cb3cbe7512ebf20fbdaa", "0x07587c046d4d4bd97c2d64edbfab1c1fe28a10e5")
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
