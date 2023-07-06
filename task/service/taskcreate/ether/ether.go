package ether

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	chainConfig "github.com/uduncloud/easynode/blockchain/config"
	"github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/config"
	"time"
)

type Ether struct {
	log        *xlog.XLog
	api        service.API
	blockChain int64
}

func NewEther(log *xlog.XLog, v *config.BlockConfig) *Ether {
	clusters := make([]*chainConfig.NodeCluster, 0, 2)

	for _, v := range v.Cluster {
		c := &chainConfig.NodeCluster{NodeUrl: v.NodeHost, NodeToken: v.NodeKey, Weight: v.Weight}
		clusters = append(clusters, c)
	}

	api := service.NewEth(clusters, log)
	return &Ether{
		blockChain: v.BlockChainCode,
		log:        log,
		api:        api,
	}
}

func (e *Ether) GetLatestBlockNumber() (int64, error) {
	log := e.log.WithFields(logrus.Fields{
		"id":         time.Now().UnixMilli(),
		"model":      "GetLatestBlockNumber",
		"blockChain": e.blockChain,
	})

	/**
	  {\"jsonrpc\":\"2.0\",\"id\":1,\"result\":\"0x1019b4c\"}
	*/
	var lastNumber int64
	jsonResult, err := e.api.LatestBlock(e.blockChain)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	number := gjson.Parse(jsonResult).Get("result").String()
	lastNumber, err = util.HexToInt2(number)
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
