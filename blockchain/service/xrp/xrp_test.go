package xrp

import (
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_xrp.json")
	return NewXRP(cfg.Cluster[310], 310, xlog.NewXLogger())
}

func TestXRP_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(310)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetBlockByNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByNumber(310, "81984330", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetBlockByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByHash(310, "B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(310, "rMwjYedjc7qqtKYVLiAccJSmCwih4LnE2q", "current")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(310, "89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockReceiptByBlockNumber(310, "81984330")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetBlockReceiptByBlockHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockReceiptByBlockHash(310, "B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestXRP_GetTransactionReceiptByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTransactionReceiptByHash(310, "89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
