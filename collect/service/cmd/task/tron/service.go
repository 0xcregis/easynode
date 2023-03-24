package tron

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	chainConfig "github.com/uduncloud/easynode/blockchain/config"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/net/tron"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/db"
	"strconv"
	"strings"
)

// todo 未全部改造完成，tron链最好使用http协议（tron节点 默认不支持json-rpc协议，需要手动开启）
// todo 当前版本 tron 是基于json-rpc 协议实现
// tron链 http 协议返回的数据 和 json-rpc 差异很大，需要特别处理

type Service struct {
	log                *xlog.XLog
	task               *db.Service
	chain              *config.Chain
	originChain        config.Chain
	txChainClient      chainService.API
	blockChainClient   chainService.API
	receiptChainClient chainService.API
}

func (s *Service) Monitor() {
}

func (s *Service) GetTx(txHash string, task *config.TxTask, eLog *logrus.Entry) *service.Tx {

	//调用接口
	//resp, err := s.txChainClient.GetTxByHash(int64(s.chain.BlockChainCode), txHash)
	resp, err := tron.Eth_GetTransactionByHash(s.txChainClient, txHash, s.log)
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
	return tx
}

func (s *Service) GetReceipt(txHash string, task *config.ReceiptTask, eLog *logrus.Entry) *service.Receipt {

	//调用接口
	resp, err := tron.Eth_GetTransactionReceiptByHash(s.receiptChainClient, txHash, s.log)
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
	return service.GetReceiptFromJson(resp)
}

func (s *Service) GetReceiptByBlock(blockHash, number string, task *config.ReceiptTask, eLog *logrus.Entry) []*service.Receipt {

	//调用接口
	var resp string
	var err error
	if len(number) > 0 {
		if !strings.HasPrefix(number, "0x") {
			n, _ := strconv.ParseInt(number, 10, 64)
			number = fmt.Sprintf("0x%x", n)
		}
		resp, err = tron.Eth_GetBlockReceiptByBlockNumber(s.receiptChainClient, number, s.log)
	} else if len(number) == 0 && len(blockHash) > 0 {
		resp, err = tron.Eth_GetBlockReceiptByBlockHash(s.receiptChainClient, blockHash, s.log)
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
	return service.GetReceiptListFromJson(resp)
}

func (s *Service) GetBlockByNumber(blockNumber string, task *config.BlockTask, eLog *logrus.Entry) (*service.Block, []*service.Tx) {
	if !strings.HasPrefix(blockNumber, "0x") {
		n, _ := strconv.ParseInt(blockNumber, 10, 64)
		blockNumber = fmt.Sprintf("0x%x", n)
	}

	//调用接口
	resp, err := tron.Eth_GetBlockByNumber(s.blockChainClient, blockNumber, s.log)
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
	return block, txList
}

func (s *Service) GetBlockByHash(blockHash string, cfg *config.BlockTask, eLog *logrus.Entry) (*service.Block, []*service.Tx) {
	//调用接口
	resp, err := tron.Eth_GetBlockByHash(s.blockChainClient, blockHash, s.log)
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
	return block, txList
}

func (s *Service) BalanceCluster(key string, clusterList []*config.FromCluster) (*config.FromCluster, error) {
	return nil, nil
}

func NewService(c *config.Chain, taskDb *config.TaskDb, sourceDb *config.SourceDb, x *xlog.XLog) service.BlockChainInterface {
	t := db.NewTaskService(taskDb, sourceDb, x)

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
		task:               t,
		chain:              c,
		originChain:        c.CopyChain(),
		txChainClient:      txClient,
		blockChainClient:   blockClient,
		receiptChainClient: receiptClient,
	}
}
