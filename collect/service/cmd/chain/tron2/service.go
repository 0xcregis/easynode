package tron2

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
	"strconv"
	"strings"
)

// tron链 http 协议返回的数据

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

func (s *Service) GetTx(txHash string, task *config.TxTask, eLog *logrus.Entry) *service.TxInterface {

	//调用接口
	resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), txHash)
	if err != nil {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "tx is empty", txHash)
		return nil
	}

	hash := gjson.Parse(resp).Get("txID").String()
	receipt := s.GetReceipt(hash, nil, eLog)

	fullTx := make(map[string]interface{}, 2)
	fullTx["tx"] = resp
	fullTx["receipt"] = receipt.Receipt

	r := &service.TxInterface{TxHash: hash, Tx: fullTx}
	return r
}

func (s *Service) GetReceipt(txHash string, task *config.ReceiptTask, eLog *logrus.Entry) *service.ReceiptInterface {

	//调用接口
	resp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), txHash)
	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil
	}

	hash := gjson.Parse(resp).Get("id").String()

	var receipt service.TronReceipt
	_ = json.Unmarshal([]byte(resp), &receipt)

	s.buildContract(&receipt)

	bs, _ := json.Marshal(receipt)
	r := &service.ReceiptInterface{TransactionHash: hash, Receipt: string(bs)}
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
	} else if len(number) == 0 && len(blockHash) > 0 {
		resp, err = s.receiptChainClient.GetBlockReceiptByBlockHash(int64(s.chain.BlockChainCode), blockHash)
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

	txs := gjson.Parse(resp).Array()
	receiptList := make([]*service.ReceiptInterface, 0, len(txs))
	for _, v := range txs {
		hash := v.Get("id").String()

		var receipt service.TronReceipt
		_ = json.Unmarshal([]byte(v.String()), &receipt)

		s.buildContract(&receipt)

		bs, _ := json.Marshal(receipt)
		r := &service.ReceiptInterface{TransactionHash: hash, Receipt: string(bs)}
		receiptList = append(receiptList, r)
	}

	return receiptList
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry) (*service.BlockInterface, []*service.TxInterface) {
	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}

	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	blockID := gjson.Parse(resp).Get("blockID").String()
	number := gjson.Parse(resp).Get("block_header.raw_data.number").String()
	list := gjson.Parse(resp).Get("transactions").Array()

	//根据区块高度，获取交易log
	receipts := s.GetReceiptByBlock(blockID, number, nil, eLog)
	mp := make(map[string]interface{}, len(receipts))
	for _, v := range receipts {
		mp[v.TransactionHash] = v.Receipt
	}

	txs := make([]*service.TxInterface, 0, len(list))
	for _, tx := range list {
		// 补充字段
		hash := tx.Get("txID").String()
		fullTx := make(map[string]interface{}, 2)
		fullTx["tx"] = tx.String()
		if v, ok := mp[hash]; ok {
			fullTx["receipt"] = v
		}
		t := &service.TxInterface{TxHash: hash, Tx: fullTx}
		txs = append(txs, t)
	}
	r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}
	return r, txs
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry) (*service.BlockInterface, []*service.TxInterface) {
	//调用接口
	resp, err := s.blockChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), blockHash)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		return nil, nil
	}

	blockID := gjson.Parse(resp).Get("blockID").String()
	number := gjson.Parse(resp).Get("block_header.raw_data.number").String()

	//根据区块高度，获取交易log
	receipts := s.GetReceiptByBlock(blockID, number, nil, eLog)
	mp := make(map[string]interface{}, len(receipts))
	for _, v := range receipts {
		mp[v.TransactionHash] = v.Receipt
	}

	list := gjson.Parse(resp).Get("transactions").Array()
	txs := make([]*service.TxInterface, 0, len(list))
	for _, tx := range list {
		// 补充字段
		hash := tx.Get("txID").String()
		fullTx := make(map[string]interface{}, 2)
		fullTx["tx"] = tx.String()
		if v, ok := mp[hash]; ok {
			fullTx["receipt"] = v
		}
		t := &service.TxInterface{TxHash: hash, Tx: fullTx}
		txs = append(txs, t)
	}
	r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}
	return r, txs
}

func (s *Service) buildContract(receipt *service.TronReceipt) {
	for _, g := range receipt.Log {

		if len(g.Topics) < 3 || g.Topics[0] != service.TronTopic {
			continue
		}

		mp := make(map[string]interface{}, 2)
		contractAddr := fmt.Sprintf("41%v", g.Address)
		token, err := s.getToken(int64(s.chain.BlockChainCode), contractAddr, contractAddr)
		if err != nil {
			continue
		}
		m := gjson.Parse(token).Map()
		if v, ok := m["decimals"]; ok {
			mp["contractDecimals"] = v.String()
		} else {
			continue
		}

		mp["data"] = g.Data
		bs, _ := json.Marshal(mp)
		g.Data = string(bs)
	}
}

func (s *Service) getToken(blockChain int64, from string, contract string) (string, error) {

	token, err := s.store.GetContract(blockChain, contract)
	if err == nil {
		return token, nil
	}

	go func() {
		token, err = s.txChainClient.TokenBalance(blockChain, from, contract, "")
		if err != nil {
			s.log.Warnf("TokenBalance fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
			return
		}
		err = s.store.StoreContract(blockChain, contract, token)
		if err != nil {
			s.log.Warnf("StoreContract fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
		}
	}()

	return token, errors.New("waiting from network")
}

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	return nil, nil
}

func NewService(c *config.Chain, x *xlog.XLog, store service.StoreTaskInterface) service.BlockChainInterface {

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
		blockNodeCluster[205] = list
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
		txNodeCluster[205] = list
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
		receiptNodeCluster[205] = list
	}

	txClient := chainService.NewTron(txNodeCluster, x)
	blockClient := chainService.NewTron(blockNodeCluster, x)
	receiptClient := chainService.NewTron(receiptNodeCluster, x)
	return &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
	}
}
