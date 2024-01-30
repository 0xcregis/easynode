package polygon

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/common/util"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

func Init() blockchain.API {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_polygon.json")
	return NewPolygonPos(cfg.Cluster[66], 66, xlog.NewXLogger())
}

func Init2() blockchain.ExApi {
	cfg := config.LoadConfig("./../../../cmd/blockchain/config_polygon.json")
	return NewPolygonPos2(cfg.Cluster[64], 64, xlog.NewXLogger())
}

func TestPolygonPos_GetLatestBlock(t *testing.T) {
	s := Init()
	resp, err := s.LatestBlock(201)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockByNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByNumber(201, fmt.Sprintf("0x%x", 45611820), false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockByHash(201, "0xfe88f073cc89fa63752de1a0fa9cc0e78bc89c295736d7d2e8ade0ad87936b00", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_Balance(t *testing.T) {
	s := Init()
	resp, err := s.Balance(201, "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetTxByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTxByHash(201, "0x9f656ad21cad7853f58aa05191ec4c11bd0459f40bec1a259f089fce4c80232f")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init()
	resp, err := s.GetBlockReceiptByBlockNumber(201, fmt.Sprintf("0x%x", 45611899))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetTransactionReceiptByHash(t *testing.T) {
	s := Init()
	resp, err := s.GetTransactionReceiptByHash(66, "0x9a022bff505dec115478d2f368092918c4bbbf82c63f541e576ab0407485885a")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GasPrice(t *testing.T) {
	s := Init2()
	resp, err := s.GasPrice(64)
	if err != nil {
		t.Error(err)
	} else {
		gas := gjson.Parse(resp).Get("result").String()
		gas, _ = util.HexToInt(gas)
		t.Log(gas)
	}
}
