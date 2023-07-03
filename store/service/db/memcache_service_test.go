package db

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/store/config"
	"log"
	"testing"
)

func Init2() *CacheService {
	cfg := config.LoadConfig("./../../../cmd/store/config.json")
	return NewCacheService(cfg.Chains, xlog.NewXLogger())
}

func TestCacheService_GetMonitorAddress(t *testing.T) {
	s := Init2()
	log.Println(s.GetMonitorAddress(200))
}
