package service

import (
	"context"
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
	"github.com/gorilla/websocket"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type PolygonPos struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.BlockChain
}

func (e *PolygonPos) GetCode(chainCode int64, address string) (string, error) {
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

func (e *PolygonPos) GetAddressType(chainCode int64, address string) (string, error) {
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

func (e *PolygonPos) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {

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
func (e *PolygonPos) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
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

func (e *PolygonPos) UnSubscribe(chainCode int64, subId string) (string, error) {
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

func (e *PolygonPos) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
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

func (e *PolygonPos) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
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

func (e *PolygonPos) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
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

func (e *PolygonPos) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
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

func (e *PolygonPos) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
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
	req = fmt.Sprintf(req, number, flag)
	return e.SendReq(chainCode, req)
}

func (e *PolygonPos) GetTxByHash(chainCode int64, hash string) (string, error) {
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

func (e *PolygonPos) SendJsonRpc(chainCode int64, req string) (string, error) {
	return e.SendReq(chainCode, req)
}

func NewPolygonPos(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain)
	if blockChainClient == nil {
		return nil
	}
	e := &PolygonPos{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	e.StartWDT()
	return e
}

func (e *PolygonPos) StartWDT() {
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

func (e *PolygonPos) MonitorCluster() any {
	return e.nodeCluster
}

func (e *PolygonPos) Balance(chainCode int64, address string, tag string) (string, error) {
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

func (e *PolygonPos) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	start := time.Now()
	defer func() {
		e.log.Printf("TokenBalance,Duration=%v", time.Since(start))
	}()
	cluster := e.BalanceCluster()
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

func (e *PolygonPos) Nonce(chainCode int64, address string, tag string) (string, error) {
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

func (e *PolygonPos) LatestBlock(chainCode int64) (string, error) {
	req := `
			 {
				 "id": 1,
				 "jsonrpc": "2.0",
				 "method": "eth_blockNumber"
			}
			`
	return e.SendReq(chainCode, req)
}

func (e *PolygonPos) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
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

func (e *PolygonPos) SendReqByWs(blockChain int64, receiverCh chan string, sendCh chan string) (string, error) {
	cluster := e.BalanceCluster()
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	//if blockChain == 200 {
	host, err := e.blockChainClient.EthSubscribe(cluster.NodeUrl, cluster.NodeToken)
	if err != nil {
		return "", err
	}

	conn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	if err = conn.SetWriteDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatalf("SetWriteDeadline: %v", err)
		return "", err
	}

	if err = conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		log.Fatalf("SetReadDeadline: %v", err)
		return "", err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(sendCh chan string) {
		interrupt := true
		for interrupt {
			select {
			case sendMsg := <-sendCh:
				//订阅命令 或取消订阅命令
				if err := conn.WriteMessage(websocket.TextMessage, []byte(sendMsg)); err != nil {
					log.Fatalf("WriteMessage: %v", err)
				}
			case <-ctx.Done():
				interrupt = false
			}
		}
	}(sendCh)

	for {
		//接受消息
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Fatalf("ReadMessage: %v", err)
		}

		// 异常数据处理
		r := gjson.ParseBytes(p)
		if r.Get("error").Exists() {
			mp := make(map[string]interface{}, 0)
			mp["result"] = r.Get("error").String()
			mp["cmd"] = 0
			rs, _ := json.Marshal(mp)
			receiverCh <- string(rs)
			cancel()
			break
		}

		//订阅命名响应
		//{"id":1,"result":"0x9a52eeddc2b289f985c0e23a7d8427c8","jsonrpc":"2.0"}

		if r.Get("result").Exists() && strings.HasPrefix(r.Get("result").String(), "0x") {
			//订阅成功
			mp := make(map[string]interface{}, 0)
			mp["result"] = r.Get("result").String()
			mp["cmd"] = 1
			rs, _ := json.Marshal(mp)
			receiverCh <- string(rs)
			continue
		}

		//取消订阅命令响应
		//{ "id": 1, "result": true, "jsonrpc": "2.0" }
		if r.Get("result").IsBool() {
			if r.Get("result").Bool() {
				//取消订阅成功
				mp := make(map[string]interface{}, 0)
				mp["result"] = r.Get("result").String()
				mp["cmd"] = 3
				rs, _ := json.Marshal(mp)
				receiverCh <- string(rs)
				cancel()
				break
			} else {
				//取消订阅失败
				mp := make(map[string]interface{}, 0)
				mp["result"] = r.Get("result").String()
				mp["cmd"] = 2
				rs, _ := json.Marshal(mp)
				receiverCh <- string(rs)
				continue
			}
		}

		receiverCh <- string(p)
	}

	return "", nil
	//}
	//return "", errors.New("blockChainCode is error")
}

func (e *PolygonPos) SendReq(blockChain int64, reqBody string) (string, error) {
	cluster := e.BalanceCluster()
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	//if blockChain == 200 {
	resp, err := e.blockChainClient.SendRequestToChain(cluster.NodeUrl, cluster.NodeToken, reqBody)
	if err != nil {
		cluster.ErrorCount += 1
	}
	return resp, err
	//}

	//return "", errors.New("blockChainCode is error")
}

func (e *PolygonPos) BalanceCluster() *config.NodeCluster {
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
