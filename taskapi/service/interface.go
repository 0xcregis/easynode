package service

type DbApiInterface interface {
	AddNodeTask(task *NodeTask) error
	GetActiveNodesFromDB() ([]string, error)
	QueryTxFromCh(blockChain int64, txHash string) (*Tx, error)
}
