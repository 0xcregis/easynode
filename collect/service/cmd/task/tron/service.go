package tron

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/net/tron"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/db"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	log              *xlog.XLog
	task             *db.Service
	chain            *config.Chain
	originChain      config.Chain
	blockClusters    []*config.FromCluster
	txClusters       []*config.FromCluster
	receiptsClusters []*config.FromCluster
}

func (s *Service) Monitor() {
	go func() {
		for true {
			<-time.After(10 * time.Minute)
			if s.originChain.BlockTask != nil {
				s.blockClusters = s.originChain.BlockTask.FromCluster
			}
			if s.originChain.TxTask != nil {
				s.txClusters = s.originChain.TxTask.FromCluster
			}
			if s.originChain.ReceiptTask != nil {
				s.receiptsClusters = s.originChain.ReceiptTask.FromCluster
			}
		}
	}()
}

func (s *Service) GetTx(txHash string, task *config.TxTask, eLog *logrus.Entry) *service.Tx {
	//选择节点
	cluster, err := s.BalanceCluster(txHash, s.txClusters)
	if err != nil {
		eLog.Errorf("BalanceCluster|GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}
	//调用接口
	resp, err := tron.Eth_GetTransactionByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetTransactionByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		cluster.ErrorCount++
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetTransactionByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", txHash)
		cluster.ErrorCount++
		return nil
	}

	cluster.ErrorCount = 0
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	tx := service.GetTxFromJson(resp)
	return tx
}

func (s *Service) GetReceipt(txHash string, task *config.ReceiptTask, eLog *logrus.Entry) *service.Receipt {
	//选择节点
	cluster, err := s.BalanceCluster(txHash, s.receiptsClusters)
	if err != nil {
		eLog.Errorf("GetReceipt|BalanceCluster|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}
	//调用接口
	resp, err := tron.Eth_GetTransactionReceiptByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetTransactionReceiptByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		cluster.ErrorCount++
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetTransactionReceiptByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		cluster.ErrorCount++
		return nil
	}

	cluster.ErrorCount = 0
	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	return service.GetReceiptFromJson(resp)
}

func (s *Service) GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, eLog *logrus.Entry) []*service.Receipt {

	//选择节点
	cluster, err := s.BalanceCluster(blockHash, s.receiptsClusters)
	if err != nil {
		eLog.Errorf("GetReceiptByBlock|BalanceCluster|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, err.Error(), number, blockHash)
		return nil
	}

	//调用接口
	var resp string
	if len(number) > 0 {
		if !strings.HasPrefix(number, "0x") {
			n, _ := strconv.ParseInt(number, 10, 64)
			number = fmt.Sprintf("0x%x", n)
		}
		resp, err = tron.Eth_GetBlockReceiptByBlockNumber(cluster.Host, cluster.Key, number, s.log)
	} else if len(number) == 0 && len(blockHash) > 0 {
		resp, err = tron.Eth_GetBlockReceiptByBlockHash(cluster.Host, cluster.Key, blockHash, s.log)
	}

	if err != nil {
		eLog.Errorf("Eth_GetBlockReceiptByBlockNumberOrEth_GetBlockReceiptByBlockHash|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, err.Error(), number, blockHash)
		cluster.ErrorCount++
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockReceiptByBlockNumberOrEth_GetBlockReceiptByBlockHash|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		cluster.ErrorCount++
		return nil
	}

	cluster.ErrorCount = 0
	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	return service.GetReceiptListFromJson(resp)
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry) (*service.Block, []*service.Tx) {
	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}
	//选择节点
	cluster, err := s.BalanceCluster(blockNumber, s.blockClusters)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BalanceCluster|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//调用接口
	resp, err := tron.Eth_GetBlockByNumber(cluster.Host, cluster.Key, blockNumber, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		cluster.ErrorCount++
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		cluster.ErrorCount++
		return nil, nil
	}

	cluster.ErrorCount = 0
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)
	return block, txList
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry) (*service.Block, []*service.Tx) {
	//选择节点
	cluster, err := s.BalanceCluster(blockHash, s.blockClusters)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BalanceCluster|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}
	//调用接口
	resp, err := tron.Eth_GetBlockByHash(cluster.Host, cluster.Key, blockHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		cluster.ErrorCount++
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		cluster.ErrorCount++
		return nil, nil
	}

	cluster.ErrorCount = 0
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)
	return block, txList
}

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	var cluster *config.FromCluster
	//clusterList := s.chain.BlockTask.FromCluster
	temp := make([]*config.FromCluster, 0, 10)
	for _, v := range clusterList {
		if v.ErrorCount < 50 {
			temp = append(temp, v)
		}
	}

	l := len(temp)
	if l > 1 {
		m := rand.Intn(l)
		cluster = temp[m]
	} else if l == 1 {
		cluster = temp[0]
	} else {
		return nil, errors.New("not found active cluster node when try exec task on node")
	}
	return cluster, nil
}

func NewService(c *config.Chain, taskDb *config.TaskDb, sourceDb *config.SourceDb, x *xlog.XLog) service.BlockChainInterface {
	t := db.NewTaskService(taskDb, sourceDb, x)

	var bk []*config.FromCluster
	if c.BlockTask != nil {
		bk = c.BlockTask.FromCluster
	} else {
		bk = nil
	}

	var tx []*config.FromCluster
	if c.TxTask != nil {
		tx = c.TxTask.FromCluster
	} else {
		tx = nil
	}

	var rp []*config.FromCluster
	if c.ReceiptTask != nil {
		rp = c.ReceiptTask.FromCluster
	} else {
		rp = nil
	}
	return &Service{
		log:              x,
		task:             t,
		chain:            c,
		originChain:      c.CopyChain(),
		blockClusters:    bk,
		txClusters:       tx,
		receiptsClusters: rp,
	}
}
