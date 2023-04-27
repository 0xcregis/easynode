package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"time"
)

type TaskCreateFile struct {
	config *config.Config
	log    *xlog.XLog
	sendCh chan []*kafka.Message
}

func NewFileTaskCreateService(config *config.Config, sendCh chan []*kafka.Message, xg *xlog.XLog) service.StoreTaskInterface {
	return &TaskCreateFile{
		config: config,
		log:    xg,
		sendCh: sendCh,
	}
}

func (t *TaskCreateFile) AddNodeTask(list []*service.NodeTask) error {
	resultList := make([]*kafka.Message, 0)
	for _, v := range list {
		bs, _ := json.Marshal(v)
		msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: 0, Key: []byte(v.NodeId), Value: bs}
		resultList = append(resultList, msg)
	}
	t.sendCh <- resultList
	return nil
}

func (t *TaskCreateFile) UpdateLastNumber(blockChainCode int64, latestNumber int64) error {

	mp := make(map[int64]*service.BlockNumber, 2)

	bs, err := util.ReadLatestBlock()
	if err != nil {
		bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode, LogTime: time.Now()}
		mp[blockChainCode] = &bn
		//return err
	} else {
		_ = json.Unmarshal(bs, &mp)
		if _, ok := mp[blockChainCode]; ok {
			mp[blockChainCode].LatestNumber = latestNumber
			mp[blockChainCode].LogTime = time.Now()
		} else {
			bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode, LogTime: time.Now()}
			mp[blockChainCode] = &bn
		}
	}

	bs, _ = json.Marshal(mp)
	err = util.WriteLatestBlock(string(bs))
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateFile) UpdateRecentNumber(blockChainCode int64, recentNumber int64) error {
	mp := make(map[int64]*service.BlockNumber, 2)
	bs, err := util.ReadLatestBlock()
	if err != nil {
		bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode, LogTime: time.Now()}
		mp[blockChainCode] = &bn
	} else {
		_ = json.Unmarshal(bs, &mp)
		if _, ok := mp[blockChainCode]; ok {
			mp[blockChainCode].RecentNumber = recentNumber
			mp[blockChainCode].LogTime = time.Now()
		} else {
			bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode, LogTime: time.Now()}
			mp[blockChainCode] = &bn
		}
	}

	bs, _ = json.Marshal(mp)
	err = util.WriteLatestBlock(string(bs))
	if err != nil {
		return err
	}

	return nil
}

func (t *TaskCreateFile) GetRecentNumber(blockCode int64) (int64, int64, error) {

	bs, err := util.ReadLatestBlock()
	if err != nil {
		return 0, 0, err
	}

	mp := make(map[int64]*service.BlockNumber, 2)
	_ = json.Unmarshal(bs, &mp)

	if v, ok := mp[blockCode]; ok {
		return v.RecentNumber, v.LatestNumber, nil
	} else {
		return 0, 0, errors.New("no record")
	}

}
