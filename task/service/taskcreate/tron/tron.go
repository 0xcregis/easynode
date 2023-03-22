package tron

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	blockChainConfig "github.com/uduncloud/easynode/blockchain/config"
	"github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/task/config"
	"strconv"
	"time"
)

type Tron struct {
	log *xlog.XLog
}

func NewTron(log *xlog.XLog) *Tron {
	return &Tron{
		log: log,
	}
}

func (e *Tron) GetLastBlockNumber(v *config.BlockConfig) (int64, error) {
	log := e.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "GetLastBlockNumber",
	})
	var lastNumber int64
	clusters := map[int64][]*blockChainConfig.NodeCluster{205: {{NodeUrl: v.NodeHost, NodeToken: v.NodeKey}}}
	/**
	{
	    "code": 0,
	    "data": "{\"blockId\":\"0000000002f52f21275a4e244b191f29dc289bfc66bce08a18c3d5051fcf7203\",\"number\":49622817}"
	}

	*/
	jsonResult, err := service.NewTron(clusters, e.log).LatestBlock(205)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	//获取链的最新区块高度
	code := gjson.Parse(jsonResult).Get("code").Int()
	data := gjson.Parse(jsonResult).Get("data").String()
	if code != 0 {
		log.Errorf("Eth_GetBlockNumber|err=%v", data)
		return 0, errors.New(data)
	}

	number := gjson.Parse(data).Get("number").String()
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
