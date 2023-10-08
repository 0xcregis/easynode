package btc

import (
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_btc.json")
	return NewBtc(cfg.Cluster[300], 300, xlog.NewXLogger())
}

func TestBtc_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(300, "1PL6qjNjEMRhTLAnHEFJWwnvjjKGWAwFws", "")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(300)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByNumber(300, "808682", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByHash(300, "00000000000000000000fbd80dd4a8502f6a06d10c6179c602d5a0ea24f43d38", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(300, "10b54fd708ab2e5703979b4ba27ca0339882abc2062e77fbe51e625203a49642")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
