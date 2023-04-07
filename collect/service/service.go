package service

import (
	"github.com/sirupsen/logrus"
	"github.com/uduncloud/easynode/collect/config"
)

type Common interface {
	Start()
	Stop()
}

type StoreTaskInterface interface {
	GetTaskWithTx(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithReceipt(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithBlock(blockChain int, nodeId string) ([]*NodeTask, error)
	AddNodeTask(list []*NodeTask) error
	UpdateNodeTaskStatus(key string, status int) error
	UpdateNodeTaskStatusWithBatch(keys []string, status int) error

	StoreExecTask(key string, task *NodeTask)
}

// BlockChainInterface 公链接口
type BlockChainInterface interface {
	GetTx(txHash string, task *config.TxTask, log *logrus.Entry) *Tx
	GetReceipt(txHash string, task *config.ReceiptTask, log *logrus.Entry) *Receipt
	GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, log *logrus.Entry) []*Receipt
	GetBlockByNumber(blockNumber string, task *config.BlockTask, log *logrus.Entry) (*Block, []*Tx)
	GetBlockByHash(blockHash string, cfg *config.BlockTask, log *logrus.Entry) (*Block, []*Tx)
	BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error)
	Monitor()
}
