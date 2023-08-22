package blockchain

import "github.com/0xcregis/easynode/blockchain/config"

type ChainConn interface {
	SendRequestToChain(host string, token string, query string) (string, error)
	SendRequestToChainByHttp(host string, token string, query string) (string, error)
	Subscribe(host string, token string) (string, error)
	UnSubscribe(host string, token string) (string, error)
	GetTokenBalance(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetTokenBalanceByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
}

type ChainNet interface {
	SendReq(blockChain int64, reqBody string) (string, error)
	SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error)
}

type ChainCluster interface {
	BalanceCluster() *config.NodeCluster
}
