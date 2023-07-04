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
	"github.com/uduncloud/easynode/common/util"
	"strconv"
	"strings"
	"time"
)

// tron链 http 协议返回的数据

type Service struct {
	log                *xlog.XLog
	chain              *config.Chain
	nodeId             string
	store              service.StoreTaskInterface
	txChainClient      chainService.API
	blockChainClient   chainService.API
	receiptChainClient chainService.API
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

	fullTx := make(map[string]interface{}, 2)
	fullTx["tx"] = resp

	receipt, err := s.GetReceipt(hash, nil, eLog)
	if err != nil {
		eLog.Errorf("GetTx|BlockChainCode=%v,err=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), txHash)
	} else {
		if receipt != nil {
			fullTx["receipt"] = receipt.Receipt
		}
	}

	//tron 交易中缺失 完整的blockHash
	if receipt != nil {
		block, _ := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), fmt.Sprintf("%v", receipt.BlockNumber), false)
		if len(block) > 0 {
			blockId := gjson.Parse(block).Get("blockID").String()
			fullTx["blockId"] = blockId
		}
	}

	r := &service.TxInterface{TxHash: hash, Tx: fullTx}
	return r
}

func (s *Service) GetReceipt(txHash string, task *config.ReceiptTask, eLog *logrus.Entry) (*service.ReceiptInterface, error) {

	//调用接口
	resp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), txHash)
	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil, errors.New("receipt is null")
	}

	hash := gjson.Parse(resp).Get("id").String()

	var receipt service.TronReceipt
	_ = json.Unmarshal([]byte(resp), &receipt)

	p := s.buildContract(&receipt)

	if p != nil {
		bs, _ := json.Marshal(receipt)
		r := &service.ReceiptInterface{TransactionHash: hash, Receipt: string(bs), BlockNumber: receipt.BlockNumber, BlockTimeStamp: receipt.BlockTimeStamp}
		return r, nil
	} else {
		//nodeId, _ := util.GetLocalNodeId()
		task := service.NodeTask{Id: time.Now().UnixNano(), NodeId: s.nodeId, BlockChain: s.chain.BlockChainCode, TxHash: receipt.Id, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
		_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), receipt.Id, task)
		return nil, errors.New("receipt is null")
	}
}

func (s *Service) GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, eLog *logrus.Entry) ([]*service.ReceiptInterface, error) {

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
		return nil, err
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		return nil, errors.New("receipt is null")
	}

	txs := gjson.Parse(resp).Array()
	receiptList := make([]*service.ReceiptInterface, 0, len(txs))
	for _, v := range txs {
		hash := v.Get("id").String()

		var receipt service.TronReceipt
		_ = json.Unmarshal([]byte(v.String()), &receipt)

		p := s.buildContract(&receipt)

		if p != nil {
			bs, _ := json.Marshal(p)
			r := &service.ReceiptInterface{TransactionHash: hash, Receipt: string(bs)}
			receiptList = append(receiptList, r)
		} else {
			//nodeId, _ := util.GetLocalNodeId()
			task := service.NodeTask{Id: time.Now().UnixNano(), NodeId: s.nodeId, BlockChain: s.chain.BlockChainCode, TxHash: receipt.Id, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
			_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), receipt.Id, task)
		}
	}

	return receiptList, nil
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry, flag bool) (*service.BlockInterface, []*service.TxInterface) {
	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}

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

	blockID := gjson.Parse(resp).Get("blockID").String()
	number := gjson.Parse(resp).Get("block_header.raw_data.number").String()

	if len(blockID) == 0 && len(number) == 0 {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,resp=%v", s.chain.BlockChainName, "resp is error", resp)
		return nil, nil
	}

	r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}

	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	list := gjson.Parse(resp).Get("transactions").Array()
	//根据区块高度，获取交易log
	mp := make(map[string]interface{}, 10)
	receipts, err := s.GetReceiptByBlock(blockID, number, nil, eLog)
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainCode=%v,err=%v,BlockNumber=%v", s.chain.BlockChainCode, err.Error(), number)
	} else {
		for _, v := range receipts {
			mp[v.TransactionHash] = v.Receipt
		}
	}

	txs := make([]*service.TxInterface, 0, len(list))
	addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	for _, tx := range list {
		// 补充字段
		hash := tx.Get("txID").String()
		fullTx := make(map[string]interface{}, 2)
		fullTx["tx"] = tx.String()
		fullTx["blockId"] = blockID
		if v, ok := mp[hash]; ok {
			fullTx["receipt"] = v
		}

		bs, _ := json.Marshal(fullTx)
		if s.CheckAddress(bs, addressList) {
			t := &service.TxInterface{TxHash: hash, Tx: fullTx}
			txs = append(txs, t)
		}
	}
	//r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}
	return r, txs
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry, flag bool) (*service.BlockInterface, []*service.TxInterface) {
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

	blockID := gjson.Parse(resp).Get("blockID").String()
	number := gjson.Parse(resp).Get("block_header.raw_data.number").String()

	if len(blockID) == 0 && len(number) == 0 {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,resp=%v", s.chain.BlockChainName, "resp is error", resp)
		return nil, nil
	}

	r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}

	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	//根据区块高度，获取交易log
	mp := make(map[string]interface{}, 10)
	receipts, err := s.GetReceiptByBlock(blockID, number, nil, eLog)
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainCode=%v,err=%v,BlockHash=%v", s.chain.BlockChainCode, err.Error(), blockID)
	} else {
		for _, v := range receipts {
			mp[v.TransactionHash] = v.Receipt
		}
	}

	list := gjson.Parse(resp).Get("transactions").Array()
	addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	txs := make([]*service.TxInterface, 0, len(list))
	for _, tx := range list {
		// 补充字段
		hash := tx.Get("txID").String()
		fullTx := make(map[string]interface{}, 2)
		fullTx["tx"] = tx.String()
		fullTx["blockId"] = blockID
		if v, ok := mp[hash]; ok {
			fullTx["receipt"] = v
		}

		bs, _ := json.Marshal(fullTx)
		if s.CheckAddress(bs, addressList) {
			t := &service.TxInterface{TxHash: hash, Tx: fullTx}
			txs = append(txs, t)
		}
	}
	//r := &service.BlockInterface{BlockHash: blockID, BlockNumber: number, Block: resp}
	return r, txs
}

func (s *Service) buildContract(receipt *service.TronReceipt) *service.TronReceipt {

	has := true

	//仅合约交易 有log
	for _, g := range receipt.Log {

		//trc20合约
		if len(g.Topics) == 3 && g.Topics[0] == service.TronTopic {
			mp := make(map[string]interface{}, 2)
			contractAddr := fmt.Sprintf("41%v", g.Address)
			token, err := s.getToken(int64(s.chain.BlockChainCode), contractAddr, contractAddr)
			if err != nil {
				has = false
				break
			}
			m := gjson.Parse(token).Map()
			if v, ok := m["decimals"]; ok {
				mp["contractDecimals"] = v.String()
			} else {
				has = false
				break
			}

			mp["data"] = g.Data
			bs, _ := json.Marshal(mp)
			g.Data = string(bs)
		}
	}

	if has {
		return receipt
	} else {
		//nodeId, _ := util.GetLocalNodeId()
		//task := service.NodeTask{Id: time.Now().UnixNano(), NodeId: nodeId, BlockChain: s.chain.BlockChainCode, TxHash: receipt.Id, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
		//_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), receipt.Id, task)
		return nil
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

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	return nil, nil
}

func getCoreAddr(addr string) string {
	if strings.HasPrefix(addr, "0x") {
		return strings.TrimLeft(addr, "0x") //去丢0x
	}

	if strings.HasPrefix(addr, "41") {
		return strings.TrimLeft(addr, "41") //去丢41
	}

	if strings.HasPrefix(addr, "0x41") {
		return strings.TrimLeft(addr, "0x41") //去丢41
	}
	return addr
}

func (s *Service) CheckAddress(txValue []byte, addrList []string) bool {
	if len(addrList) < 1 || len(txValue) < 1 {
		return false
	}
	txAddressList := make(map[string]int64, 10)

	root := gjson.ParseBytes(txValue)
	tx := root.Get("tx").String()
	txRoot := gjson.Parse(tx)
	contracts := txRoot.Get("raw_data.contract").Array()
	if len(contracts) < 1 {
		return false
	}
	r := contracts[0]
	txType := r.Get("type").String()

	var fromAddr, toAddr string
	var logs []gjson.Result
	var internalTransactions []gjson.Result

	fromAddr = r.Get("parameter.value.owner_address").String()
	txAddressList[getCoreAddr(fromAddr)] = 1

	//TransferContract,TransferAssetContract
	if r.Get("parameter.value.to_address").Exists() {
		toAddr = r.Get("parameter.value.to_address").String()
	}

	//DelegateResourceContract,UnDelegateResourceContract
	if r.Get("parameter.value.receiver_address").Exists() {
		toAddr = r.Get("parameter.value.receiver_address").String()
	}

	//TriggerSmartContract
	if r.Get("parameter.value.contract_address").Exists() {
		toAddr = r.Get("parameter.value.contract_address").String()
	}

	//AccountCreateContract
	if r.Get("parameter.value.account_address").Exists() {
		toAddr = r.Get("parameter.value.account_address").String()
	}

	txAddressList[getCoreAddr(toAddr)] = 1

	if txType == "TriggerSmartContract" {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		if receiptRoot.Get("receipt.result").String() != "SUCCESS" {
			return false
		}
		logs = receiptRoot.Get("log").Array()
		internalTransactions = receiptRoot.Get("internal_transactions").Array()

		//合约交易 合约调用下的TRC20
		if len(logs) > 0 {
			for _, v := range logs {
				topics := v.Get("topics").Array()
				//Transfer()
				if len(topics) >= 3 && topics[0].String() == service.TronTopic {
					from, _ := util.Hex2Address(topics[1].String())
					if len(from) > 0 {
						txAddressList[getCoreAddr(from)] = 1
					}
					to, _ := util.Hex2Address(topics[2].String())
					if len(to) > 0 {
						txAddressList[getCoreAddr(to)] = 1
					}
				}
			}
		}

		//合约调用下的内部交易TRX转帐和TRC10转账：
		if len(internalTransactions) > 0 {
			for _, v := range internalTransactions {
				fromAddr = v.Get("caller_address").String()
				toAddr = v.Get("transferTo_address").String()
				txAddressList[getCoreAddr(fromAddr)] = 1
				txAddressList[getCoreAddr(toAddr)] = 1
			}
		}
	}

	mp := make(map[string]int64, len(addrList))
	for _, v := range addrList {
		addr := getCoreAddr(v)
		mp[addr] = 1
	}

	has := false
	for k, _ := range txAddressList {
		//monitorAddr := getCoreAddr(v)
		if _, ok := mp[k]; ok {
			has = true
			break
		}
	}
	return has
}

func NewService(c *config.Chain, x *xlog.XLog, store service.StoreTaskInterface, nodeId string) service.BlockChainInterface {

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
		blockClient = chainService.NewTron(list, x)
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
		txClient = chainService.NewTron(list, x)
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
		receiptClient = chainService.NewTron(list, x)
	}

	return &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
		nodeId:             nodeId,
	}
}
