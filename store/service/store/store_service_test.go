package store

import (
	"context"
	"github.com/sunjiangjun/xlog"
	storeConfig "github.com/uduncloud/easynode/store/config"
	"log"
	"testing"
	"time"
)

func Init() *Service {
	cfg := storeConfig.LoadConfig("./../../../cmd/store/config.json")
	log.Printf("%+v\n", cfg)
	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/store/store", 24*time.Hour)
	return NewStoreService(&cfg, xLog)
}
func TestService_ReadTxFromKafka(t *testing.T) {
	s := Init()
	s.ReadTxFromKafka(200, s.config.Chains[0].KafkaCfg, context.Background())
}
