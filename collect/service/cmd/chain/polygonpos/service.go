package polygonpos

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	chainService "github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/common/util"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type Service struct {
	log                    *xlog.XLog
	chain                  *config.Chain
	nodeId                 string
	transferTopic          string
	nftTransferSingleTopic string
	//nftTransferBatchTopic  string
	store              collect.StoreTaskInterface
	txChainClient      blockchain.API
	blockChainClient   blockchain.API
	receiptChainClient blockchain.API
}

func (s *Service) GetMultiBlockByNumber(blockNumber string, log *logrus.Entry, flag bool) ([]*collect.BlockInterface, []*collect.TxInterface) {
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

	result := gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := GetBlockFromJson(result)

	if len(block.BlockHash) < 1 {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, resp, blockHash)
		return nil, nil
	}

	r := &collect.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block, BlockTime: block.BlockTime}
	if !flag { //仅区块，不涉及交易
		return r, nil
	}

	m := make(map[string]*collect.ReceiptInterface, 10)

	list, err := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, eLog)
	if err == nil && len(list) > 0 {
		for _, v := range list {
			m[v.TransactionHash] = v
		}
	} else {
		for _, v := range txList {
			receipt, err := s.GetReceipt(v.TxHash, eLog)
			if err != nil {
				//retry on monitor
				task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: v.TxHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
				_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), v.TxHash, task)
				eLog.Warnf("GetBlockByHash|BlockChainCode=%v,err=%v,blockHash=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockHash, v.TxHash)
			} else {
				m[receipt.TransactionHash] = receipt
			}
		}
	}

	for _, v := range txList {
		if receipt, ok := m[v.TxHash]; ok {
			bs, _ := json.Marshal(receipt.Receipt)
			v.Receipt = string(bs)
		}
	}

	//get monitor address list
	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		eLog.Errorf("GetBlockByHash|BlockChainName=%v,err=%v,blockHash=%v", s.chain.BlockChainName, err, blockHash)
	}
	if len(addressList) < 1 {
		eLog.Warnf("the tx is ignored due to monitor address is empty,blockHash=%v", blockHash)
		return r, nil
	}

	//filter tx
	addressMp := rebuildAddress(addressList)
	txs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx)
		if s.CheckAddress(bs, addressMp) {
			t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
			txs = append(txs, t)
		} else {
			eLog.Warnf("the tx is ignored,hash=%v", tx.TxHash)
		}
	}
	return r, txs
}

func (s *Service) GetBlockByNumber(blockNumber string, eLog *logrus.Entry, flag bool) (*collect.BlockInterface, []*collect.TxInterface) {

	start := time.Now()
	defer func() {
		eLog.Printf("GetBlockByNumber.Duration =%v,blockNumber:%v", time.Since(start), blockNumber)
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

	result := gjson.Parse(resp).Get("result").String()

	//解析数据
	block, txList := GetBlockFromJson(result)

	if len(block.BlockHash) < 1 {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, resp, blockNumber)
		return nil, nil
	}

	r := &collect.BlockInterface{BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, Block: block, BlockTime: block.BlockTime}
	if !flag { //仅区块数据，不涉及交易
		return r, nil
	}

	m := make(map[string]*collect.ReceiptInterface, 10)
	list, err := s.GetReceiptByBlock(block.BlockHash, block.BlockNumber, eLog)
	if err == nil && len(list) > 0 {
		for _, v := range list {
			m[v.TransactionHash] = v
		}
	} else {
		for _, v := range txList {
			receipt, err := s.GetReceipt(v.TxHash, eLog)
			if err != nil {
				//retry on monitor
				task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: v.TxHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
				_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), v.TxHash, task)
				eLog.Warnf("GetBlockByNumber|BlockChainCode=%v,err=%v,blockNumber=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), blockNumber, v.TxHash)
			} else {
				m[receipt.TransactionHash] = receipt
			}
		}
	}

	for _, v := range txList {
		if receipt, ok := m[v.TxHash]; ok {
			bs, _ := json.Marshal(receipt.Receipt)
			v.Receipt = string(bs)
		}
	}

	//get monitor address list
	addressList, err := s.store.GetMonitorAddress(int64(s.chain.BlockChainCode))
	if err != nil {
		eLog.Errorf("GetBlockByNumber|BlockChainName=%v,err=%v,blockNumber=%v", s.chain.BlockChainName, err, blockNumber)
	}
	if len(addressList) < 1 {
		eLog.Warnf("the tx is ignored due to monitor address is empty,blockNumber=%v", blockNumber)
		return r, nil
	}

	//filter tx
	addressMp := rebuildAddress(addressList)
	txs := make([]*collect.TxInterface, 0, len(txList))
	for _, tx := range txList {
		bs, _ := json.Marshal(tx)
		if s.CheckAddress(bs, addressMp) {
			t := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
			txs = append(txs, t)
		} else {
			eLog.Warnf("the tx is ignored,hash=%v", tx.TxHash)
		}
	}
	return r, txs
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
	result := gjson.Parse(resp).Get("result").String()

	//解析数据
	tx := GetTxFromJson(result)

	if len(tx.TxHash) < 1 {
		eLog.Errorf("GetTx|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, resp, txHash)
		return nil
	}

	//fix tx.time
	if len(tx.TxTime) < 1 {
		b, _ := s.GetBlockByNumber(tx.BlockNumber, eLog, false)
		if b != nil {
			tx.TxTime = b.BlockTime
		}
	}

	receipt, err := s.GetReceipt(tx.TxHash, eLog)
	if err != nil {
		//retry on monitor
		task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: tx.TxHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
		_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), tx.TxHash, task)
		eLog.Warnf("GetTx|BlockChainCode=%v,err=%v,txHash=%v", s.chain.BlockChainCode, err.Error(), tx.TxHash)
	} else {
		if receipt != nil {
			bs, _ := json.Marshal(receipt.Receipt)
			tx.Receipt = string(bs)
		}
	}

	r := &collect.TxInterface{TxHash: tx.TxHash, Tx: tx}
	return r
}

func (s *Service) GetReceiptByBlock(blockHash, number string, eLog *logrus.Entry) ([]*collect.ReceiptInterface, error) {

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
		return nil, err
	}

	//处理数据
	if resp == "" {
		eLog.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		return nil, errors.New("receipt is null")
	}

	resp = gjson.Parse(resp).Get("result").String()
	// 解析数据
	receiptList := GetReceiptListFromJson(resp)

	if len(receiptList) == 0 {
		eLog.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,blocknumber=%v, blockHash=%v", s.chain.BlockChainName, "receipts is null", number, blockHash)
		return nil, errors.New("receipt is null")
	}

	rs := make([]*collect.ReceiptInterface, 0, len(receiptList))
	for _, v := range receiptList {
		txHash := v.TransactionHash
		v, err = s.buildContract(v)
		if err == nil {
			r := &collect.ReceiptInterface{TransactionHash: txHash, Receipt: v}
			rs = append(rs, r)
		} else {
			//收据数据异常，则加入重试机制
			//nodeId, _ := util.GetLocalNodeId()
			task := collect.NodeTask{Id: time.Now().UnixNano(), BlockChain: s.chain.BlockChainCode, NodeId: s.nodeId, TxHash: txHash, TaskType: 1, TaskStatus: 0, CreateTime: time.Now(), LogTime: time.Now()}
			_ = s.store.StoreErrTxNodeTask(int64(s.chain.BlockChainCode), txHash, task)
		}
	}
	return rs, nil
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
	if resp == "" {
		eLog.Errorf("GetReceipt|BlockChainName=%v,err=%v,txHash=%v", s.chain.BlockChainName, "receipt is empty", txHash)
		return nil, errors.New("receipt is null")
	}

	resp = gjson.Parse(resp).Get("result").String()

	// 解析数据
	receipt := GetReceiptFromJson(resp)

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
		//save contract of erc-20 to receipt.log
		if len(g.Topics) == 3 && g.Topics[0] == s.transferTopic {
			//处理 普通资产和 20 协议 资产转移
			mp := make(map[string]interface{}, 2)
			token, err := s.getToken(int64(s.chain.BlockChainCode), receipt.From, g.Address, "20")
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
			mp["token"] = token
			mp["eip"] = 20
			mp["data"] = g.Data
			bs, _ := json.Marshal(mp)
			g.Data = string(bs)
		}

		//erc721
		//save contract of erc-721 to receipt.log
		if len(g.Topics) == 4 && g.Topics[0] == s.transferTopic {
			//处理 普通资产和 721 协议 资产转移
			mp := make(map[string]interface{}, 2)
			token, _ := s.getToken(int64(s.chain.BlockChainCode), receipt.From, g.Address, "721")
			mp["token"] = token
			mp["eip"] = 721
			mp["data"] = g.Data
			bs, _ := json.Marshal(mp)
			g.Data = string(bs)
		}

		//erc1155
		// save contract of erc-1155 to receipt.log
		if len(g.Topics) == 4 && g.Topics[0] == s.nftTransferSingleTopic {
			//处理 普通资产和 1155 协议 资产转移
			mp := make(map[string]interface{}, 3)
			token, _ := s.getToken(int64(s.chain.BlockChainCode), receipt.From, g.Address, "1155")
			mp["token"] = token
			mp["eip"] = 1155
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

func (s *Service) getToken(blockChain int64, from string, contract string, eip string) (string, error) {
	token, err := s.store.GetContract(blockChain, contract)
	if err == nil && len(token) > 5 {
		return token, nil
	}

	token, err = s.txChainClient.Token(blockChain, contract, "", eip)
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

	fromAddr := root.Get("from").String()
	txAddressList[getCoreAddr(fromAddr)] = 1

	toAddr := root.Get("to").String()
	txAddressList[getCoreAddr(toAddr)] = 1

	if root.Get("receipt").Exists() {
		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		list := receiptRoot.Get("logs").Array()
		for _, v := range list {
			topics := v.Get("topics").Array()
			//Transfer(),erc20
			if len(topics) == 3 && topics[0].String() == s.transferTopic {
				from, _ := util.Hex2Address(topics[1].String())
				if len(from) > 0 {
					txAddressList[getCoreAddr(from)] = 1
				}
				to, _ := util.Hex2Address(topics[2].String())
				if len(to) > 0 {
					txAddressList[getCoreAddr(to)] = 1
				}
			}

			//Transfer(),erc721
			if len(topics) == 4 && topics[0].String() == s.transferTopic {
				from, _ := util.Hex2Address(topics[1].String())
				if len(from) > 0 {
					txAddressList[getCoreAddr(from)] = 1
				}
				to, _ := util.Hex2Address(topics[2].String())
				if len(to) > 0 {
					txAddressList[getCoreAddr(to)] = 1
				}
			}

			//Transfer(),erc1155
			if len(topics) == 4 && topics[0].String() == s.nftTransferSingleTopic {
				operator, _ := util.Hex2Address(topics[1].String())
				if len(operator) > 0 {
					txAddressList[getCoreAddr(operator)] = 1
				}

				from, _ := util.Hex2Address(topics[2].String())
				if len(from) > 0 {
					txAddressList[getCoreAddr(from)] = 1
				}

				to, _ := util.Hex2Address(topics[3].String())
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
	for k := range txAddressList {
		//monitorAddr := getCoreAddr(v)
		if _, ok := addrList[k]; ok {
			has = true
			break
		}
	}
	return has
}

func NewService(c *config.Chain, x *xlog.XLog, store collect.StoreTaskInterface, nodeId string, transferTopic, nftTransferSingleTopic string) collect.BlockChainInterface {

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
		log:                    x,
		chain:                  c,
		store:                  store,
		txChainClient:          txClient,
		blockChainClient:       blockClient,
		receiptChainClient:     receiptClient,
		nodeId:                 nodeId,
		transferTopic:          transferTopic,
		nftTransferSingleTopic: nftTransferSingleTopic,
	}
}
