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

// KafkaInterface 接口
type KafkaInterface interface {
	//ProductMessage 生产消息到指定kafka.topic
	ProductMessage(msg *kafka.Message)
}

// NodeInterface 节点信息接口
type NodeInterface interface {
	//PushNodeInfo  定时的上传node info (20s)
	PushNodeInfo(node *NodeInfo)
}

type TaskInterface interface {
	GetTaskWithTx(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithReceipt(blockChain int, nodeId string) ([]*NodeTask, error)
	GetTaskWithBlock(blockChain int, nodeId string) ([]*NodeTask, error)
	UpdateTaskStatus(id int64, status int) error
	AddTaskSource(source *NodeSource) error
	AddNodeTask(list []*NodeTask) error
	UpdateNodeTaskStatusWithTx(task *NodeTask) error
	UpdateNodeTaskListStatusWithTx(txHashList []string, taskType int, blockChain int, taskStatus int) error
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
