package btc

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

type Btc struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.ChainConn
}

func (e *Btc) Token(chainCode int64, contractAddr string, abi string, eip string) (string, error) {
	return "", nil
}

func (e *Btc) GetCode(chainCode int64, address string) (string, error) {
	return "", nil
}

func (e *Btc) GetAddressType(chainCode int64, address string) (string, error) {
	return "", nil
}

func (e *Btc) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", nil
}

// SubscribeLogs {"jsonrpc":"2.0","id": 1, "method": "eth_subscribe", "params": ["logs", {"address": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", "topics": ["0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"]}]}
func (e *Btc) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	return "", nil
}

func (e *Btc) UnSubscribe(chainCode int64, subId string) (string, error) {
	return "", nil
}

func (e *Btc) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	return "", nil
}

func (e *Btc) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	return "", nil
}

func (e *Btc) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	return "", nil
}

func (e *Btc) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	req := `
		{
		  "method": "getblock",
          "cache":true,
		  "params": [
			"%v",
 			%v
		  ]
		}`

	hasTx := 1
	if flag {
		hasTx = 2
	}

	req = fmt.Sprintf(req, hash, hasTx)
	return e.SendReq(chainCode, req)
}

func (e *Btc) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	req := `{
				"method": "getblockhash",
       		    "cache":true,
				"params": [
					%v
				]
			}
			`
	req = fmt.Sprintf(req, number)
	resp, err := e.SendReq(chainCode, req)

	if err != nil {
		return resp, err
	}

	hash := gjson.Parse(resp).Get("result").String()
	if len(hash) < 10 {
		return "", fmt.Errorf("not found blockHash by blockNumber:%v", number)
	}

	return e.GetBlockByHash(chainCode, hash, flag)
}

func (e *Btc) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `
			{
				"method": "getrawtransaction",
                "cache":true,
				"params": [
					"%v",
					2
				]
			}
		`
	req = fmt.Sprintf(req, hash)
	return e.SendReq(chainCode, req)
}

func (e *Btc) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendReq(chainCode, req)
}

func NewBtc(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain, xlog)
	if blockChainClient == nil {
		return nil
	}
	e := &Btc{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func (e *Btc) StartWDT() {
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

func (e *Btc) MonitorCluster() any {
	return e.nodeCluster
}

func (e *Btc) Balance(chainCode int64, address string, tag string) (string, error) {
	uris := make([]string, 0)
	for _, v := range e.nodeCluster {
		uris = append(uris, v.Utxo)
	}

	list, err := GetBalanceEx(uris, address)
	if err != nil {
		return "", err
	}

	bs, _ := json.Marshal(list)
	return string(bs), nil
}

func (e *Btc) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	return "", nil
}

func (e *Btc) Nonce(chainCode int64, address string, tag string) (string, error) {
	return "", nil
}

func (e *Btc) LatestBlock(chainCode int64) (string, error) {
	req := `{
				"method": "getblockcount"
			}
			`
	return e.SendReq(chainCode, req)
}

func (e *Btc) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	req := `{
			"method": "sendrawtransaction",
			"params": [
				"%v"
				]
			}`
	req = fmt.Sprintf(req, signedTx)
	return e.SendReq(chainCode, req)
}

func (e *Btc) SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error) {
	return "", nil
}

func (e *Btc) SendReq(blockChain int64, reqBody string) (resp string, err error) {
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

	cluster := e.BalanceCluster()
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

func (e *Btc) BalanceCluster() *config.NodeCluster {
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
