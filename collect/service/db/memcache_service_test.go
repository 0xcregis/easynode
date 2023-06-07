package db

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"log"
	"testing"
	"time"
)

func Init() service.StoreTaskInterface {
	cfg := config.LoadConfig("./../../../cmd/collect/config.json")
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/collect/task", 24*time.Hour)
	return NewTaskCacheService(cfg.Chains[0], x)
}

func TestService_GetTaskByTx(t *testing.T) {
}

func TestService_UpdateTaskStatus(t *testing.T) {
	s := Init()
	s.GetNodeTask(200, "200_blockId_0xfaef00a7e1c58c9d6f77512a8707010ae94c2ca4078489f2e22aae0d63e71aa3")
}

func TestService_StoreContract(t *testing.T) {
	s := Init()
	log.Println(s.StoreContract(200, "0x123", "1234"))
}

func TestService_StoreErrTxNodeTask(t *testing.T) {
	s := Init()
	n := service.NodeTask{NodeId: "1bf67775-80f3-4482-a960-c0af3a964cba", Id: 1685068437248198000, TxHash: "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", BlockChain: 200, TaskStatus: 0, TaskType: 1}
	n.CreateTime = time.Now()
	n.LogTime = time.Now()
	log.Println(s.StoreErrTxNodeTask(200, "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", n))
}

func TestService_GetErrTxNodeTask(t *testing.T) {
	s := Init()
	log.Println(s.GetErrTxNodeTask(200, "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d"))
}

func TestService_StoreNodeTask(t *testing.T) {
	s := Init()
	n := service.NodeTask{NodeId: "1bf67775-80f3-4482-a960-c0af3a964cba", Id: 1685068437248198000, TxHash: "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", BlockChain: 200, TaskStatus: 0, TaskType: 1}
	n.CreateTime = time.Now()
	n.LogTime = time.Now()
	s.StoreNodeTask("200_txId_0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", &n)
}
