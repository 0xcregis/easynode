package chain

type BlockChain interface {
	EthSendRequestToChain(host string, token string, query string) (string, error)
	SendRequestToChain(host string, token string, query string) (string, error)
	GetTokenBalance(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
}
