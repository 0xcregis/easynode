package tron

import (
	"errors"
	"strconv"
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

type Tron struct {
	log        *xlog.XLog
	api        blockchain.API
	blockChain int64
}

func (e *Tron) CreateNodeTask(nodeId string, blockChain int64, number string) (*task.NodeTask, error) {
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

func NewTron(log *xlog.XLog, v *config.BlockConfig) *Tron {
	clusters := make([]*chainConfig.NodeCluster, 0, 2)
	for _, v := range v.Cluster {
		c := &chainConfig.NodeCluster{NodeUrl: v.NodeHost, NodeToken: v.NodeKey, Weight: v.Weight}
		clusters = append(clusters, c)
	}
	api := service.NewApi(v.BlockChainCode, clusters, log)
	return &Tron{
		log:        log,
		api:        api,
		blockChain: v.BlockChainCode,
	}
}

func (e *Tron) GetLatestBlockNumber() (int64, error) {
	log := e.log.WithFields(logrus.Fields{
		"id":         time.Now().UnixMilli(),
		"model":      "GetLastBlockNumber",
		"blockChain": e.blockChain,
	})
	var lastNumber int64
	/**
	  {\"blockId\":\"0000000002f52f21275a4e244b191f29dc289bfc66bce08a18c3d5051fcf7203\",\"number\":49622817}
	*/
	jsonResult, err := e.api.LatestBlock(e.blockChain)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	number := gjson.Parse(jsonResult).Get("number").String()
	lastNumber, err = strconv.ParseInt(number, 0, 64)
	//lastNumber, err = util.HexToInt(number)
	if err != nil {
		log.Errorf("HexToInt|err=%v", err)
		return 0, err
	}
	if lastNumber > 1 {
		//_ = s.UpdateLastNumber(v.BlockChainCode, lastNumber)
		return lastNumber, nil
	} else {
		return 0, errors.New("GetLastBlockNumber is fail")
	}
}
