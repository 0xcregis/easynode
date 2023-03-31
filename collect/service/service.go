package service

import (
	"github.com/sirupsen/logrus"
	"github.com/uduncloud/easynode/collect/config"
)

type Common interface {
	Start()
	Stop()
}

type MonitorDbInterface interface {
	CheckTable()
}

type TaskDbInterface interface {
	GetTaskWithTx(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithReceipt(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithBlock(blockChain int, nodeId string) ([]*NodeTask, error)
	UpdateTaskStatus(id int64, status int) error
	AddNodeTask(list []*NodeTask) error
	GetNodeTaskTable() string
	GetNodeTaskWithTxs(txHash []string, taskType int, blockChain int, taskStatus int) ([]int64, error)
	GetNodeTaskByBlockNumber(number string, taskType int, blockChain int) (*NodeTask, error)
	GetNodeTaskByBlockHash(hash string, taskType int, blockChain int) (*NodeTask, error)
	UpdateNodeTaskStatus(id int64, status int) error
	UpdateNodeTaskStatusWithBatch(ids []int64, status int) error
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
