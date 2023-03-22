package taskcreate

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/common/sql"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/taskcreate/ether"
	"github.com/uduncloud/easynode/task/service/taskcreate/tron"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type Service struct {
	config       *config.Config
	nodeSourceDb *gorm.DB
	taskDb       *gorm.DB
	nodeInfoDb   *gorm.DB
	log          *xlog.XLog
}

func (s *Service) Start() {

	go func() {
		for true {
			//处理自产生区块任务，包括：区块
			s.CreateBlockTask()
			<-time.After(10 * time.Second)
		}
	}()
}

func (s *Service) CreateBlockTask() {
	log := s.log.WithFields(logrus.Fields{
		"model": "CreateBlockTask",
		"id":    time.Now().UnixMilli(),
	})
	blockConfigs := s.config.BlockConfigs
	if len(blockConfigs) < 1 {
		log.Warnf("config.BlockConfigs|info=%v", "chain config is null")
		return
	}

	for _, v := range blockConfigs {
		log.Infof("blockConfig=%+v", v)
		//获取公链最新区块高度
		if v.BlockChainCode == 200 {
			//ether
			err := s.GetLastBlockNumberForEther(v)
			if err != nil {
				log.Errorf("GetLastBlockNumberForEther|err=%v", err)
			}
		} else if v.BlockChainCode == 205 {
			//tron
			err := s.GetLastBlockNumberForTron(v)
			if err != nil {
				log.Errorf("GetLastBlockNumberForTron|err=%v", err)
			}
		}

		//生成最新的区块任务
		err := s.NewBlockTask(*v, log)
		if err != nil {
			log.Errorf("NewBlockTask|err=%v", err.Error())
		}

	}
}

func (s *Service) GetLastBlockNumberForEther(v *config.BlockConfig) error {
	lastNumber, err := ether.NewEther(s.log).GetLastBlockNumber(v)
	if err != nil {
		return err
	}
	if lastNumber > 1 {
		_ = s.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) GetLastBlockNumberForTron(v *config.BlockConfig) error {
	lastNumber, err := tron.NewTron(s.log).GetLastBlockNumber(v)
	if err != nil {
		return err
	}
	if lastNumber > 1 {
		_ = s.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) NewBlockTask(v config.BlockConfig, log *logrus.Entry) error {
	if v.BlockMin < 1 {
		panic("Min blockNumber is not less 1")
	}

	//已经下发的最新区块高度
	UsedMaxNumber, lastBlockNumber, err := s.GetRecentNumber(v.BlockChainCode)
	if err != nil {
		log.Errorf("GetRecentNumber|err=%v", err)
		return err
	}

	//如果从未下发该链区块任务，则 使用配置的最小区块高度
	if UsedMaxNumber == 0 {
		UsedMaxNumber = v.BlockMin
	}

	//获取指定区块高度
	//UsedMaxNumber~BlockMax

	//如果没有配置最大高度，则最大高度 时时读取链上最新高度
	if v.BlockMax < 1 {
		v.BlockMax = lastBlockNumber
	}

	if UsedMaxNumber >= v.BlockMax {
		return errors.New("UsedMaxNumber > BlockMax")
	}
	list := make([]*service.NodeSource, 0)

	UsedMaxNumber++

	for UsedMaxNumber <= v.BlockMax {
		ns := &service.NodeSource{BlockChain: v.BlockChainCode, BlockNumber: fmt.Sprintf("%v", UsedMaxNumber), SourceType: 2}
		list = append(list, ns)
		UsedMaxNumber++
	}

	recentNumber := UsedMaxNumber - 1

	//生成任务并保存
	err = s.AddNodeSource(list)
	if err != nil {
		return err
	}
	//更新最新下发的区块高度
	err = s.UpdateRecentNumber(v.BlockChainCode, recentNumber)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) AddNodeSource(list []*service.NodeSource) error {
	err := s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Clauses(clause.Insert{Modifier: "IGNORE"}).Omit("id,create_time").CreateInBatches(&list, 10).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateLastNumber(blockChainCode int64, latestNumber int64) error {
	bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode}
	err := s.nodeSourceDb.Table(service.BLOCK_NUMBER_TABLE).Omit("id,create_time,log_time,recent_number").Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"latest_number"})}).Create(&bn).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateRecentNumber(blockChainCode int64, recentNumber int64) error {
	bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode}
	err := s.nodeSourceDb.Table(service.BLOCK_NUMBER_TABLE).Omit("id,create_time,log_time,latest_number").Clauses(clause.OnConflict{DoUpdates: clause.AssignmentColumns([]string{"recent_number"})}).Create(&bn).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetRecentNumber(blockCode int64) (int64, int64, error) {
	var Num int64
	err := s.nodeSourceDb.Table(service.BLOCK_NUMBER_TABLE).Where("chain_code=?", blockCode).Count(&Num).Error
	if err != nil {
		return 0, 0, err
	}

	if Num < 1 {
		return 0, 0, nil
	}

	var temp service.BlockNumber
	err = s.nodeSourceDb.Table(service.BLOCK_NUMBER_TABLE).Where("chain_code=?", blockCode).First(&temp).Error
	if err != nil {
		return 0, 0, errors.New("no record")
	}
	return temp.RecentNumber, temp.LatestNumber, nil
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/create_task", 24*time.Hour)
	s, err := sql.Open(config.NodeSourceDb.User, config.NodeSourceDb.Password, config.NodeSourceDb.Addr, config.NodeSourceDb.DbName, config.NodeSourceDb.Port, xg)
	if err != nil {
		panic(err)
	}

	info, err := sql.Open(config.NodeInfoDb.User, config.NodeInfoDb.Password, config.NodeInfoDb.Addr, config.NodeInfoDb.DbName, config.NodeInfoDb.Port, xg)
	if err != nil {
		panic(err)
	}

	task, err := sql.Open(config.NodeTaskDb.User, config.NodeTaskDb.Password, config.NodeTaskDb.Addr, config.NodeTaskDb.DbName, config.NodeTaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &Service{
		config:       config,
		nodeSourceDb: s,
		nodeInfoDb:   info,
		taskDb:       task,
		log:          xg,
	}
}
