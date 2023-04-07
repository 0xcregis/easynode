package db

import (
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/service"
	"testing"
	"time"
)

func Init() service.StoreTaskInterface {
	//cfg := config.LoadConfig("./../../../config_tron.json")
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/collect/task", 24*time.Hour)
	receiverCh := make(chan []*kafka.Message, 0)
	return NewMySQLTaskService(receiverCh, x)
}

func TestService_GetTaskByTx(t *testing.T) {
}

func TestService_UpdateTaskStatus(t *testing.T) {
	s := Init()
	s.UpdateNodeTaskStatus("", 2)
}
