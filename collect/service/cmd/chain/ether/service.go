package ether

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	chainConfig "github.com/uduncloud/easynode/blockchain/config"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/common/util"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	log                *xlog.XLog
	chain              *config.Chain
	store              service.StoreTaskInterface
	txChainClient      chainService.API
	blockChainClient   chainService.API
	receiptChainClient chainService.API
}

func (s *Service) Monitor() {
}

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	return nil, nil
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry, flag bool) (*service.BlockInterface, []*service.TxInterface) {
	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByHash.Duration =%v", time.Now().Sub(start))
	}()
	//调用接口
	resp, err := s.blockChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), blockHash, flag)
	//resp, err := ether.Eth_GetBlockByHash(cluster.Host, cluster.Key, blockHash, s.log)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)
	r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}
	//list := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, nil, eLog)

	for _, v := range txList {
		receipt := s.GetReceipt(v.TxHash, nil, eLog)
		//for _, r := range list {
		if receipt != nil && v.TxHash == receipt.TransactionHash {
			bs, _ := json.Marshal(receipt.Receipt)
			v.Receipt = string(bs)
		}
		//}
	}
	txs := make([]*service.TxInterface, 0, len(txList))
	for _, tx := range txList {
		t := &service.TxInterface{TxHash: tx.TxHash, Tx: tx}
		txs = append(txs, t)
	}
	//r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry, flag bool) (*service.BlockInterface, []*service.TxInterface) {

	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByNumber.Duration =%v", time.Now().Sub(start))
	}()

	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}

	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber, flag)
	//resp, err := ether.Eth_GetBlockByNumber(cluster.Host, cluster.Key, blockNumber, s.log)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := service.GetBlockFromJson(resp)
	r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	if !flag { //仅区块数据，不涉及交易
		return r, nil
	}

	//list := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, nil, eLog)
	for _, v := range txList {
		receipt := s.GetReceipt(v.TxHash, nil, eLog)
		//for _, r := range list {
		if receipt != nil && v.TxHash == receipt.TransactionHash {
			bs, _ := json.Marshal(receipt.Receipt)
			v.Receipt = string(bs)
		}
		//}
	}
	txs := make([]*service.TxInterface, 0, len(txList))
	for _, tx := range txList {
		t := &service.TxInterface{TxHash: tx.TxHash, Tx: tx}
		txs = append(txs, t)
	}
	//r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return r, txs
}

func (s *Service) GetTx(txHash string, task *config.TxTask, eLog *logrus.Entry) *service.TxInterface {

	//调用接口
	resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), txHash)
	//resp, err := ether.Eth_GetTransactionByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", txHash)
		return nil
	}
	resp = gjson.Parse(resp).Get("result").String()

	//解析数据
	tx := service.GetTxFromJson(resp)

	rp := s.GetReceipt(tx.TxHash, nil, eLog)

	if rp != nil {
		bs, _ := json.Marshal(rp.Receipt)
		tx.Receipt = string(bs)
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
		eLog.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, err.Error(), number, blockHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		return nil
	}
	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	receiptList := service.GetReceiptListFromJson(resp)
	rs := make([]*service.ReceiptInterface, 0, len(receiptList))
	for _, v := range receiptList {
		s.buildContract(v)
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
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil
	}

	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	receipt := service.GetReceiptFromJson(resp)
	s.buildContract(receipt)
	r := &service.ReceiptInterface{TransactionHash: receipt.TransactionHash, Receipt: receipt}
	return r
}

func (s *Service) buildContract(receipt *service.Receipt) {
	for _, g := range receipt.Logs {

		if len(g.Topics) < 3 || g.Topics[0] != service.EthTopic {
			continue
		}

		//过滤721协议
		if len(g.Topics) == 4 && g.Topics[0] == service.EthTopic {
			if len(g.Data) == 2 {
				continue
			}
		}

		mp := make(map[string]interface{}, 2)
		token, err := s.getToken(int64(s.chain.BlockChainCode), receipt.From, g.Address)
		if err != nil {
			nodeId, _ := util.GetLocalNodeId()
			task := service.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: nodeId, TxHash: receipt.TransactionHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
			_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), receipt.TransactionHash, task)
			break
		}

		m := gjson.Parse(token).Map()
		if v, ok := m["decimals"]; ok {
			mp["contractDecimals"] = v.String()
		} else {
			nodeId, _ := util.GetLocalNodeId()
			task := service.NodeTask{Id: time.Now().UnixNano(), NodeId: nodeId, BlockChain: s.chain.BlockChainCode, TxHash: receipt.TransactionHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
			_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), receipt.TransactionHash, task)
			continue
		}

		mp["data"] = g.Data
		bs, _ := json.Marshal(mp)
		g.Data = string(bs)
	}
}

func (s *Service) getToken(blockChain int64, from string, contract string) (string, error) {
	token, err := s.store.GetContract(blockChain, contract)
	if err == nil && len(token) > 5 {
		return token, nil
	}

	token, err = s.txChainClient.TokenBalance(blockChain, from, contract, "")
	if err != nil {
		s.log.Errorf("TokenBalance fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
		//请求失败的合约记录，监控服务重试
		_ = s.store.StoreContract(blockChain, contract, "")
		return "", err
	}

	if len(token) > 0 {
		err = s.store.StoreContract(blockChain, contract, token)
		if err != nil {
			s.log.Warnf("StoreContract fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
		}
		return token, nil
	}
	return "", errors.New("wait for response")
}

func NewService(c *config.Chain, x *xlog.XLog, store service.StoreTaskInterface) service.BlockChainInterface {

	var blockClient chainService.API
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
		blockClient = chainService.NewEth(list, x)
	}

	var txClient chainService.API
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

		txClient = chainService.NewEth(list, x)
	}

	var receiptClient chainService.API
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
		receiptClient = chainService.NewEth(list, x)
	}

	return &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
	}
}
