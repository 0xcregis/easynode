package bnb

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/common/util"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Bnb struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.ChainConn
	nftClient        blockchain.NFT
}

func (e *Bnb) TokenURI(chainCode int64, contractAddress string, tokenId string, eip int64) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	resp, err := e.nftClient.TokenURI(cluster.NodeUrl, cluster.NodeToken, contractAddress, tokenId, eip)
	if err != nil {
		cluster.ErrorCount += 1
		return "", err
	}
	return resp, nil
}

func (e *Bnb) BalanceOf(chainCode int64, contractAddress string, address string, tokenId string, eip int64) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	resp, err := e.nftClient.BalanceOf(cluster.NodeUrl, cluster.NodeToken, contractAddress, address, tokenId, eip)
	if err != nil {
		cluster.ErrorCount += 1
		return "", err
	}
	return resp, nil
}

func (e *Bnb) OwnerOf(chainCode int64, contractAddress string, tokenId string, eip int64) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	resp, err := e.nftClient.OwnerOf(cluster.NodeUrl, cluster.NodeToken, contractAddress, tokenId, eip)
	if err != nil {
		cluster.ErrorCount += 1
		return "", err
	}
	return resp, nil
}

func (e *Bnb) TotalSupply(chainCode int64, contractAddress string, eip int64) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	resp, err := e.nftClient.TotalSupply(cluster.NodeUrl, cluster.NodeToken, contractAddress, eip)
	if err != nil {
		cluster.ErrorCount += 1
		return "", err
	}
	return resp, nil
}

func (e *Bnb) Token(chainCode int64, contractAddr string, abi string, eip string) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	var resp map[string]any
	var err error
	if eip == "721" {
		resp, err = e.blockChainClient.GetToken721(cluster.NodeUrl, cluster.NodeToken, contractAddr, "")
	} else if eip == "1155" {
		resp, err = e.blockChainClient.GetToken1155(cluster.NodeUrl, cluster.NodeToken, contractAddr, "")
	} else if eip == "20" {
		resp, err = e.blockChainClient.GetToken20(cluster.NodeUrl, cluster.NodeToken, contractAddr, "")
	} else {
		return "", fmt.Errorf("unknow the eip:%v", eip)
	}

	if err != nil {
		cluster.ErrorCount += 1
	}

	bs, _ := json.Marshal(resp)
	return string(bs), err
}

func (e *Bnb) GetCode(chainCode int64, address string) (string, error) {
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getCode",
				"params": [
					"%v",
					"latest"
				]
 			}`
	query = fmt.Sprintf(query, address)
	return e.SendReq(chainCode, query)
}

func (e *Bnb) GetAddressType(chainCode int64, address string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetAddressType,Duration=%v", time.Since(start))
	}()
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getCode",
				"params": [
					"%v",
					"latest"
				]
			}`
	query = fmt.Sprintf(query, address)
	resp, err := e.SendReq(chainCode, query)
	if err != nil {
		return "", err
	}

	code := gjson.Parse(resp).Get("result").String()
	if len(code) > 5 {
		//合约地址
		return "0x12", nil
	} else {
		//外部地址
		return "0x11", nil
	}
}

func (e *Bnb) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	//query := `{"jsonrpc":  "2.0",  "id":  1,  "method":  "eth_subscribe",  "params":  ["newPendingTransactions"]}`
	return "", fmt.Errorf("can not support sub for the blockchain:%v", chainCode)
}

// SubscribeLogs {"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs", {"address": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]}
func (e *Bnb) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	//query := ` {"jsonrpc":  "2.0",  "id":  1,  "method":  "eth_subscribe",  "params":  ["logs",  {"address":  "%v",  "topics":  []}]}`
	//query = fmt.Sprintf(query, address)
	return "", fmt.Errorf("can not support sub for the blockchain:%v", chainCode)
}

func (e *Bnb) UnSubscribe(chainCode int64, subId string) (string, error) {
	query := `{"id": 1, "method": "eth_unsubscribe", "params": ["%v"]}`

	query = fmt.Sprintf(query, subId)
	resp, err := e.SendReq(chainCode, query)
	if err != nil {
		return "", err
	}

	/**
	{
	  "jsonrpc": "2.0",
	  "id": 1,
	  "result": true
	}
	*/
	if gjson.Parse(resp).Get("result").Bool() {
		return resp, nil
	} else {
		return "", errors.New(resp)
	}
}

func (e *Bnb) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getBlockReceipts",
				"params": [
					"%v"
				]
			}`
	query = fmt.Sprintf(query, number)
	return e.SendReq(chainCode, query)
}

func (e *Bnb) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getBlockReceipts",
				"params": [
					"%v"
				]
			}`

	query = fmt.Sprintf(query, hash)
	return e.SendReq(chainCode, query)
}

func (e *Bnb) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTransactionReceiptByHash,Duration=%v", time.Since(start))
	}()
	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getTransactionReceipt",
				"params": [
					"%v"
				]
			}`
	query = fmt.Sprintf(query, hash)
	return e.SendReq(chainCode, query)
}

func (e *Bnb) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_getBlockByHash",
		  "params": [
			"%v",
 			%v
		  ]
		}`

	req = fmt.Sprintf(req, hash, flag)
	return e.SendReq(chainCode, req)
}

func (e *Bnb) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	req := `
			{
			  "id": 1,
			  "jsonrpc": "2.0",
			  "method": "eth_getBlockByNumber",
			  "params": [
				"%v",
				%v
			  ]
			}
			`
	number, _ = util.Int2Hex(number)
	req = fmt.Sprintf(req, number, flag)
	return e.SendReq(chainCode, req)
}

func (e *Bnb) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
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
	return e.SendReq(chainCode, req)
}

func (e *Bnb) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendReq(chainCode, req)
}

func NewNftBnb(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.NftApi {
	nftClient := chain.NewNFT(blockchain, xlog)
	if nftClient == nil {
		return nil
	}

	chain.NewNFT(blockchain, xlog)
	e := &Bnb{
		log:         xlog,
		nodeCluster: cluster,
		nftClient:   nftClient,
	}
	e.StartWDT()
	return e
}

func NewBnb(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	e := &Bnb{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func NewBnb2(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.ExApi {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	e := &Bnb{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func (e *Bnb) StartWDT() {
	go func() {
		t := time.NewTicker(10 * time.Minute)
		for {
			<-t.C
			for _, v := range e.nodeCluster {
				v.ErrorCount = 0
			}
		}
	}()
}

func (e *Bnb) MonitorCluster() any {
	return e.nodeCluster
}

func (e *Bnb) TraceTransaction(chainCode int64, address string) (string, error) {
	req := `{
					 "id": 1,
					 "jsonrpc": "2.0",
					 "params": [
						  "%v"
					 ],
					 "method": "trace_transaction"
				}`
	req = fmt.Sprintf(req, address)

	return e.SendReq2(chainCode, req)
}

func (e *Bnb) Balance(chainCode int64, address string, tag string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("Balance,Duration=%v", time.Since(start))
	}()
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
	return e.SendReq(chainCode, req)
}

func (e *Bnb) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("TokenBalance,Duration=%v", time.Since(start))
	}()
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	mp, err := e.blockChainClient.GetToken20(cluster.NodeUrl, cluster.NodeToken, contractAddr, address)
	if err != nil {
		return "", err
	}
	rs, _ := json.Marshal(mp)
	return string(rs), nil
}

func (e *Bnb) Nonce(chainCode int64, address string, tag string) (string, error) {
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
	return e.SendReq(chainCode, req)
}

func (e *Bnb) LatestBlock(chainCode int64) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "method": "eth_blockNumber"
			}
			`
	return e.SendReq(chainCode, req)
}

func (e *Bnb) GetAccountResourceForTron(chainCode int64, address string) (string, error) {
	return "", nil
}

func (e *Bnb) EstimateGasForTron(chainCode int64, from, to, functionSelector, parameter string) (string, error) {
	return "", nil
}

func (e *Bnb) EstimateGas(chainCode int64, from, to, data string) (string, error) {
	req := `
	{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_estimateGas",
		  "params": [
			{
			  "from":"%v",
			  "to": "%v",
			  "data": "%v"
			}
		  ]
	}`
	req = fmt.Sprintf(req, from, to, data)
	res, err := e.SendJsonRpc(chainCode, req)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (e *Bnb) GasPrice(chainCode int64) (string, error) {
	req := `
		{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "eth_gasPrice"
		}
		`

	//req = fmt.Sprintf(req, from, to, functionSelector, parameter)
	res, err := e.SendJsonRpc(chainCode, req)
	if err != nil {
		return "", err
	}
	return res, nil
}

func (e *Bnb) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
					 "id": 1,
					 "jsonrpc": "2.0",
					 "params": [
						  "%v"
					 ],
					 "method": "eth_sendRawTransaction"
				}`
	req = fmt.Sprintf(req, signedTx)
	return e.SendReq(chainCode, req)
}

func (e *Bnb) SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", nil
}

func (e *Bnb) SendReq(blockChain int64, reqBody string) (resp string, err error) {
	reqBody = strings.Replace(reqBody, "\t", "", -1)
	reqBody = strings.Replace(reqBody, "\n", "", -1)
	var uri string
	defer func() {
		if err != nil {
			e.log.Errorf("method:%v,blockChain:%v,req:%v,err:%v,uri:%v", "SendReq", blockChain, reqBody, err, uri)
		} else {
			e.log.Printf("method:%v,blockChain:%v,req:%v,resp:%v", "SendReq", blockChain, reqBody, "ok")
		}
	}()
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	uri = fmt.Sprintf("%v/%v", cluster.NodeUrl, cluster.NodeToken)

	resp, err = e.blockChainClient.SendRequestToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	if err != nil {
		cluster.ErrorCount += 1
	}
	return resp, err
}

func (e *Bnb) SendReq2(blockChain int64, reqBody string) (resp string, err error) {
	reqBody = strings.Replace(reqBody, "\t", "", -1)
	reqBody = strings.Replace(reqBody, "\n", "", -1)
	var uri string
	defer func() {
		if err != nil {
			e.log.Errorf("method:%v,blockChain:%v,req:%v,err:%v,uri:%v", "SendReq", blockChain, reqBody, err, uri)
		} else {
			e.log.Printf("method:%v,blockChain:%v,req:%v,resp:%v", "SendReq", blockChain, reqBody, "ok")
		}
	}()
	cluster := e.BalanceCluster(true)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	uri = fmt.Sprintf("%v/%v", cluster.NodeUrl, cluster.NodeToken)

	resp, err = e.blockChainClient.SendRequestToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	if err != nil {
		cluster.ErrorCount += 1
	}
	return resp, err
}

func (e *Bnb) BalanceCluster(trace bool) *config.NodeCluster {
	var resultCluster *config.NodeCluster
	l := len(e.nodeCluster)

	if l > 1 {
		//如果有多个节点，则根据权重计算
		mp := make(map[string][]int64, 0)
		originCluster := make(map[string]*config.NodeCluster, 1)

		var sum int64
		for _, v := range e.nodeCluster {
			if v.Weight == 0 {
				//如果没有设置weight,则默认设定5
				v.Weight = 5
			}
			if !trace || trace && v.Trace {
				sum += v.Weight
				key := fmt.Sprintf("%v/%v", v.NodeUrl, v.NodeToken)
				mp[key] = []int64{v.Weight, sum}
				originCluster[key] = v
			}
		}

		f := math.Mod(float64(time.Now().Unix()), float64(sum))
		var nodeId string

		for k, v := range mp {
			if len(v) == 2 && f <= float64(v[1]) && f >= float64(v[1]-v[0]) {
				nodeId = k
				break
			}
		}
		resultCluster = originCluster[nodeId]
	} else if l == 1 {
		//如果 仅有一个节点，则只能使用该节点
		resultCluster = e.nodeCluster[0]
		if trace && !resultCluster.Trace {
			return nil
		}
	} else {
		return nil
	}
	return resultCluster
}
