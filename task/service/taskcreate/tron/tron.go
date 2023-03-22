package tron

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/net/tron"
	"github.com/uduncloud/easynode/task/util"
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
	jsonResult, err := tron.Eth_GetBlockNumber(v.NodeHost, v.NodeKey)
	if err != nil {
		log.Errorf("Eth_GetBlockNumber|err=%v", err)
		return 0, err
	} else {
		log.Printf("Eth_GetBlockNumber|resp=%v", jsonResult)
	}

	//获取链的最新区块高度
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
