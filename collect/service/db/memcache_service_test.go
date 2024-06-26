package db

import (
	"log"
	"testing"
	"time"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/sunjiangjun/xlog"
)

func Init() collect.StoreTaskInterface {
	cfg := config.LoadConfig("./../../../cmd/collect/collect_config_bnb_test.json")
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/collect/task", 24*time.Hour)
	return NewTaskCacheService(cfg.Chains[0], x)
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
	n := collect.NodeTask{NodeId: "1bf67775-80f3-4482-a960-c0af3a964cba", Id: 1685068437248198000, TxHash: "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", BlockChain: 200, TaskStatus: 0, TaskType: 1}
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
	n := collect.NodeTask{NodeId: "1bf67775-80f3-4482-a960-c0af3a964cba", Id: 1685068437248198000, TxHash: "0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", BlockChain: 200, TaskStatus: 0, TaskType: 1}
	n.CreateTime = time.Now()
	n.LogTime = time.Now()
	s.StoreNodeTask("200_txId_0x5e856ee7b4b43efc94a0fd18960e2a085243b2cd5b0db798d6709bd67c39ac0d", &n, true)
}

func TestService_ResetNodeTask(t *testing.T) {
	s := Init()
	s.ResetNodeTask(205, "205_blockId_00000000020b63869ef3d3034d74200e46ba92271c286277b707ab47065b134f", "205_blockId_00000000020b63869ef3d3034d74200e46ba92271c286277b707ab47065b134f")
}

func TestService_UpdateNodeTaskStatus(t *testing.T) {
	s := Init()
	s.UpdateNodeTaskStatus("2610_tx_36140539", 1)
}

func TestService_StoreLatestBlock(t *testing.T) {
	s := Init()
	str := "{\\\"id\\\":1687945592603474952}"
	s.StoreLatestBlock(200, "LatestBlock", str, "15605643")
}
