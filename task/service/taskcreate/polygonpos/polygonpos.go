package polygonpos

import (
	"errors"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/common/util"
	"github.com/0xcregis/easynode/task/config"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type PolygonPos struct {
	log        *xlog.XLog
	api        blockchain.API
	blockChain int64
}

func NewPolygonPos(log *xlog.XLog, v *config.BlockConfig) *PolygonPos {
	clusters := make([]*chainConfig.NodeCluster, 0, 2)

	for _, v := range v.Cluster {
		c := &chainConfig.NodeCluster{NodeUrl: v.NodeHost, NodeToken: v.NodeKey, Weight: v.Weight}
		clusters = append(clusters, c)
	}

	api := service.NewApi(v.BlockChainCode, clusters, log)
	return &PolygonPos{
		blockChain: v.BlockChainCode,
		log:        log,
		api:        api,
	}
}

func (e *PolygonPos) GetLatestBlockNumber() (int64, error) {
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
