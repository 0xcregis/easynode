package service

type API interface {
	Balance(chainCode int64, address string, tag string) (string, error)
	TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error)
	Nonce(chainCode int64, address string, tag string) (string, error)
	LatestBlock(chainCode int64) (string, error)
	SendRawTransaction(chainCode int64, signedTx string) (string, error)

	GetBlockByHash(chainCode int64, hash string) (string, error)
	GetBlockByNumber(chainCode int64, number string) (string, error)
	GetTxByHash(chainCode int64, hash string) (string, error)
	SendJsonRpc(chainCode int64, req string) (string, error)

	GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error)
	GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error)
	GetTransactionReceiptByHash(chainCode int64, hash string) (string, error)

	SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error)
	SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error)
	UnSubscribe(chainCode int64, subId string) (string, error)
}
