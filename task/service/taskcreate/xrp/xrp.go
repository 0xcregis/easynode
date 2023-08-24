package xrp

import (
	"errors"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/task"
	"github.com/0xcregis/easynode/task/config"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type XRP struct {
	log        *xlog.XLog
	api        blockchain.API
	blockChain int64
}

func (e *XRP) CreateNodeTask(nodeId string, blockChain int64, number string) (*task.NodeTask, error) {
	t := &task.NodeTask{
		NodeId:      nodeId,
		BlockNumber: number,
		BlockChain:  blockChain,
		TaskType:    2,
		TaskStatus:  0,
		CreateTime:  time.Now(),
		LogTime:     time.Now(),
		Id:          time.Now().UnixNano(),
	}
	return t, nil
}

func NewXRP(log *xlog.XLog, v *config.BlockConfig) *XRP {
	clusters := make([]*chainConfig.NodeCluster, 0, 2)
	for _, v := range v.Cluster {
		c := &chainConfig.NodeCluster{NodeUrl: v.NodeHost, NodeToken: v.NodeKey, Weight: v.Weight}
		clusters = append(clusters, c)
	}
	api := service.NewApi(v.BlockChainCode, clusters, log)
	return &XRP{
		blockChain: v.BlockChainCode,
		log:        log,
		api:        api,
	}
}

func (e *XRP) GetLatestBlockNumber() (int64, error) {
	log := e.log.WithFields(logrus.Fields{
		"id":         time.Now().UnixMilli(),
		"model":      "GetLatestBlockNumber",
		"blockChain": e.blockChain,
	})

	/**
	 {
	    "result": {
	        "ledger_hash": "586B03A3291D0E501035E861CAD220E5861500ADBAD8799C39E4F6657A83126C",
	        "ledger_index": 82011530,
	        "status": "success"
	    }
	}
	*/
	var lastNumber int64
	jsonResult, err := e.api.LatestBlock(e.blockChain)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	lastNumber = gjson.Parse(jsonResult).Get("result.ledger_index").Int()

	if lastNumber > 1 {
		//_ = s.UpdateLastNumber(v.BlockChainCode, lastNumber)
		return lastNumber, nil
	} else {
		return 0, errors.New("GetLastBlockNumber is fail")
	}
}
