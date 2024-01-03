package bnb

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/blockchain_config.json")
	return NewBnb(cfg.Cluster[2610], 2610, xlog.NewXLogger())
}

func TestBnb_Token(t *testing.T) {
	s := Init()
	resp, err := s.Token(202, "0x55d398326f99059fF775485246999027B3197955", "", "20")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_TokenBalance(t *testing.T) {
	s := Init()
	resp, err := s.TokenBalance(2610, "0xf0e4939183a76746e602a12c389ab183be4290b1", "0x337610d27c682e347c9cd60bd4b3b107c9d34ddd", "")

	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(202)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetBlockByNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByNumber(202, "0x1effbd0", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetBlockByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByHash(202, "0x940520e526c180ff862615e8f9f04212192c0dfc79b5c8ba1340a6d7e83bde06", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(202, "0x97F210e4a942f50AFC633E64b2F6088ef7065fB0", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(202, "0x4f4ee0761b29670f2261d5e7632f1215c9807f0e7fbec6a151e1201b97546383")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockReceiptByBlockNumber(202, fmt.Sprintf("0x%x", 32504450))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestBnb_GetTransactionReceiptByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTransactionReceiptByHash(202, "0x4f4ee0761b29670f2261d5e7632f1215c9807f0e7fbec6a151e1201b97546383")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
