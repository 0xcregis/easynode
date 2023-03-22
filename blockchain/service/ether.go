package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/chain"
	"github.com/uduncloud/easynode/blockchain/chain/ether"
	"github.com/uduncloud/easynode/blockchain/config"
	"math/rand"
)

type Ether struct {
	log              *xlog.XLog
	nodeCluster      map[int64][]*config.NodeCluster
	blockChainClient chain.BlockChain
}

func (e *Ether) GetBlockByHash(chainCode int64, hash string) (string, error) {
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_getBlockByHash",
		  "params": [
			"%v",
 			false
		  ]
		}`

	req = fmt.Sprintf(req, hash)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) GetBlockByNumber(chainCode int64, number string) (string, error) {
	req := `
			{
			  "id": 1,
			  "jsonrpc": "2.0",
			  "method": "eth_getBlockByNumber",
			  "params": [
				"%v",
				false
			  ]
			}
			`
	req = fmt.Sprintf(req, number)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) GetTxByHash(chainCode int64, hash string) (string, error) {
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "params": [
			"%v"
		  ],
		  "method": "eth_getTransactionByHash"
		}
		`
	req = fmt.Sprintf(req, hash)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendEthReq(chainCode, req)
}

func NewEth(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) API {

	var blockChainClient chain.BlockChain
	for k, _ := range cluster {
		if k == 200 {
			blockChainClient = ether.NewChainClient()
		}
	}
	return &Ether{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
}
func (e *Ether) Balance(chainCode int64, address string, tag string) (string, error) {
	if len(tag) < 1 {
		tag = "latest"
	}
	req := `{
				 "id": 1,
				 "jsonrpc": "2.0",
				 "params": [
					  "%v",
					  "%v"
				 ],
				 "method": "eth_getBalance"
			}`

	req = fmt.Sprintf(req, address, tag)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	cluster := e.BalanceCluster(chainCode)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	mp, err := e.blockChainClient.GetTokenBalance(cluster.NodeUrl, cluster.NodeToken, contractAddr, address)
	if err != nil {
		return "", err
	}
	rs, _ := json.Marshal(mp)
	return string(rs), nil
}

func (e *Ether) Nonce(chainCode int64, address string, tag string) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "params": [
					  "%v",
					  "%v"
				 ],
				 "method": "eth_getTransactionCount"
			}
			`
	req = fmt.Sprintf(req, address, tag)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) LatestBlock(chainCode int64) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "method": "eth_blockNumber"
			}
			`
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
					 "id": 1,
					 "jsonrpc": "2.0",
					 "params": [
						  "%v"
					 ],
					 "method": "eth_sendRawTransaction"
				}`
	req = fmt.Sprintf(req, signedTx)
	return e.SendEthReq(chainCode, req)
}

func (e *Ether) SendEthReq(blockChain int64, reqBody string) (string, error) {
	cluster := e.BalanceCluster(blockChain)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	if blockChain == 200 {
		return e.blockChainClient.EthSendRequestToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	}

	return "", errors.New("blockChainCode is error")
}

func (e *Ether) BalanceCluster(blockChain int64) *config.NodeCluster {
	cluster, ok := e.nodeCluster[blockChain]
	if !ok {
		//不存在节点
		return nil
	}

	//todo 后期重新设计和优化
	var resultCluster *config.NodeCluster
	//temps := make([]*config.NodeCluster, 0, 10)
	//for _, v := range cluster {
	//	if v.ErrorCount < 50 {
	//		temps = append(temps, v)
	//	}
	//}

	l := len(cluster)
	if l > 1 {
		m := rand.Intn(l)
		resultCluster = cluster[m]
	} else if l == 1 {
		resultCluster = cluster[0]
	} else {
		return nil
	}
	return resultCluster

}
