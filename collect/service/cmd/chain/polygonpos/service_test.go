package polygonpos

import (
	"log"
	"testing"

	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (service.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("17658423", x.WithFields(logrus.Fields{}), true)
	log.Printf("%+v", b)
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("0xb7959ab3ed0f8b424f897ba8ec28168358079be0030a9337d20252f7d4c18cfd", x.WithFields(logrus.Fields{}))
	log.Println(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("0x2eadb796722102575e8a1fe75241cba96f605659c160a65b3ad9473027b74c21", x.WithFields(logrus.Fields{}))
	if err != nil {
		panic(err)
	}
	log.Printf("%+v", r)
}

func TestService_BalanceCluster(t *testing.T) {

}

func TestService_Monitor(t *testing.T) {
	_, cfg, _ := Init()

	c1 := cfg.Chains[0]
	c2 := c1.CopyChain()

	c1.BlockTask.FromCluster[0].ErrorCount = 11
	c2.BlockTask.FromCluster[0].ErrorCount = 12

	c1.TxTask.FromCluster[0].ErrorCount = 13
	c2.TxTask.FromCluster[0].ErrorCount = 14

	c1.ReceiptTask.FromCluster[0].ErrorCount = 15
	c2.ReceiptTask.FromCluster[0].ErrorCount = 16

	c1.BlockChainCode = 101
	c2.BlockChainCode = 102

	log.Println(c2)
	log.Println(c1)
}
