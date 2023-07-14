package db

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/store/config"
	"github.com/sunjiangjun/xlog"
)

func Init2() *CacheService {
	cfg := config.LoadConfig("./../../../cmd/store/config.json")
	return NewCacheService(cfg.Chains, xlog.NewXLogger())
}

func TestCacheService_GetMonitorAddress(t *testing.T) {
	s := Init2()
	log.Println(s.GetMonitorAddress(200))
}
