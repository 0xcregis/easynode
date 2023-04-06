package db

import (
	"encoding/json"
	"errors"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"time"
)

type TaskCreateFile struct {
	config *config.Config
	log    *xlog.XLog
}

func NewFileTaskCreateService(config *config.Config, xg *xlog.XLog) service.DbTaskCreateInterface {
	return &TaskCreateFile{
		config: config,
		log:    xg,
	}
}

func (t *TaskCreateFile) AddNodeTask(list []*service.NodeTask) error {
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
