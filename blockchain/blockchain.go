package blockchain

import "github.com/0xcregis/easynode/blockchain/config"

type ChainConn interface {
	SendRequestToChain(host string, token string, query string) (string, error)
	SendRequestToChainByHttp(host string, token string, query string) (string, error)
	Subscribe(host string, token string) (string, error)
	UnSubscribe(host string, token string) (string, error)
	GetToken20(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetToken20ByHttp(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetToken721(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
	GetToken1155(host string, token string, contractAddress string, userAddress string) (map[string]interface{}, error)
}

type NFT interface {
	TokenURI(host string, key string, contractAddress string, tokenId string, eip int64) (string, error)
	BalanceOf(host string, key string, contractAddress string, address string, tokenId string, eip int64) (string, error)
	OwnerOf(host string, key string, contractAddress string, tokenId string, eip int64) (string, error)
	TotalSupply(host string, key string, contractAddress string, eip int64) (string, error)
}

type ChainNet interface {
	SendReq(blockChain int64, reqBody string) (string, error)
	SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error)
}

type ChainCluster interface {
	BalanceCluster() *config.NodeCluster
}
