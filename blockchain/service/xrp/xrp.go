package xrp

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
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type XRP struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.ChainConn
}

func (e *XRP) Token(chainCode int64, contractAddr string, abi string, eip string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (e *XRP) GetCode(chainCode int64, address string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

func (e *XRP) GetAddressType(chainCode int64, address string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

func (e *XRP) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

// SubscribeLogs {"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs", {"address": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]}
func (e *XRP) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

func (e *XRP) UnSubscribe(chainCode int64, subId string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

func (e *XRP) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	req := `{
					"method": "ledger",
					"params": [
						{
							"ledger_index": "%v",
							"accounts": false,
							"full": false,
							"transactions": true,
							"expand": %v,
							"owner_funds": true
						}
					]
				}`
	req = fmt.Sprintf(req, number, true)
	resp, err := e.SendReq(chainCode, req, false)
	if err != nil {
		return "", err
	}

	arrays := gjson.Parse(resp).Get("result.ledger.transactions").Array()
	list := make([]any, 0, 5)
	for _, root := range arrays {
		r := make(map[string]any, 5)
		account := root.Get("Account").String()
		r["account"] = account
		hash := root.Get("hash").String()
		r["hash"] = hash
		date := root.Get("date").Int()
		r["date"] = date
		ledgerIndex := root.Get("ledger_index").Int()
		r["ledgerIndex"] = ledgerIndex

		transactionIndex := root.Get("meta.TransactionIndex").Int()
		r["transactionIndex"] = transactionIndex
		transactionResult := root.Get("meta.TransactionResult").String()
		r["transactionResult"] = transactionResult

		status := root.Get("status").String()
		r["status"] = status

		list = append(list, r)
	}

	bs, _ := json.Marshal(list)
	return string(bs), nil
}

func (e *XRP) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	req := `{
					"method": "ledger",
					"params": [
						{
							"ledger_hash": "%v",
							"accounts": false,
							"full": false,
							"transactions": true,
							"expand": %v,
							"owner_funds": true
						}
					]
				}`
	req = fmt.Sprintf(req, hash, true)
	resp, err := e.SendReq(chainCode, req, false)
	if err != nil {
		return "", err
	}

	arrays := gjson.Parse(resp).Get("result.ledger.transactions").Array()
	list := make([]any, 0, 5)
	for _, root := range arrays {
		r := make(map[string]any, 5)
		account := root.Get("Account").String()
		r["account"] = account
		hash := root.Get("hash").String()
		r["hash"] = hash
		ledgerIndex := root.Get("ledger_index").Int()
		r["ledgerIndex"] = ledgerIndex
		date := root.Get("date").Int()
		r["date"] = date
		transactionIndex := root.Get("meta.TransactionIndex").Int()
		r["transactionIndex"] = transactionIndex
		transactionResult := root.Get("meta.TransactionResult").String()
		r["transactionResult"] = transactionResult

		status := root.Get("status").String()
		r["status"] = status

		list = append(list, r)
	}

	bs, _ := json.Marshal(list)
	return string(bs), nil
}

func (e *XRP) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTransactionReceiptByHash,Duration=%v", time.Since(start))
	}()
	req := `{
		"method": "tx",
		"params": [
				{
					"transaction": "%v"
				}
			]
		}`
	req = fmt.Sprintf(req, hash)
	resp, err := e.SendReq(chainCode, req, false)
	if err != nil {
		return "", err
	}

	root := gjson.Parse(resp).Get("result")

	r := make(map[string]any, 5)
	account := root.Get("Account").String()
	r["account"] = account
	hash = root.Get("hash").String()
	r["hash"] = hash
	ledgerIndex := root.Get("ledger_index").Int()
	r["ledgerIndex"] = ledgerIndex
	date := root.Get("date").Int()
	r["date"] = date
	transactionIndex := root.Get("meta.TransactionIndex").Int()
	r["transactionIndex"] = transactionIndex
	transactionResult := root.Get("meta.TransactionResult").String()
	r["transactionResult"] = transactionResult

	status := root.Get("status").String()
	r["status"] = status
	bs, _ := json.Marshal(r)
	return string(bs), nil
}

func (e *XRP) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	req := `{
					"method": "ledger",
					"params": [
						{
							"ledger_hash": "%v",
							"accounts": false,
							"full": false,
							"transactions": true,
							"expand": %v,
							"owner_funds": true
						}
					]
				}`
	req = fmt.Sprintf(req, hash, flag)
	return e.SendReq(chainCode, req, false)
}

// GetBlockByNumber number: ledger version or shortcut string
func (e *XRP) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	req := `{
					"method": "ledger",
					"params": [
						{
							"ledger_index": "%v",
							"accounts": false,
							"full": false,
							"transactions": true,
							"expand": %v,
							"owner_funds": true
						}
					]
				}`
	req = fmt.Sprintf(req, number, flag)
	return e.SendReq(chainCode, req, false)
}

func (e *XRP) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `{
		"method": "tx",
		"params": [
				{
					"transaction": "%v"
				}
			]
		}`
	req = fmt.Sprintf(req, hash)
	return e.SendReq(chainCode, req, false)
}

func (e *XRP) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendReq(chainCode, req, false)
}

func NewXRP(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	e := &XRP{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func (e *XRP) StartWDT() {
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

func (e *XRP) MonitorCluster() any {
	return e.nodeCluster
}

// Balance ,tag:validated or closed or current or ledger version
func (e *XRP) Balance(chainCode int64, address string, tag string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("Balance,Duration=%v", time.Since(start))
	}()
	if len(tag) < 1 {
		tag = "validated"
	}
	req := `{
    "method": "account_info",
    "params": [
			{
				"account": "%v",
				"ledger_index": "%v",
				"queue": false
			}
  		  ]
		}`

	req = fmt.Sprintf(req, address, tag)
	resp, err := e.SendReq(chainCode, req, false)
	if err != nil {
		return "", err
	}

	return gjson.Parse(resp).Get("result.account_data").String(), nil
}

func (e *XRP) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
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

func (e *XRP) Nonce(chainCode int64, address string, tag string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", chainCode)
}

func (e *XRP) LatestBlock(chainCode int64) (string, error) {
	req := `{
				"method": "ledger_closed",
				"params": [
					{}
				]
			}`
	return e.SendReq(chainCode, req, false)
}

func (e *XRP) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
			"method": "submit",
			"params": [
							{
								"tx_blob": "%v"
							}
					  ]
			}`
	req = fmt.Sprintf(req, signedTx)
	return e.SendReq(chainCode, req, false)
}

func (e *XRP) SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", fmt.Errorf("blockchain:%v not implement the method", blockChain)
}

func (e *XRP) SendReq(blockChain int64, reqBody string, trace bool) (resp string, err error) {
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
	cluster := e.BalanceCluster(trace)
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	uri = fmt.Sprintf("%v/%v", cluster.NodeUrl, cluster.NodeToken)

	resp, err = e.blockChainClient.SendRequestToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	if err != nil {
		cluster.ErrorCount += 1
		//return "", errors.New("blockChainCode is error")
	}
	return resp, err
}

func (e *XRP) BalanceCluster(trace bool) *config.NodeCluster {
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
