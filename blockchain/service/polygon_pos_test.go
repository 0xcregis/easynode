package service

import (
	"fmt"
	"testing"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/polygonpos"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func Init4() blockchain.API {
	cfg := config.LoadConfig("./../../cmd/blockchain/config_polygon.json")
	return NewPolygonPos(cfg.Cluster[201], polygonpos.NewChainClient(), xlog.NewXLogger())
}

func TestPolygonPos_GetLatestBlock(t *testing.T) {
	s := Init4()
	resp, err := s.LatestBlock(201)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockByNumber(t *testing.T) {
	s := Init4()
	resp, err := s.GetBlockByNumber(201, fmt.Sprintf("0x%x", 45611820), false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockByHash(t *testing.T) {
	s := Init4()
	resp, err := s.GetBlockByHash(201, "0xfe88f073cc89fa63752de1a0fa9cc0e78bc89c295736d7d2e8ade0ad87936b00", false)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_Balance(t *testing.T) {
	s := Init4()
	resp, err := s.Balance(201, "0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270", "latest")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetTxByHash(t *testing.T) {
	s := Init4()
	resp, err := s.GetTxByHash(201, "0x9f656ad21cad7853f58aa05191ec4c11bd0459f40bec1a259f089fce4c80232f")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetBlockReceiptByBlockNumber(t *testing.T) {
	s := Init4()
	resp, err := s.GetBlockReceiptByBlockNumber(201, fmt.Sprintf("0x%x", 45611899))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}

func TestPolygonPos_GetTransactionReceiptByHash(t *testing.T) {
	s := Init4()
	resp, err := s.GetTransactionReceiptByHash(201, "0x9f656ad21cad7853f58aa05191ec4c11bd0459f40bec1a259f089fce4c80232f")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(resp)
	}
}
