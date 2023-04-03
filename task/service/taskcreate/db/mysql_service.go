package db

import (
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/common/sql"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type TaskCreateDb struct {
	config        *config.Config
	taskDb        *gorm.DB
	blockNumberDb *gorm.DB
	log           *xlog.XLog
}

func NewMySQLTaskCreateService(config *config.Config, xg *xlog.XLog) service.DbTaskCreateInterface {
	task, err := sql.Open(config.NodeTaskDb.User, config.NodeTaskDb.Password, config.NodeTaskDb.Addr, config.NodeTaskDb.DbName, config.NodeTaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	blockNumber, err := sql.Open(config.BlockNumberDb.User, config.BlockNumberDb.Password, config.BlockNumberDb.Addr, config.BlockNumberDb.DbName, config.BlockNumberDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &TaskCreateDb{
		config:        config,
		taskDb:        task,
		blockNumberDb: blockNumber,
		log:           xg,
	}
}

func (t *TaskCreateDb) getNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", t.config.NodeTaskDb.Table, time.Now().Format(service.DayFormat))
	return table
}

func (t *TaskCreateDb) AddNodeTask(list []*service.NodeTask) error {
	err := t.taskDb.Table(t.getNodeTaskTable()).Clauses(clause.Insert{Modifier: "IGNORE"}).Omit("id,log_time,create_time").CreateInBatches(&list, 10).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateDb) UpdateLastNumber(blockChainCode int64, latestNumber int64) error {
	bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode}
	err := t.blockNumberDb.Table(t.config.BlockNumberDb.Table).Omit("id,create_time,log_time,recent_number").Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"latest_number"})}).Create(&bn).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateDb) UpdateRecentNumber(blockChainCode int64, recentNumber int64) error {
	bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode}
	err := t.blockNumberDb.Table(t.config.BlockNumberDb.Table).Omit("id,create_time,log_time,latest_number").Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"recent_number"})}).Create(&bn).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *TaskCreateDb) GetRecentNumber(blockCode int64) (int64, int64, error) {
	var Num int64
	err := t.blockNumberDb.Table(t.config.BlockNumberDb.Table).Where("chain_code=?", blockCode).Count(&Num).Error
	if err != nil {
		return 0, 0, err
	}

	if Num < 1 {
		return 0, 0, nil
	}

	var temp service.BlockNumber
	err = t.blockNumberDb.Table(t.config.BlockNumberDb.Table).Where("chain_code=?", blockCode).First(&temp).Error
	if err != nil {
		return 0, 0, errors.New("no record")
	}
	return temp.RecentNumber, temp.LatestNumber, nil
}
