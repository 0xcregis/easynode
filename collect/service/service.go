package service

import (
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/uduncloud/easynode/collect/config"
)

type Common interface {
	Start()
	Stop()
}

type StoreTaskInterface interface {
	SendNodeTask(list []*NodeTask) []*kafka.Message
	UpdateNodeTaskStatus(key string, status int) error
	UpdateNodeTaskStatusWithBatch(keys []string, status int) error
	GetNodeTask(key string) (*NodeTask, error)
	ResetNodeTask(oldKey, key string) error
	StoreNodeTask(key string, task *NodeTask)

	StoreContract(blockchain int64, contract string, data string) error
	GetContract(blockchain int64, contract string) (string, error)
	GetAllKeyForContract(blockchain int64, key string) ([]string, error)

	StoreErrTxNodeTask(blockchain int64, key string, data any) error
	GetErrTxNodeTask(blockchain int64, key string) (string, error)
	GetAllKeyForErrTx(blockchain int64, key string) ([]string, error)
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
