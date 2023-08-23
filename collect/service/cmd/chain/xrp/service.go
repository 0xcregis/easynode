package xrp

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	chainService "github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Service struct {
	log                *xlog.XLog
	chain              *config.Chain
	nodeId             string
	transferTopic      string
	store              collect.StoreTaskInterface
	txChainClient      blockchain.API
	blockChainClient   blockchain.API
	receiptChainClient blockchain.API
}

func (s *Service) GetMultiBlockByNumber(blockNumber string, log *logrus.Entry, flag bool) ([]*collect.BlockInterface, []*collect.TxInterface) {
	log.Errorf("GetMultiBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "not found the method", blockNumber)
	return nil, nil
}

func (s *Service) Monitor() {
	go func() {
		for {
			<-time.After(30 * time.Second)
			if s.txChainClient != nil {
				txCluster := s.txChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "tx", txCluster)
			}

			if s.blockChainClient != nil {
				blockCluster := s.blockChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "block", blockCluster)
			}

			if s.receiptChainClient != nil {
				receiptCluster := s.receiptChainClient.MonitorCluster()
				_ = s.store.StoreClusterNode(int64(s.chain.BlockChainCode), "receipt", receiptCluster)
			}
		}
	}()
}

func (s *Service) GetBlockByHash(blockHash string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {
	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByHash.Duration =%v,blockHash:%v", time.Since(start), blockHash)
	}()
	//调用接口
	resp, err := s.blockChainClient.GetBlockByHash(int64(s.chain.BlockChainCode), blockHash, flag)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err.Error(), blockHash)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "block is empty", blockHash)
		return nil, nil
	}

	root := gjson.Parse(resp).Get("result")

	blockHash = root.Get("ledger_hash").String()
	if len(blockHash) < 1 {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, resp, blockHash)
		return nil, nil
	}
	block := root.Get("ledger").String()
	block = GetBlockHead(block)
	blockNumber := root.Get("ledger_index").String()

	r := &collect.BlockInterface{BlockHash: blockHash, BlockNumber: blockNumber, Block: block}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	txs := root.Get("ledger.transactions").Array()
	txList := make([]*collect.TxInterface, 0, len(txs))
	for _, v := range txs {
		hash := v.Get("hash").String()
		txStr := v.String()
		tx := &collect.TxInterface{TxHash: hash, Tx: txStr}
		txList = append(txList, tx)
	}

	for _, v := range txList {
		m := make(map[string]any, 2)
		m["tx"] = v.Tx
		m["hash"] = v.TxHash
		m["block"] = block
		m["blockNumber"] = blockNumber
		m["blockHash"] = blockHash
		//receipt, err := s.GetReceipt(v.TxHash, eLog)
		//if err != nil {
		//	eLog.Errorf("GetBlockByHash|BlockChainCode=%v,err=%v,blockHash=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockHash, v.TxHash)
		//} else {
		//	if receipt != nil {
		//		bs, _ := json.Marshal(receipt.Receipt)
		//		m["receipt"] = string(bs)
		//	}
		//}

		v.Tx = m
	}

	addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	addressMp := rebuildAddress(addressList)
	resultTxs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx.Tx)
		if s.CheckAddress(bs, addressMp) {
			//t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx.Tx}
			resultTxs = append(resultTxs, tx)
		}
	}
	return r, resultTxs
}

func (s *Service) GetBlockByNumber(blockNumber string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {

	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByNumber.Duration =%v,blockNumber:%v", time.Since(start), blockNumber)
	}()
	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber, flag)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	root := gjson.Parse(resp).Get("result")

	blockHash := root.Get("ledger_hash").String()
	if len(blockHash) < 1 {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, resp, blockNumber)
		return nil, nil
	}
	block := root.Get("ledger").String()
	block = GetBlockHead(block)
	blockNumber = root.Get("ledger_index").String()

	r := &collect.BlockInterface{BlockHash: blockHash, BlockNumber: blockNumber, Block: block}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	txs := root.Get("ledger.transactions").Array()
	txList := make([]*collect.TxInterface, 0, len(txs))
	for _, v := range txs {
		hash := v.Get("hash").String()
		txStr := v.String()
		tx := &collect.TxInterface{TxHash: hash, Tx: txStr}
		txList = append(txList, tx)
	}

	for _, v := range txList {
		m := make(map[string]any, 2)
		m["tx"] = v.Tx
		m["hash"] = v.TxHash
		m["block"] = block
		m["blockNumber"] = blockNumber
		m["blockHash"] = blockHash
		//receipt, err := s.GetReceipt(v.TxHash, eLog)
		//if err != nil {
		//	eLog.Errorf("GetBlockByNumber|BlockChainCode=%v,err=%v,blockNumber=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockNumber, v.TxHash)
		//} else {
		//	if receipt != nil {
		//		bs, _ := json.Marshal(receipt.Receipt)
		//		m["receipt"] = string(bs)
		//	}
		//}

		v.Tx = m
	}

	addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	addressMp := rebuildAddress(addressList)
	resultTxs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx.Tx)
		if s.CheckAddress(bs, addressMp) {
			//t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx.Tx}
			resultTxs = append(resultTxs, tx)
		}
	}
	return r, resultTxs
}

func (s *Service) GetTx(txHash string, eLog *logrus.Entry) *collect.TxInterface {

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
	root := gjson.Parse(resp).Get("result")

	//解析数据

	hash := root.Get("hash").String()
	if len(hash) < 1 {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil
	}

	m := make(map[string]any, 2)
	m["tx"] = root.String()
	m["hash"] = hash

	//receipt, err := s.GetReceipt(hash, eLog)
	//if err != nil {
	//	eLog.Errorf("GetTx|BlockChainCode=%v,err=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), hash)
	//} else {
	//	if receipt != nil {
	//		bs, _ := json.Marshal(receipt.Receipt)
	//		m["receipt"] = string(bs)
	//	}
	//}
	//bs, _ := json.Marshal(m)

	addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	addressMp := rebuildAddress(addressList)
	bs, _ := json.Marshal(m)
	if s.CheckAddress(bs, addressMp) {
		r := &collect.TxInterface{TxHash: hash, Tx: m}
		return r
	}
	//r := &collect.TxInterface{TxHash: hash, Tx: m}
	return nil
}

func (s *Service) GetReceiptByBlock(blockHash, number string, eLog *logrus.Entry) ([]*collect.ReceiptInterface, error) {

	//调用接口
	var resp string
	var err error
	if len(number) > 1 {
		resp, err = s.receiptChainClient.GetBlockReceiptByBlockNumber(int64(s.chain.BlockChainCode), number)
	} else if len(blockHash) > 1 {
		resp, err = s.receiptChainClient.GetBlockReceiptByBlockHash(int64(s.chain.BlockChainCode), blockHash)
	} else {
		err = fmt.Errorf("params is error")
	}

	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,blockHash=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockHash, number)
		return nil, err
	}

	//处理数据
	if len(resp) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, "receipt is empty", blockHash)
		return nil, errors.New("receipt is null")
	}

	// 解析数据
	roots := gjson.Parse(resp).Array()
	list := make([]*collect.ReceiptInterface, 0, len(roots))
	for _, root := range roots {
		TransactionHash := root.Get("hash").String()
		date := root.Get("date").Int()
		ledgerIndex := root.Get("ledgerIndex").Int()
		receipt := collect.ReceiptInterface{TransactionHash: TransactionHash, BlockNumber: ledgerIndex, BlockTimeStamp: date, Receipt: resp}
		list = append(list, &receipt)
	}
	return list, nil
}

func (s *Service) GetReceipt(txHash string, eLog *logrus.Entry) (*collect.ReceiptInterface, error) {

	//调用接口
	resp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), txHash)
	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	//处理数据
	if len(resp) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil, errors.New("receipt is null")
	}

	// 解析数据
	root := gjson.Parse(resp)
	TransactionHash := root.Get("hash").String()
	if len(TransactionHash) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil, errors.New("receipt is null")
	}
	date := root.Get("date").Int()
	ledgerIndex := root.Get("ledgerIndex").Int()

	return &collect.ReceiptInterface{TransactionHash: TransactionHash, BlockNumber: ledgerIndex, BlockTimeStamp: date, Receipt: resp}, nil
}

func getCoreAddr(addr string) string {
	addr = strings.ToLower(addr)
	if strings.HasPrefix(addr, "0x") {
		return strings.Replace(addr, "0x", "", 1) //去丢0x
	}
	return addr
}

func rebuildAddress(addrList []string) map[string]int64 {
	mp := make(map[string]int64, len(addrList))
	for _, v := range addrList {
		addr := getCoreAddr(v)
		mp[addr] = 1
	}
	return mp
}

func (s *Service) CheckAddress(tx []byte, addrList map[string]int64) bool {

	//if len(addrList) < 1 || len(tx) < 1 {
	//	return false
	//}

	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(tx)

	txBody := root.Get("tx").String()
	txRoot := gjson.Parse(txBody)

	fromAddr := txRoot.Get("Account").String()
	txAddressList[getCoreAddr(fromAddr)] = 1

	if txRoot.Get("Destination").Exists() {
		toAddr := txRoot.Get("Destination").String()
		txAddressList[getCoreAddr(toAddr)] = 1
	}

	//mp := make(map[string]int64, len(addrList))
	//for _, v := range addrList {
	//	addr := getCoreAddr(v)
	//	mp[addr] = 1
	//}

	has := false
	for k := range txAddressList {
		//monitorAddr := getCoreAddr(v)
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	return has
}

func NewService(c *config.Chain, x *xlog.XLog, store collect.StoreTaskInterface, nodeId string, transferTopic string) collect.BlockChainInterface {

	var blockClient blockchain.API
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
		blockClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	var txClient blockchain.API
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

		txClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	var receiptClient blockchain.API
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
		receiptClient = chainService.NewApi(int64(c.BlockChainCode), list, x)
	}

	return &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
		nodeId:             nodeId,
		transferTopic:      transferTopic,
	}
}
