package tron

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_tron.json")
	return NewTron(cfg.Cluster[198], 198, xlog.NewXLogger())
}

func TestTron_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(205, "TXeZAknJe2gbqSJyZYXbNMVvQgsKQbSoxX", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_TokenBalance(t *testing.T) {
	s := Init()
	resp, err := s.TokenBalance(198, "TDGLAKnr2SeYJHht6YxrtZqfRrVR9RFdwV", "TBo8ZFTG13PZZTgSVbuTVBi5FrCZjDedFU", "")

	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_Nonce(t *testing.T) {
	s := Init()
	resp, err := s.Nonce(205, "TWGZbjofbTLY3UCjCV4yiLkRg89zLqwRgi", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(205)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(205, "d0ff91487dd11ab6bd2cffa4af97bb472ede4f1713786fa2b15bf32011d0b681")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetTransactionReceiptByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTransactionReceiptByHash(205, "d0ff91487dd11ab6bd2cffa4af97bb472ede4f1713786fa2b15bf32011d0b681")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetBlockByNumber(t *testing.T) {
	c := Init()
	log.Println(c.GetBlockByNumber(205, "45611899", false))
}

func TestTron_GetBlockByHash(t *testing.T) {
	c := Init()
	log.Println(c.GetBlockByHash(205, "0000000002f2f5f62d94d85ec1abf2c0dfc26d72da4f5e5d5a2624d51e231425", false))
}

func TestTron_GetBlockReceiptByBlockNumber(t *testing.T) {
	c := Init()
	log.Println(c.GetBlockReceiptByBlockNumber(205, "45611899"))
}

func TestTron_SendRawTransaction(t *testing.T) {
	c := Init()
	tx := `
{
  "raw_data": {
    "contract": [
      {
        "parameter": {
          "value": {
            "amount": 1000,
            "owner_address": "41608f8da72479edc7dd921e4c30bb7e7cddbe722e",
            "to_address": "41e9d79cc47518930bc322d9bf7cddd260a0260a8d"
          },
          "type_url": "type.googleapis.com/protocol.TransferContract"
        },
        "type": "TransferContract"
      }
    ],
    "ref_block_bytes": "5e4b",
    "ref_block_hash": "47c9dc89341b300d",
    "expiration": 1591089627000,
    "timestamp": 1591089567635
  },
  "raw_data_hex": "0a025e4b220847c9dc89341b300d40f8fed3a2a72e5a66080112620a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412310a1541608f8da72479edc7dd921e4c30bb7e7cddbe722e121541e9d79cc47518930bc322d9bf7cddd260a0260a8d18e8077093afd0a2a72e"
}
  `

	resp, err := c.SendRawTransaction(198, tx)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
