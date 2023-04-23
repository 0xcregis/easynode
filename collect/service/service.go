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
	SendNodeTask(list []*NodeTask) error
	UpdateNodeTaskStatus(key string, status int) error
	UpdateNodeTaskStatusWithBatch(keys []string, status int) error
	GetNodeTask(key string) (*NodeTask, error)
	ResetNodeTask(oldKey, key string) error
	StoreExecTask(key string, task *NodeTask)
}

// BlockChainInterface 公链接口
type BlockChainInterface interface {
	GetTx(txHash string, task *config.TxTask, log *logrus.Entry) *TxInterface
	GetReceipt(txHash string, task *config.ReceiptTask, log *logrus.Entry) *ReceiptInterface
	GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, log *logrus.Entry) []*ReceiptInterface
	GetBlockByNumber(blockNumber string, task *config.BlockTask, log *logrus.Entry) (*BlockInterface, []*TxInterface)
	GetBlockByHash(blockHash string, cfg *config.BlockTask, log *logrus.Entry) (*BlockInterface, []*TxInterface)
	BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error)
	Monitor()
}
