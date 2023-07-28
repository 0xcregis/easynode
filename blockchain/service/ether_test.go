package service

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/ether"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init3() blockchain.API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_ether.json")
	return NewEth(cfg.Cluster[200], ether.NewChainClient(), xlog.NewXLogger())
}

func TestEther_GetLatestBlock(t *testing.T) {
	s := Init3()
	resp, err := s.LatestBlock(200)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByNumber(t *testing.T) {
	s := Init3()
	resp, err := s.GetBlockByNumber(200, "0xF3F088", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByHash(t *testing.T) {
	s := Init3()
	resp, err := s.GetBlockByHash(200, "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_Balance(t *testing.T) {
	s := Init3()
	resp, err := s.Balance(200, "0x06d9ca334a8a74474e9b6ee31280c494321ae759", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetTxByHash(t *testing.T) {
	s := Init3()
	resp, err := s.GetTxByHash(200, "0x840d9a67084505cd06221e1b7f4690356e7e789ea9827cf196e8f9a875d8e42d")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init3()
	resp, err := s.GetBlockReceiptByBlockNumber(200, fmt.Sprintf("0x%x", 17790088))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetTransactionReceiptByHash(t *testing.T) {
	s := Init3()
	resp, err := s.GetTransactionReceiptByHash(200, "0x840d9a67084505cd06221e1b7f4690356e7e789ea9827cf196e8f9a875d8e42d")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
