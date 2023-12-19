package filecoin

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	chainService "github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/common/ethtypes"
	"github.com/0xcregis/easynode/common/util"
	"github.com/filecoin-project/go-address"
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
	monitorAddress     map[string]int64
	lock               sync.RWMutex
}

func (s *Service) GetMultiBlockByNumber(blockNumber string, log *logrus.Entry, flag bool) ([]*collect.BlockInterface, []*collect.TxInterface) {

	start := time.Now()
	defer func() {
		log.Printf("GetMultiBlockByNumber.Duration =%v,blockNumber:%v", time.Since(start), blockNumber)
	}()

	//调用接口
	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), blockNumber, flag)
	if err != nil {
		log.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err.Error(), blockNumber)
		return nil, nil
	}

	//处理数据
	if resp == "" {
		log.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, "block is empty", blockNumber)
		return nil, nil
	}

	list := gjson.Parse(resp).Array()
	blockList := make([]*collect.BlockInterface, 0, len(list))

	//解析数据
	for _, v := range list {
		blockHash := v.Get("Cid").String()
		block := v.Get("BlockHead").String()
		mp := make(map[string]any, 2)
		mp["block"] = block
		mp["number"] = blockNumber
		mp["blockHash"] = blockHash
		r := &collect.BlockInterface{BlockHash: blockHash, BlockNumber: blockNumber, Block: mp}
		blockList = append(blockList, r)
	}

	if !flag { //仅区块数据，不涉及交易
		return blockList, nil
	}
	//
	////list := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, nil, eLog)
	//for _, v := range txList {
	//	receipt, err := s.GetReceipt(v.TxHash, log)
	//	if err != nil {
	//		log.Errorf("GetBlockByNumber|BlockChainCode=%v,err=%v,blockNumber=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockNumber, v.TxHash)
	//	} else {
	//		if receipt != nil && v.TxHash == receipt.TransactionHash {
	//			bs, _ := json.Marshal(receipt.Receipt)
	//			v.Receipt = string(bs)
	//		}
	//	}
	//}
	//
	//addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	//addressMp := rebuildAddress(addressList)
	//txs := make([]*collect.TxInterface, 0, len(txList))
	//for _, tx := range txList {
	//	bs, _ := json.Marshal(tx)
	//	if s.CheckAddress(bs, addressMp) {
	//		t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
	//		txs = append(txs, t)
	//	} else {
	//		log.Warnf("the tx is ignored,hash=%v", tx.TxHash)
	//	}
	//}
	//r := &service.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block}
	return nil, nil
}

func (s *Service) Monitor() {
	go func() {
		for {
			<-time.After(60 * time.Second)
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

			// reload monitor address
			s.reload()
		}
	}()
}

func (s *Service) reload() {
	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		s.log.Warnf("ReloadMonitorAddress BlockChainName=%v,err=%v", s.chain.BlockChainName, err)
	}

	s.lock.Lock()
	s.monitorAddress = rebuildAddress(addressList)
	s.lock.Unlock()
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

	root := gjson.Parse(resp)

	blockHash = root.Get("Cid").String()
	if len(blockHash) < 1 {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, resp, blockHash)
		return nil, nil
	}
	block := root.Get("BlockHead").String()
	blockNumber := gjson.Parse(block).Get("Height").String()

	r := &collect.BlockInterface{BlockHash: blockHash, BlockNumber: blockNumber, Block: block}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	txs := root.Get("BlockMessages").String()
	msgList := gjson.Parse(txs).Array()
	txList := make([]*collect.TxInterface, 0, len(msgList))
	for _, v := range msgList {
		hash := v.Get("CID./").String()
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
		receipt, err := s.GetReceipt(v.TxHash, eLog)
		if err != nil {
			task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: v.TxHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
			_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), v.TxHash, task)
			eLog.Warnf("GetBlockByHash|BlockChainCode=%v,err=%v,blockHash=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockHash, v.TxHash)
		} else {
			if receipt != nil {
				bs, _ := json.Marshal(receipt.Receipt)
				m["receipt"] = string(bs)
			}
		}

		v.Tx = m
	}

	//addressList, _ := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	//addressMp := rebuildAddress(addressList)
	resultTxs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx.Tx)
		if s.CheckAddress(bs, s.monitorAddress) {
			//t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx.Tx}
			resultTxs = append(resultTxs, tx)
		}
	}
	return r, resultTxs
}

func (s *Service) GetBlockByNumber(blockNumber string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {
	return nil, nil
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

	hash := root.Get("CID./").String()
	if len(hash) < 1 {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil
	}

	m := make(map[string]any, 2)
	m["tx"] = root.String()
	m["hash"] = hash

	receipt, err := s.GetReceipt(hash, eLog)
	if err != nil {
		task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: txHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
		_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), txHash, task)
		eLog.Errorf("GetTx|BlockChainCode=%v,err=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), hash)
	} else {

		if receipt != nil {
			if v, ok := receipt.Receipt.(*collect.Receipt); ok {
				resp, err := s.getBlockByNumber(v.BlockNumber)
				if err == nil {
					m["block"] = resp
					m["blockNumber"] = v.BlockNumber
				}
			}

			bs, _ := json.Marshal(receipt.Receipt)
			m["receipt"] = string(bs)
		}
	}
	//bs, _ := json.Marshal(m)
	r := &collect.TxInterface{TxHash: hash, Tx: m}
	return r
}

func (s *Service) getBlockByNumber(number string) (any, error) {

	resp, err := s.blockChainClient.GetBlockByNumber(int64(s.chain.BlockChainCode), number, false)

	if err != nil {
		return nil, err
	}

	//处理数据
	if resp == "" {
		return nil, errors.New("resp data is empty")
	}

	list := gjson.Parse(resp).Array()
	if len(list) < 1 {
		return nil, errors.New("resp data is empty")
	}

	//解析数据
	r := list[0].Get("BlockHead").String()
	return r, nil
}

func (s *Service) GetReceiptByBlock(blockHash, number string, eLog *logrus.Entry) ([]*collect.ReceiptInterface, error) {
	return nil, nil
}

func (s *Service) GetReceipt(txHash string, eLog *logrus.Entry) (*collect.ReceiptInterface, error) {

	//调用接口
	resp, err := s.receiptChainClient.GetTransactionReceiptByHash(int64(s.chain.BlockChainCode), txHash)
	//resp, err := ether.Eth_GetTransactionReceiptByHash(cluster.Host, cluster.Key, txHash, s.log)
	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	//处理数据
	if len(resp) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil, errors.New("receipt is null")
	}

	result := gjson.Parse(resp).Get("result").String()

	if len(result) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil, errors.New("receipt is null")
	}

	// 解析数据
	receipt := GetReceiptFromJson(result)

	bh, err := ethtypes.ParseEthHash(receipt.BlockHash)
	if err == nil {
		receipt.BlockHash = bh.ToCid().String()
	} else {
		eLog.Errorf("GetReceipt|BlockChainName=%v,parse.blockhash.err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	th, err := ethtypes.ParseEthHash(receipt.TransactionHash)
	if err == nil {
		receipt.TransactionHash = th.ToCid().String()
	} else {
		eLog.Errorf("GetReceipt|BlockChainName=%v,parse.TransactionHash.err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	fromAddr, err := ethtypes.ParseEthAddress(receipt.From)
	if err == nil {
		var addr address.Address
		addr, err = fromAddr.ToFilecoinAddress()
		if err == nil {
			receipt.From = addr.String()
		}
	}

	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,parse.from.err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	toAddr, err := ethtypes.ParseEthAddress(receipt.To)
	if err == nil {
		var addr address.Address
		addr, err = toAddr.ToFilecoinAddress()
		if err == nil {
			receipt.To = addr.String()
		}
	}
	if err != nil {
		eLog.Errorf("GetReceipt|BlockChainName=%v,parse.to.err=%v,txHash=%v", s.chain.BlockChainName, err.Error(), txHash)
		return nil, err
	}

	if receipt == nil || len(receipt.TransactionHash) < 1 {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil, errors.New("receipt is null")
	}

	receipt, err = s.buildContract(receipt)
	if err == nil {
		return &collect.ReceiptInterface{TransactionHash: receipt.TransactionHash, Receipt: receipt}, nil
	} else {
		//task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: txHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
		//_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), txHash, task)
		return nil, errors.New("receipt is null")
	}
}

func (s *Service) buildContract(receipt *collect.Receipt) (*collect.Receipt, error) {
	if receipt == nil {
		return nil, errors.New("receipt is null")
	}
	has := true

	// 仅有 合约交易，才能有logs
	for _, g := range receipt.Logs {

		//erc20
		if len(g.Topics) == 3 && g.Topics[0] == s.transferTopic {
			//处理 普通资产和 20 协议 资产转移
			mp := make(map[string]interface{}, 2)
			token, err := s.getToken(int64(s.chain.BlockChainCode), receipt.From, g.Address)
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
		return receipt, nil
	} else {
		return receipt, errors.New("can not get contract")
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

	if len(addrList) < 1 || len(tx) < 1 {
		return false
	}

	txAddressList := make(map[string]int64, 10)
	root := gjson.ParseBytes(tx)

	txBody := root.Get("tx").String()
	txRoot := gjson.Parse(txBody)

	fromAddr := txRoot.Get("From").String()
	txAddressList[getCoreAddr(fromAddr)] = 1

	toAddr := txRoot.Get("To").String()
	txAddressList[getCoreAddr(toAddr)] = 1

	if root.Get("receipt").Exists() {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		list := receiptRoot.Get("logs").Array()
		for _, v := range list {
			topics := v.Get("topics").Array()
			//Transfer()
			if len(topics) >= 3 && topics[0].String() == s.transferTopic {
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

	//mp := make(map[string]int64, len(addrList))
	//for _, v := range addrList {
	//	addr := getCoreAddr(v)
	//	mp[addr] = 1
	//}

	has := false
	s.lock.RLock()
	for k := range txAddressList {
		//monitorAddr := getCoreAddr(v)
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	s.lock.RUnlock()
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

	s := &Service{
		log:                x,
		chain:              c,
		store:              store,
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
		nodeId:             nodeId,
		transferTopic:      transferTopic,
		lock:               sync.RWMutex{},
	}

	s.reload()
	return s
}
