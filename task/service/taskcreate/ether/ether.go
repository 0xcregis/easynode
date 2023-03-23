package ether

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	blockChainConfig "github.com/uduncloud/easynode/blockchain/config"
	"github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/util"
	"time"
)

type Ether struct {
	log *xlog.XLog
}

func NewEther(log *xlog.XLog) *Ether {
	return &Ether{
		log: log,
	}
}

func (e *Ether) GetLastBlockNumber(v *config.BlockConfig) (int64, error) {
	log := e.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "GetLastBlockNumber",
	})

	clusters := map[int64][]*blockChainConfig.NodeCluster{200: {{NodeUrl: v.NodeHost, NodeToken: v.NodeKey}}}
	/**
	  {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x1019b4c\"}
	*/
	var lastNumber int64
	jsonResult, err := service.NewEth(clusters, e.log).LatestBlock(200)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	number := gjson.Parse(jsonResult).Get("result").String()
	lastNumber, err = util.HexToInt(number)
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
