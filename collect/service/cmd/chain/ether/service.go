package ether

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	chainConfig "github.com/uduncloud/easynode/blockchain/config"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	log                *xlog.XLog
	chain              *config.Chain
	txChainClient      chainService.API
	blockChainClient   chainService.API
	receiptChainClient chainService.API
}

func (s *Service) Monitor() {
}

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	return nil, nil
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry) (*service.BlockInterface, []*service.TxInterface) {
	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByHash.Duration =%v", time.Now().Sub(start))
	}()
	//调用接口
	resp, err := s.blockChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), blockHash)
	//resp, err := ether.Eth_GetBlockByHash(cluster.Host, cluster.Key, blockHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)

	list := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, nil, eLog)

	for _, v := range txList {
		for _, r := range list {
			if v.TxHash == r.TransactionHash {
				bs, _ := json.Marshal(r.Receipt)
				v.Receipt = string(bs)
				break
			}
		}
	}
	txs := make([]*service.TxInterface, 0, len(txList))
	for _, tx := range txList {
		//	// 补充字段
		//
		//	tp, err := s.txChainClient.GetAddressType(int64(s.chain.BlockChainCode), tx.ToAddr)
		//	if err == nil {
		//		tx.Type = tp
		//	}
		//
		//	rp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), tx.TxHash)
		//
		//	if err == nil {
		//		tx.Receipt = rp
		//	}
		t := &service.TxInterface{TxHash: tx.TxHash, Tx: tx}
		txs = append(txs, t)
	}
	r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry) (*service.BlockInterface, []*service.TxInterface) {

	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByNumber.Duration =%v", time.Now().Sub(start))
	}()

	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}

	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber)
	//resp, err := ether.Eth_GetBlockByNumber(cluster.Host, cluster.Key, blockNumber, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)

	list := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, nil, eLog)
	for _, v := range txList {
		for _, r := range list {
			if v.TxHash == r.TransactionHash {
				bs, _ := json.Marshal(r.Receipt)
				v.Receipt = string(bs)
				break
			}
		}
	}
	txs := make([]*service.TxInterface, 0, len(txList))
	for _, tx := range txList {
		//	// 补充字段
		//
		//	tp, err := s.txChainClient.GetAddressType(int64(s.chain.BlockChainCode), tx.ToAddr)
		//	if err == nil {
		//		tx.Type = tp
		//	}
		//
		//	rp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), tx.TxHash)
		//
		//	if err == nil {
		//		tx.Receipt = rp
		//	}
		t := &service.TxInterface{TxHash: tx.TxHash, Tx: tx}
		txs = append(txs, t)
	}
	r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetTx(txHash string, task *config.TxTask, eLog *logrus.Entry) *service.TxInterface {

	//调用接口
	resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), txHash)
	//resp, err := ether.Eth_GetTransactionByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetTransactionByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetTransactionByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", txHash)
		return nil
	}
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	tx := service.GetTxFromJson(resp)

	// 补充字段

	tp, err := s.txChainClient.GetAddressType(int64(s.chain.BlockChainCode), tx.ToAddr)
	if err == nil {
		tx.Type = tp
	}

	rp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), tx.TxHash)

	if err == nil {
		tx.Receipt = rp
	}
	r := &service.TxInterface{TxHash: tx.TxHash, Tx: tx}
	return r
}

func (s *Service) GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, eLog *logrus.Entry) []*service.ReceiptInterface {

	//调用接口
	var resp string
	var err error
	if len(number) > 0 {
		if !strings.HasPrefix(number, "0x") {
			n, _ := strconv.ParseInt(number, 10, 64)
			number = fmt.Sprintf("0x%x", n)
		}
		resp, err = s.receiptChainClient.GetBlockReceiptByBlockNumber(int64(s.chain.BlockChainCode), number)
		//resp, err = ether.Eth_GetBlockReceiptByBlockNumber(cluster.Host, cluster.Key, number, s.log)
	} else if len(number) == 0 && len(blockHash) > 0 {
		resp, err = s.receiptChainClient.GetBlockReceiptByBlockHash(int64(s.chain.BlockChainCode), blockHash)
		//resp, err = ether.Eth_GetBlockReceiptByBlockHash(cluster.Host, cluster.Key, blockHash, s.log)
	}

	if err != nil {
		eLog.Errorf("Eth_GetBlockReceiptByBlockNumberOrEth_GetBlockReceiptByBlockHash|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, err.Error(), number, blockHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetBlockReceiptByBlockNumberOrEth_GetBlockReceiptByBlockHash|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		return nil
	}
	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	receiptList := service.GetReceiptListFromJson(resp)
	rs := make([]*service.ReceiptInterface, 0, len(receiptList))
	for _, v := range receiptList {
		r := &service.ReceiptInterface{TransactionHash: v.TransactionHash, Receipt: v}
		rs = append(rs, r)
	}
	return rs
}

func (s *Service) GetReceipt(txHash string, task *config.ReceiptTask, eLog *logrus.Entry) *service.ReceiptInterface {

	//调用接口
	resp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), txHash)
	//resp, err := ether.Eth_GetTransactionReceiptByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("Eth_GetTransactionReceiptByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("Eth_GetTransactionReceiptByHash|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	receipt := service.GetReceiptFromJson(resp)

	r := &service.ReceiptInterface{TransactionHash: receipt.TransactionHash, Receipt: receipt}
	return r
}

func NewService(c *config.Chain, x *xlog.XLog) service.BlockChainInterface {
	blockNodeCluster := map[int64][]*chainConfig.NodeCluster{}
	if c.BlockTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.BlockTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}
		blockNodeCluster[200] = list
	}

	txNodeCluster := map[int64][]*chainConfig.NodeCluster{}
	if c.TxTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.TxTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}
		txNodeCluster[200] = list
	}

	receiptNodeCluster := map[int64][]*chainConfig.NodeCluster{}
	if c.ReceiptTask != nil {
		list := make([]*chainConfig.NodeCluster, 0, 4)
		for _, v := range c.ReceiptTask.FromCluster {
			temp := &chainConfig.NodeCluster{
				NodeUrl:   v.Host,
				NodeToken: v.Key,
				Weight:    v.Weight,
			}
			list = append(list, temp)
		}
		receiptNodeCluster[200] = list
	}

	txClient := chainService.NewEth(txNodeCluster, x)
	blockClient := chainService.NewEth(blockNodeCluster, x)
	receiptClient := chainService.NewEth(receiptNodeCluster, x)

	return &Service{
		log:                x,
		chain:              c,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
	}
}
