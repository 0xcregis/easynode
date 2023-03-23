package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/blockchain/chain"
	"github.com/uduncloud/easynode/blockchain/chain/tron"
	"github.com/uduncloud/easynode/blockchain/config"
	"math/rand"
	"strconv"
)

type Tron struct {
	log              *xlog.XLog
	nodeCluster      map[int64][]*config.NodeCluster
	blockChainClient chain.BlockChain
}

func (t *Tron) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyid")
}

func NewTron(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) API {
	var blockChainClient chain.BlockChain
	for k, _ := range cluster {
		if k == 205 {
			blockChainClient = tron.NewChainClient()
		}
	}

	return &Tron{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
}

func (t *Tron) GetBlockByHash(chainCode int64, hash string) (string, error) {
	req := `{"value": "%v"}`
	req = fmt.Sprintf(req, hash)
	res, err := t.SendReq(chainCode, req, "wallet/getblockbyid")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	mp := gjson.Parse(res).Map()
	delete(mp, "transactions")
	r, _ := json.Marshal(mp)
	return string(r), nil
	//} else {
	//	return res, nil
	//}
}

func (t *Tron) GetBlockByNumber(chainCode int64, number string) (string, error) {
	req := `{"num": %v}`

	n, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return "", err
	}
	req = fmt.Sprintf(req, n)
	res, err := t.SendReq(chainCode, req, "wallet/getblockbynum")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	mp := gjson.Parse(res).Map()
	delete(mp, "transactions")
	r, _ := json.Marshal(mp)
	return string(r), nil
	//} else {
	//	return res, nil
	//}

}

func (t *Tron) GetTxByHash(chainCode int64, hash string) (string, error) {
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyid")
}

func (t *Tron) SendJsonRpc(chainCode int64, req string) (string, error) {
	cluster := t.BalanceCluster(chainCode)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "/jsonrpc")
	return t.blockChainClient.EthSendRequestToChain(url, cluster.NodeToken, req)
}

func (t *Tron) Balance(chainCode int64, address string, tag string) (string, error) {
	req := `{"address":"%v",  "visible": true}`
	req = fmt.Sprintf(req, address)
	res, err := t.SendReq(chainCode, req, "wallet/getaccount")
	if err != nil {
		return "", err
	}

	r := gjson.Parse(res)
	if r.Get("Error").Exists() {
		return "", errors.New(r.Get("Error").String())
	} else {
		returnStr := `{"balance":%v}`
		balance := r.Get("balance").Int()
		returnStr = fmt.Sprintf(returnStr, balance)
		return returnStr, nil
	}
}

func (t *Tron) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	cluster := t.BalanceCluster(chainCode)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "wallet/triggerconstantcontract")
	mp, err := t.blockChainClient.GetTokenBalanceByHttp(url, cluster.NodeToken, contractAddr, address)
	if err != nil {
		return "", err
	}
	rs, _ := json.Marshal(mp)
	return string(rs), nil
}

func (t *Tron) Nonce(chainCode int64, address string, tag string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) LatestBlock(chainCode int64) (string, error) {
	res, err := t.SendReq(chainCode, "", "wallet/getnowblock")
	if err != nil {
		return "", err
	}

	r := gjson.Parse(res)

	blockId := r.Get("blockID").String()
	number := r.Get("block_header.raw_data.number").Int()

	returnStr := `{"blockId":"%v","number":%v}`
	returnStr = fmt.Sprintf(returnStr, blockId, number)
	return returnStr, nil
}

func (t *Tron) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	return t.SendReq(chainCode, signedTx, "wallet/broadcasttransaction")
}

func (t *Tron) SendReq(blockChain int64, reqBody string, url string) (string, error) {
	cluster := t.BalanceCluster(blockChain)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	if blockChain == 205 {
		url = fmt.Sprintf("%v/%v", cluster.NodeUrl, url)
		return t.blockChainClient.SendRequestToChain(url, cluster.NodeToken, reqBody)
	}

	return "", errors.New("blockChainCode is error")
}

func (t *Tron) BalanceCluster(blockChain int64) *config.NodeCluster {
	cluster, ok := t.nodeCluster[blockChain]
	if !ok {
		//不存在节点
		return nil
	}
	//todo 后期重构节点筛选算法
	//根据 采集节点、任务节点的节点使用数据，综合判断出最佳节点
	//目前暂使用随机算法 找到节点
	if len(cluster) > 1 {
		l := len(cluster)
		return cluster[rand.Intn(l)]
	} else if len(cluster) == 1 {
		return cluster[0]
	} else {
		return nil
	}
}
