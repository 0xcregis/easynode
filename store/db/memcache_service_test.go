package db

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/config"
	"github.com/sunjiangjun/xlog"
)

func Init2() *CacheService {
	cfg := config.LoadConfig("./../../cmd/store/config_tron.json")
	return NewCacheService(cfg.Chains, xlog.NewXLogger())
}

func TestCacheService_GetMonitorAddress(t *testing.T) {
	s := Init2()
	log.Println(s.GetMonitorAddress(200))
}

func TestCacheService_SetMonitorAddress(t *testing.T) {
	s := Init2()
	err := s.SetMonitorAddress(205, []*store.MonitorAddress{{Token: "063470c2-2152-44e5-973c-3fcb21cd4579", Address: "0x41a614f803b6fd780986a42c78ec9c7f77e6ded13c", BlockChain: 205}})
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}
}
