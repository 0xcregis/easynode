package bnb

import (
	"testing"
)

func TestBnb_TokenURI(t *testing.T) {
	c := NewNFTClient()
	resp, err := c.TokenURI("https://binance.llamarpc.com", "", "0x55d398326f99059fF775485246999027B3197955", "3742", 721)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetToken(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken20("https://binance.llamarpc.com", "", "0x55d398326f99059fF775485246999027B3197955", "0x403f9D1EA51D55d0341ce3c2fBF33E09846F2C74")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetToken721(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken721("https://binance.llamarpc.com", "", "0x55d398326f99059fF775485246999027B3197955", "0x403f9D1EA51D55d0341ce3c2fBF33E09846F2C74")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetToken1155(t *testing.T) {
	c := NewChainClient()
	resp, err := c.GetToken1155("https://binance.llamarpc.com", "", "0x55d398326f99059fF775485246999027B3197955", "0x403f9D1EA51D55d0341ce3c2fBF33E09846F2C74")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_SendRequestToChain(t *testing.T) {
	c := NewChainClient()
	query := `
	{ "jsonrpc":"2.0", "method":"eth_blockNumber","params":[],"id":1}
	`
	resp, err := c.SendRequestToChain("https://binance.llamarpc.com", "", query)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
