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
	return NewTron(cfg.Cluster[195], 195, xlog.NewXLogger())
}

func Init2() *Tron {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_tron.json")
	return NewTron2(cfg.Cluster[195], 195, xlog.NewXLogger())
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
	tx := `0A8A010A0202DB2208C89D4811359A28004098A4E0A6B52D5A730802126F0A32747970652E676F6F676C65617069732E636F6D2F70726F746F636F6C2E5472616E736665724173736574436F6E747261637412390A07313030303030311215415A523B449890854C8FC460AB602DF9F31FE4293F1A15416B0580DA195542DDABE288FEC436C7D5AF769D24206412418BF3F2E492ED443607910EA9EF0A7EF79728DAAAAC0EE2BA6CB87DA38366DF9AC4ADE54B2912C1DEB0EE6666B86A07A6C7DF68F1F9DA171EEE6A370B3CA9CBBB00`
	resp, err := c.SendRawTransaction(195, tx)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestTron_GetAccountResource(t *testing.T) {
	c := Init2()
	resp, err := c.GetAccountResource(195, "TXeZAknJe2gbqSJyZYXbNMVvQgsKQbSoxX")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
