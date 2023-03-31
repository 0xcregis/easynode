package db

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"log"
	"testing"
	"time"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../../config_tron.json")
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/collect/task", 24*time.Hour)
	return NewMySQLTaskService(cfg.TaskDb, x)
}

func TestService_GetTaskByBlock(t *testing.T) {

	s := Init()
	log.Println(s.GetTaskWithBlock(200, "0ea23437-ccd9-4abb-a528-abd1e28182f0"))
}

func TestService_GetTaskByTx(t *testing.T) {
}

func TestService_UpdateTaskStatus(t *testing.T) {
	s := Init()
	s.UpdateNodeTaskStatus(1, 2)
}

func TestService_AddTaskSource(t *testing.T) {
	s := Init()
	s.AddNodeSource(&service.NodeSource{
		SourceType: 1,
		BlockChain: 200,
		TxHash:     "0x8fc90a6c3ee3001cdcbbb685b4fbe67b1fa2bec575b15b0395fea5540d0901ae",
	})
}

func TestService_AddTaskSourceList(t *testing.T) {
	s := Init()
	err := s.AddNodeTaskList([]*service.NodeSource{&service.NodeSource{
		SourceType: 10,
		BlockChain: 200,
		TxHash:     "0x8fc90a6c3ee3001cdcbbb685b4fbe67b1fa2bec575b15b0395fea5540d0901ae",
	}})

	if err != nil {
		panic(err)
	}
}

func TestService_GetNodeTaskByBlockHash(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTaskByBlockHash("0x3e14a137361b2002758bd80bf4ce87d863dd2a355e1138bd82e05c60d70f45c9", 2, 200))
}

func TestService_GetNodeTaskByBlockNumber(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTaskByBlockNumber("16073507", 2, 200))
}

func TestService_GetNodeTaskWithTxs(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTaskWithTxs([]string{"0x14c731a2721e19d8927246f6d1df4aa222b9327022bf6c858e69e8035eac3a19"}, 1, 205, 1))
}
