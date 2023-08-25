package service

import (
	"context"
	"log"
	"testing"
	"time"

	storeConfig "github.com/0xcregis/easynode/store/config"
	"github.com/sunjiangjun/xlog"
)

func Init() *StoreHandler {
	cfg := storeConfig.LoadConfig("./../../cmd/store/config.json")
	log.Printf("%+v\n", cfg)
	xLog := xlog.NewXLogger().BuildOutType(xlog.STD).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/store/store", 24*time.Hour)
	return NewStoreHandler(&cfg, xLog)
}

func TestService_ReadTxFromKafka(t *testing.T) {
	s := Init()
	s.ReadTxFromKafka(301, s.config.Chains[0].KafkaCfg, context.Background())
}

func TestStoreHandler_ReadBlockFromKafka(t *testing.T) {
	s := Init()
	s.ReadBlockFromKafka(301, s.config.Chains[0].KafkaCfg, context.Background())
}

func TestStoreHandler_ReadSubTxFromKafka(t *testing.T) {
	s := Init()
	s.ReadSubTxFromKafka(301, s.config.Chains[0].KafkaCfg, context.Background())
}

func TestStoreHandler_ReadReceiptFromKafka(t *testing.T) {
	s := Init()
	s.ReadReceiptFromKafka(301, s.config.Chains[0].KafkaCfg, context.Background())
}

func TestStoreHandler_ReadBackupTxFromKafka(t *testing.T) {
	s := Init()
	s.ReadBackupTxFromKafka(200, s.config.Chains[0].KafkaCfg, context.Background())
}
