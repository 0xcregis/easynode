package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/common/sql"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TaskCreateSqlDb struct {
	config *config.Config
	taskDb *gorm.DB
	log    *xlog.XLog
}

func NewMySQLTaskCreateService(config *config.Config, xg *xlog.XLog) service.StoreTaskInterface {
	task, err := sql.Open(config.NodeTaskDb.User, config.NodeTaskDb.Password, config.NodeTaskDb.Addr, config.NodeTaskDb.DbName, config.NodeTaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &TaskCreateSqlDb{
		config: config,
		taskDb: task,
		log:    xg,
	}
}

func (t *TaskCreateSqlDb) getNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", t.config.NodeTaskDb.Table, time.Now().Format(service.DayFormat))
	return table
}

func (t *TaskCreateSqlDb) AddNodeTask(list []*service.NodeTask) error {
	err := t.taskDb.Table(t.getNodeTaskTable()).Clauses(clause.Insert{Modifier: "IGNORE"}).Omit("id,log_time,create_time").CreateInBatches(&list, 10).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateSqlDb) UpdateLastNumber(blockChainCode int64, latestNumber int64) error {

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

func (t *TaskCreateSqlDb) UpdateRecentNumber(blockChainCode int64, recentNumber int64) error {
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

func (t *TaskCreateSqlDb) GetRecentNumber(blockCode int64) (int64, int64, error) {

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
