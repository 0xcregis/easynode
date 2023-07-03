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
	SendNodeTask(list []*NodeTask, partitions []int64) []*kafka.Message

	UpdateNodeTaskStatus(key string, status int) error
	UpdateNodeTaskStatusWithBatch(keys []string, status int) error
	GetNodeTask(blockchain int64, key string) (int64, *NodeTask, error)
	DelNodeTask(blockchain int64, key string) (int64, *NodeTask, error)
	ResetNodeTask(blockchain int64, oldKey, key string) error
	StoreNodeTask(key string, task *NodeTask)
	GetAllKeyForNodeTask(blockchain int64) ([]string, error)

	StoreContract(blockchain int64, contract string, data string) error
	GetContract(blockchain int64, contract string) (string, error)
	GetAllKeyForContract(blockchain int64) ([]string, error)

	StoreErrTxNodeTask(blockchain int64, key string, data any) error
	GetErrTxNodeTask(blockchain int64, key string) (int64, string, error)
	DelErrTxNodeTask(blockchain int64, key string) (string, error)
	GetAllKeyForErrTx(blockchain int64) ([]string, error)

	GetMonitorAddress(blockChain int64) ([]string, error)

	StoreLatestBlock(blockchain int64, key string, data any, number string) error

	StoreNodeId(blockchain int64, key string, data any) error
	GetAllNodeId(blockchain int64) ([]string, error)
}

// BlockChainInterface 公链接口
type BlockChainInterface interface {
	GetTx(txHash string, task *config.TxTask, log *logrus.Entry) *TxInterface
	GetReceipt(txHash string, task *config.ReceiptTask, log *logrus.Entry) (*ReceiptInterface, error)
	GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, log *logrus.Entry) ([]*ReceiptInterface, error)
	GetBlockByNumber(blockNumber string, task *config.BlockTask, log *logrus.Entry, flag bool) (*BlockInterface, []*TxInterface)
	GetBlockByHash(blockHash string, cfg *config.BlockTask, log *logrus.Entry, flag bool) (*BlockInterface, []*TxInterface)
	BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error)
	Monitor()
	CheckAddress(tx []byte, addrList []string) bool
}
