package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	"google.golang.org/grpc"
	"log"
	"math"
	"time"
)

// grpc.trongrid.io:50051
type Tron2 struct {
	log              *xlog.XLog
	nodeCluster      map[int64][]*config.NodeCluster
	blockChainClient *client.GrpcClient
}

func (t *Tron2) Balance(chainCode int64, address string, tag string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) TokenBalance(chainCode int64, address string, contractAddr string, abi string) (string, error) {
	c := t.GetGrpcClient(chainCode)
	if c == nil {
		return "", errors.New("conn error")
	}
	mp := make(map[string]interface{}, 2)
	balance, err := c.TRC20ContractBalance(address, contractAddr)

	if err != nil {
		log.Println("err=", err)
	} else {
		mp["balance"] = balance.String()
	}

	name, err := c.TRC20GetName(contractAddr)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["name"] = name
	}

	symbol, err := c.TRC20GetSymbol(contractAddr)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["symbol"] = symbol
	}

	decimals, err := c.TRC20GetDecimals(contractAddr)
	if err != nil {
		log.Println("err=", err)
	} else {
		mp["decimals"] = decimals
	}
	bs, _ := json.Marshal(mp)
	return string(bs), nil
}

func (t *Tron2) Nonce(chainCode int64, address string, tag string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) LatestBlock(chainCode int64) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) SendRawTransaction(chainCode int64, signedTx string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetBlockByHash(chainCode int64, hash string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetBlockByNumber(chainCode int64, number string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetTxByHash(chainCode int64, hash string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) SendJsonRpc(chainCode int64, req string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetBlockReceiptByBlockNumber(chainCode int64, number string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetBlockReceiptByBlockHash(chainCode int64, hash string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetTransactionReceiptByHash(chainCode int64, hash string) (string, error) {
	c := t.GetGrpcClient(chainCode)
	if c == nil {
		return "", errors.New("conn error")
	}

	tx, err := c.GetTransactionInfoByID(hash)
	if err != nil {
		return "", err
	}

	bs, _ := json.Marshal(tx)
	return string(bs), nil
}

func (t *Tron2) SubscribePendingTx(chainCode int64, receiverCh chan string, sendCh chan string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) SubscribeLogs(chainCode int64, address string, topics []string, receiverCh chan string, sendCh chan string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) UnSubscribe(chainCode int64, subId string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (t *Tron2) GetAddressType(chainCode int64, address string) (string, error) {
	c := t.GetGrpcClient(chainCode)
	if c == nil {
		return "", errors.New("conn error")
	}

	abi, err := t.blockChainClient.GetContractABI(address)
	if err != nil {
		return "", err
	}

	if abi != nil {
		//合约地址
		return "0x12", nil
	} else {
		//外部地址
		return "0x11", nil
	}
}

func (t *Tron2) GetCode(chainCode int64, address string) (string, error) {
	c := t.GetGrpcClient(chainCode)
	if c == nil {
		return "", errors.New("conn error")
	}

	abi, err := t.blockChainClient.GetContractABI(address)
	if err != nil {
		return "", err
	}

	bs, _ := json.Marshal(abi)
	return string(bs), nil
}

func NewTron2(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) *Tron2 {
	return &Tron2{
		log:              xlog,
		nodeCluster:      cluster,
		blockChainClient: nil,
	}
}

func (t *Tron2) GetGrpcClient(blockChain int64) *client.GrpcClient {
	if t.blockChainClient != nil {
		return t.blockChainClient
	}

	nc := t.BalanceCluster(blockChain)

	conn := client.NewGrpcClient(nc.NodeUrl)
	_ = conn.SetAPIKey(nc.NodeToken) // todo 没有发现设置意义
	err := conn.Start(grpc.WithInsecure())
	if err != nil {
		t.log.Errorf("GetGrpcClient|error=%v", err.Error())
		return nil
	}
	t.blockChainClient = conn
	return conn
}

func (t *Tron2) BalanceCluster(blockChain int64) *config.NodeCluster {
	cluster, ok := t.nodeCluster[blockChain]
	if !ok {
		//不存在节点
		return nil
	}

	var resultCluster *config.NodeCluster
	l := len(cluster)

	if l > 1 {
		//如果有多个节点，则根据权重计算
		mp := make(map[string][]int64, 0)
		originCluster := make(map[string]*config.NodeCluster, 0)

		var sum int64
		for _, v := range cluster {
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
		resultCluster = cluster[0]
	} else {
		return nil
	}
	return resultCluster
}
