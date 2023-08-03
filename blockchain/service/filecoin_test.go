package service

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init5() blockchain.API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_filecoin.json")
	return NewFileCoin(cfg.Cluster[301], 301, xlog.NewXLogger())
}

func TestFileCoin_GetLatestBlock(t *testing.T) {
	s := Init5()
	resp, err := s.LatestBlock(301)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_GetBlockByNumber(t *testing.T) {
	s := Init5()
	resp, err := s.GetBlockByNumber(301, "3091306", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_GetBlockByHash(t *testing.T) {
	s := Init5()
	resp, err := s.GetBlockByHash(301, "bafy2bzacebll5xugzbbrfvyzrzvn25rvdschatyuw2bharrgx3tt6qgb7cag2", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_Balance(t *testing.T) {
	s := Init5()
	resp, err := s.Balance(301, "f047684", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_GetTxByHash(t *testing.T) {
	s := Init5()
	resp, err := s.GetTxByHash(301, "bafy2bzacebttywwub67mp65nmovcb5axtas7luvfqhl7gay3pwlplfenjttwc")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init5()
	resp, err := s.GetBlockReceiptByBlockNumber(301, fmt.Sprintf("0x%x", 17790088))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_GetTransactionReceiptByHash(t *testing.T) {
	s := Init5()
	//id := cid.MustParse("bafy2bzacecctyxrgsua4w3xi64awesrikkk5dmtprda6ffipyy2fkjwobuqmy")
	resp, err := s.GetTransactionReceiptByHash(301, "bafy2bzacecctyxrgsua4w3xi64awesrikkk5dmtprda6ffipyy2fkjwobuqmy")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestFileCoin_SendRawTransaction(t *testing.T) {
	s := Init5()
	tx := `
  {
    "Message": {
      "Version": 42,
      "To": "f01234",
      "From": "f01234",
      "Nonce": 42,
      "Value": "0",
      "GasLimit": 9,
      "GasFeeCap": "0",
      "GasPremium": "0",
      "Method": 1,
      "Params": "Ynl0ZSBhcnJheQ==",
      "CID": {
        "/": "bafy2bzacebbpdegvr3i4cosewthysg5xkxpqfn2wfcz6mv2hmoktwbdxkax4s"
      }
    },
    "Signature": {
      "Type": 2,
      "Data": "Ynl0ZSBhcnJheQ=="
    },
    "CID": {
      "/": "bafy2bzacebbpdegvr3i4cosewthysg5xkxpqfn2wfcz6mv2hmoktwbdxkax4s"
    }
  }
`
	t.Log(s.SendRawTransaction(301, tx))
}

func TestFileCoin_Nonce(t *testing.T) {
	s := Init5()
	t.Log(s.Nonce(301, "f1ys5qqiciehcml3sp764ymbbytfn3qoar5fo3iwy", ""))
}
