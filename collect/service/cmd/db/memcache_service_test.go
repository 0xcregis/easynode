package db

import (
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"testing"
	"time"
)

func Init() service.StoreTaskInterface {
	cfg := config.LoadConfig("./../../../../cmd/collect/config.json")
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/collect/task", 24*time.Hour)
	receiverCh := make(chan []*kafka.Message, 0)
	return NewTaskCacheService(cfg.Chains[0], receiverCh, x)
}

func TestService_GetTaskByTx(t *testing.T) {
}

func TestService_UpdateTaskStatus(t *testing.T) {
	s := Init()
	s.GetNodeTask("200_blockId_0xfaef00a7e1c58c9d6f77512a8707010ae94c2ca4078489f2e22aae0d63e71aa3")
}
