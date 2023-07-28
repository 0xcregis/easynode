package collect

import (
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
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

	StoreClusterNode(blockChain int64, prefix string, data any) error
	GetClusterNode(blockChain int64, prefix string) (map[string]int64, error)
	StoreClusterHealthStatus(blockChain int64, data map[string]int64) error
}

// BlockChainInterface 公链接口
type BlockChainInterface interface {
	GetTx(txHash string, log *logrus.Entry) *TxInterface
	GetReceipt(txHash string, log *logrus.Entry) (*ReceiptInterface, error)
	GetReceiptByBlock(blockHash, number string, log *logrus.Entry) ([]*ReceiptInterface, error)
	GetBlockByNumber(blockNumber string, log *logrus.Entry, flag bool) (*BlockInterface, []*TxInterface)
	GetBlockByHash(blockHash string, log *logrus.Entry, flag bool) (*BlockInterface, []*TxInterface)
	Monitor()
	CheckAddress(tx []byte, addrList map[string]int64) bool
}
