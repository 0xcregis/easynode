package filecoin

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/common/ethtypes"
	"github.com/ipfs/go-cid"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type FileCoin struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.ChainConn
}

func (e *FileCoin) Token(chainCode int64, contractAddr string, abi string, eip string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewFileCoin(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	e := &FileCoin{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func (e *FileCoin) GetCode(chainCode int64, address string) (string, error) {
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

func (e *FileCoin) GetAddressType(chainCode int64, address string) (string, error) {
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

func (e *FileCoin) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {

	query := `{"jsonrpc":  "2.0",  "id":  1,  "method":  "eth_subscribe",  "params":  ["newPendingTransactions"]}`

	var resp string
	var err error
	go func() {
		sendCh <- query
		resp, err = e.SendReqByWs(chainCode, receiverCh, sendCh)
		e.log.Errorf("func=%v,resp=%v,err=%v", "SubscribePendingTx", resp, err.Error())
	}()
	return resp, err
}

// SubscribeLogs {"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs", {"address": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]}
func (e *FileCoin) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	query := ` {"jsonrpc":  "2.0",  "id":  1,  "method":  "eth_subscribe",  "params":  ["logs",  {"address":  "%v",  "topics":  []}]}`
	query = fmt.Sprintf(query, address)
	var resp string
	var err error
	go func() {
		sendCh <- query
		resp, err = e.SendReqByWs(chainCode, receiverCh, sendCh)
		e.log.Errorf("func=%v,resp=%v,err=%v", "SubscribeLogs", resp, err.Error())
	}()
	return resp, err
}

func (e *FileCoin) UnSubscribe(chainCode int64, subId string) (string, error) {
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

func (e *FileCoin) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	return "", errors.New("does not support the method")
}

func (e *FileCoin) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	return "", errors.New("does not support the method")
}

func (e *FileCoin) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTransactionReceiptByHash,Duration=%v", time.Since(start))
	}()

	eth, err := ethtypes.EthHashFromCid(cid.MustParse(hash))
	if err != nil {
		log.Panicln(err)
	}

	query := `{
				"id": 1,
				"jsonrpc": "2.0",
				"method": "eth_getTransactionReceipt",
				"params": [
					"%v"
				]
			}`
	query = fmt.Sprintf(query, eth.String())
	return e.SendReq(chainCode, query)
}

func (e *FileCoin) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "Filecoin.ChainGetBlock",
		  "params": [
            {
            "/": "%v"
            }
		  ]
		}`

	req = fmt.Sprintf(req, hash)
	resp, err := e.SendReq(chainCode, req)
	if err != nil {
		return "", err
	}

	m := make(map[string]any, 2)
	blockHead := gjson.Parse(resp).Get("result").String()
	m["BlockHead"] = blockHead
	m["Cid"] = hash

	if flag {
		blockMessage, err := e.GetTxsByHash(chainCode, hash)
		if err != nil {
			m["BlockMessages"] = err
		} else {
			m["BlockMessages"] = gjson.Parse(blockMessage).Get("result.BlsMessages").String()
		}
	}

	bs, _ := json.Marshal(m)
	return string(bs), nil
}

func (e *FileCoin) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	req := `
			{
			  "id": 1,
			  "jsonrpc": "2.0",
			  "method": "Filecoin.ChainGetTipSetByHeight",
			  "params": [
				%v,
				null
			  ]
			}
			`
	req = fmt.Sprintf(req, number)
	resp, err := e.SendReq(chainCode, req)
	if err != nil {
		return "", err
	}

	//mp := make(map[string]any, 2)
	root := gjson.Parse(resp).Get("result")
	cids := root.Get("Cids").Array()
	//mp["Cids"] = cids

	blocks := root.Get("Blocks").Array()
	list := make([]any, 0, len(blocks))
	for index, r := range blocks {
		m := make(map[string]any, 2)
		m["BlockHead"] = r.String()
		blockId := cids[index].Get("/").String()
		m["Cid"] = blockId
		if flag { //blockhead
			blockMsg, err := e.GetTxsByHash(chainCode, blockId)
			if err != nil {
				m["BlockMessages"] = err
			} else {
				m["BlockMessages"] = gjson.Parse(blockMsg).Get("result.BlsMessages").String()
			}
		}
		list = append(list, m)
	}

	bs, _ := json.Marshal(list)
	return string(bs), nil
}

func (e *FileCoin) GetTxsByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "params": [
				{
           			 "/": "%v"
       			 }
		  ],
		  "method": "Filecoin.ChainGetBlockMessages"
		}
		`
	req = fmt.Sprintf(req, hash)
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `
		{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "params": [
				{
           			 "/": "%v"
       			 }
		  ],
		  "method": "Filecoin.ChainGetMessage"
		}
		`
	req = fmt.Sprintf(req, hash)
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) Nonce(chainCode int64, address string, tag string) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "params": [
					  "%v"
				 ],
				 "method": "Filecoin.MpoolGetNonce"
			}
			`
	req = fmt.Sprintf(req, address)
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) LatestBlock(chainCode int64) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "method": "eth_blockNumber"
			}
			`
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
					 "id": 1,
					 "jsonrpc": "2.0",
					 "params": [
						  %v
					 ],
					 "method": "Filecoin.MpoolPush"
				}`
	req = fmt.Sprintf(req, signedTx)
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) StartWDT() {
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

func (e *FileCoin) MonitorCluster() any {
	return e.nodeCluster
}

func (e *FileCoin) Balance(chainCode int64, address string, tag string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("Balance,Duration=%v", time.Since(start))
	}()
	//if len(tag) < 1 {
	//	tag = "latest"
	//}
	req := `{
			"jsonrpc": "2.0",
			"method": "Filecoin.WalletBalance",
			"id": 1,
			"params": [
				"%v"
				]
			}`

	req = fmt.Sprintf(req, address)
	return e.SendReq(chainCode, req)
}

func (e *FileCoin) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
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

func (e *FileCoin) SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error) {
	cluster := e.BalanceCluster(false)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	return "", errors.New("blockChainCode is error")
}

func (e *FileCoin) SendReq(blockChain int64, reqBody string) (resp string, err error) {
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

	//return "", errors.New("blockChainCode is error")
}

func (e *FileCoin) BalanceCluster(trace bool) *config.NodeCluster {
	var resultCluster *config.NodeCluster
	l := len(e.nodeCluster)

	if l > 1 {
		//如果有多个节点，则根据权重计算
		mp := make(map[string][]int64, 0)
		originCluster := make(map[string]*config.NodeCluster, 0)

		var sum int64
		for _, v := range e.nodeCluster {
			if v.Weight == 0 {
				//如果没有设置weight,则默认设定5
				v.Weight = 5
			}
			sum += v.Weight
			key := fmt.Sprintf("%v/%v", v.NodeUrl, v.NodeToken)
			mp[key] = []int64{v.Weight, sum}
			originCluster[key] = v
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
	} else {
		return nil
	}
	return resultCluster
}
