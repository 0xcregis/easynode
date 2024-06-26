package ether

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_ether.json")
	return NewEth(cfg.Cluster[200], 200, xlog.NewXLogger())
}

func Init2() blockchain.NftApi {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_ether.json")
	return NewNftEth(cfg.Cluster[60], 60, xlog.NewXLogger())
}

func Init3() blockchain.ExApi {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_ether.json")
	return NewEth2(cfg.Cluster[200], 200, xlog.NewXLogger())
}

func TestEther_Token(t *testing.T) {
	s := Init()
	resp, err := s.Token(200, "0x4577fcfB0642afD21b5f2502753ED6D497B830E9", "", "721")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_BalanceOf(t *testing.T) {
	s := Init2()
	resp, err := s.BalanceOf(60, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "0xE093b32E23646248990d121aA02d2B493B538E41", "2095", 721)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_TokenURI(t *testing.T) {
	s := Init2()
	resp, err := s.TokenURI(200, "0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", "2095", 721)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_TokenBalance(t *testing.T) {
	s := Init()
	resp, err := s.TokenBalance(200, "0xdac17f958d2ee523a2206206994597c13d831ec7", "0xdac17f958d2ee523a2206206994597c13d831ec7", "")

	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_Nonce(t *testing.T) {
	s := Init()
	resp, err := s.Nonce(200, "0xae2Fc483527B8EF99EB5D9B44875F005ba1FaE13", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(200)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByNumber(200, "0xF3F088", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByHash(200, "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(200, "0x06d9ca334a8a74474e9b6ee31280c494321ae759", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(200, "0x840d9a67084505cd06221e1b7f4690356e7e789ea9827cf196e8f9a875d8e42d")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockReceiptByBlockNumber(200, fmt.Sprintf("0x%x", 17790088))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_GetTransactionReceiptByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTransactionReceiptByHash(200, "0x840d9a67084505cd06221e1b7f4690356e7e789ea9827cf196e8f9a875d8e42d")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestEther_TraceTransaction(t *testing.T) {
	ex := Init3()
	resp, err := ex.TraceTransaction(200, "0x999cabe1fcca80148290827a8c655734531615cb22d30faa222ec7a67928587b")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}

}
