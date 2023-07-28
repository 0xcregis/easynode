package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Tron struct {
	log              *xlog.XLog
	nodeCluster      []*config.NodeCluster
	blockChainClient blockchain.BlockChain
}

func (t *Tron) StartWDT() {
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for {
			<-ticker.C
			for _, v := range t.nodeCluster {
				v.ErrorCount = 0
			}
		}
	}()
}

func (t *Tron) MonitorCluster() any {
	return t.nodeCluster
}

func (t *Tron) GetCode(chainCode int64, address string) (string, error) {
	req := `{ "value": "%v", "visible": true}`
	req = fmt.Sprintf(req, address)
	return t.SendReq(chainCode, req, "wallet/getcontract")
}

func (t *Tron) GetAddressType(chainCode int64, address string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetAddressType,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v", "visible": true}`
	req = fmt.Sprintf(req, address)
	resp, err := t.SendReq(chainCode, req, "wallet/getcontract")
	if err != nil {
		return "", err
	}

	if gjson.Parse(resp).Get("code_hash").Exists() {
		//合约地址
		return "0x12", nil
	} else {
		//外部地址
		return "0x11", nil
	}
}

func (t *Tron) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) UnSubscribe(chainCode int64, subId string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockReceiptByBlockNumber,Duration=%v", time.Since(start))
	}()
	req := `{"num": %v}`

	n, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return "", err
	}
	req = fmt.Sprintf(req, n)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyblocknum")
}

func (t *Tron) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	return "", nil
}

func (t *Tron) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetTransactionReceiptByHash,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactioninfobyid")
}

func NewTron(cluster []*config.NodeCluster, blockchain int64, xlog *xlog.XLog) blockchain.API {
	blockChainClient := chain.NewChain(blockchain)
	if blockChainClient == nil {
		return nil
	}
	t := &Tron{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: blockChainClient,
	}
	t.StartWDT()
	return t
}

func (t *Tron) GetBlockByHash(chainCode int64, hash string, flag bool) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockByHash,Duration=%v", time.Since(start))
	}()

	req := `{"id_or_num": "%v","detail":%v}`
	req = fmt.Sprintf(req, hash, flag)
	res, err := t.SendReq(chainCode, req, "wallet/getblock")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	//mp := gjson.Parse(res).Map()
	//delete(mp, "transactions")
	//r, _ := json.Marshal(mp)
	//return string(r), nil
	//} else {
	return res, nil
	//}
}

func (t *Tron) GetBlockByNumber(chainCode int64, number string, flag bool) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetBlockByNumber,Duration=%v", time.Since(start))
	}()

	req := `{"id_or_num": "%v","detail":%v}`

	n, err := strconv.ParseInt(number, 0, 64)
	if err != nil {
		return "", err
	}
	req = fmt.Sprintf(req, n, flag)
	res, err := t.SendReq(chainCode, req, "wallet/getblock")
	if err != nil {
		return "", err
	}

	//var delTx bool = true
	//if delTx {
	//mp := gjson.Parse(res).Map()
	//delete(mp, "transactions")
	//r, _ := json.Marshal(mp)
	//return string(r), nil
	//} else {
	return res, nil
	//}
}

func (t *Tron) GetTxByHash(chainCode int64, hash string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("GetTxByHash,Duration=%v", time.Since(start))
	}()
	req := `{ "value": "%v"}`
	req = fmt.Sprintf(req, hash)
	return t.SendReq(chainCode, req, "wallet/gettransactionbyid")
}

func (t *Tron) SendJsonRpc(chainCode int64, req string) (string, error) {
	cluster := t.BalanceCluster()
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}
	url := fmt.Sprintf("%v/%v", cluster.NodeUrl, "jsonrpc")
	return t.blockChainClient.SendRequestToChain(url, cluster.NodeToken, req)
}

func (t *Tron) Balance(chainCode int64, address string, tag string) (string, error) {
	start := time.Now()
	defer func() {
		t.log.Printf("Balance,Duration=%v", time.Since(start))
	}()
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
	start := time.Now()
	defer func() {
		t.log.Printf("TokenBalance,Duration=%v", time.Since(start))
	}()
	cluster := t.BalanceCluster()
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
	cluster := t.BalanceCluster()
	if cluster == nil {
		//不存在节点
		return "", errors.New("blockchain node has not found")
	}

	if blockChain == 205 {
		url = fmt.Sprintf("%v/%v", cluster.NodeUrl, url)
		resp, err := t.blockChainClient.SendRequestToChainByHttp(url, cluster.NodeToken, reqBody)
		if err != nil {
			cluster.ErrorCount += 1
		}
		return resp, err
	}

	return "", errors.New("blockChainCode is error")
}

func (t *Tron) BalanceCluster() *config.NodeCluster {

	var resultCluster *config.NodeCluster
	l := len(t.nodeCluster)

	if l > 1 {
		//如果有多个节点，则根据权重计算
		mp := make(map[string][]int64, 0)
		originCluster := make(map[string]*config.NodeCluster, 0)

		var sum int64
		for _, v := range t.nodeCluster {
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
		resultCluster = t.nodeCluster[0]
	} else {
		return nil
	}
	return resultCluster
}
