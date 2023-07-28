package blockchain

type BlockChain interface {
	SendRequestToChain(host string, token string, query string) (string, error)
	SendRequestToChainByHttp(host string, token string, query string) (string, error)
	EthSubscribe(host string, token string) (string, error)
	EthUnSubscribe(host string, token string) (string, error)
	GetTokenBalance(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
}
