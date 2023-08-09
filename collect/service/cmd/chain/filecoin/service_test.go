package filecoin

import (
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_filecoin.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, _ := s.GetBlockByNumber("3106211", x.WithFields(logrus.Fields{}), false)
	t.Logf("%+v", b)
}

func TestService_GetMultiBlockByNumber(t *testing.T) {
	s, _, x := Init()
	list, _ := s.GetMultiBlockByNumber("3106211", x.WithFields(logrus.Fields{}), false)
	t.Logf("block.list:%v", len(list))
}

func TestService_GetBlockByHash(t *testing.T) {
	s, _, x := Init()
	block, txList := s.GetBlockByHash("bafy2bzacecs5veov5flezd6eol7ezbnwrjr36jkym5q4i7yfbne5nnubctrps", x.WithFields(logrus.Fields{}), true)
	t.Logf("block.hash:%v, tx.list:%v", block.BlockHash, len(txList))
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("bafy2bzacea22wldxdpp5nqvckkusyogyqggny3ycauzf2n7p3r7cuist4x7jm", x.WithFields(logrus.Fields{}))
	t.Log(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("bafy2bzacea22wldxdpp5nqvckkusyogyqggny3ycauzf2n7p3r7cuist4x7jm", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}

func TestService_CheckAddress(t *testing.T) {
	s, _, _ := Init()
	str := `{"hash":"bafy2bzacebcb3rqh24ghzr2gc4vhs7qqsncl5xad2eck2dwvnb6ub3txajzks","tx":"{\"Version\":0,\"To\":\"f02146033\",\"From\":\"f3qt5uzschha7cc5e7akcgrrbgs44uthlgc6orz3b3irqwzzjumpanb3slshgjmxaz7gjh2sqsm4fbubrmjngq\",\"Nonce\":125340,\"Value\":\"52555744144206138\",\"GasLimit\":53854543,\"GasFeeCap\":\"464213390\",\"GasPremium\":\"1421564\",\"Method\":6,\"Params\":\"iggZ9WDYKlgpAAGC4gOB6AIg5ibenwmI2zZYV/fo3p65x2SUnuaMaPQjz946+gpdtUQaAC8xZIAaAEbbbfQAAAA=\",\"CID\":{\"/\":\"bafy2bzacebcb3rqh24ghzr2gc4vhs7qqsncl5xad2eck2dwvnb6ub3txajzks\"}}"}`
	s.CheckAddress([]byte(str), nil)
}

func TestService_GetReceiptByBlock(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceiptByBlock("", "3106211", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}
