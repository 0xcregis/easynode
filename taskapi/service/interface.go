package service

type DbApiInterface interface {
	AddNodeTask(task *NodeTask) error
	QueryTxFromCh(blockChain int64, txHash string) (*Tx, error)
}
