package taskcreate

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/taskcreate/db"
	"github.com/uduncloud/easynode/task/service/taskcreate/ether"
	"github.com/uduncloud/easynode/task/service/taskcreate/tron"
	"time"
)

type Service struct {
	config    *config.Config
	dbService service.DbTaskCreateInterface
	log       *xlog.XLog
	nodeId    string
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
		_ = s.dbService.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) GetLastBlockNumberForTron(v *config.BlockConfig) error {
	lastNumber, err := tron.NewTron(s.log).GetLastBlockNumber(v)
	if err != nil {
		return err
	}
	if lastNumber > 1 {
		_ = s.dbService.UpdateLastNumber(v.BlockChainCode, lastNumber)
	}
	return nil
}

func (s *Service) NewBlockTask(v config.BlockConfig, log *logrus.Entry) error {
	if v.BlockMin < 1 {
		panic("Min blockNumber is not less 1")
	}

	//已经下发的最新区块高度
	UsedMaxNumber, lastBlockNumber, err := s.dbService.GetRecentNumber(v.BlockChainCode)
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
	list := make([]*service.NodeTask, 0)

	UsedMaxNumber++

	for UsedMaxNumber <= v.BlockMax {
		//ns := &service.NodeSource{BlockChain: v.BlockChainCode, BlockNumber: fmt.Sprintf("%v", UsedMaxNumber), SourceType: 2}

		task := &service.NodeTask{
			NodeId:      s.nodeId,
			BlockNumber: fmt.Sprintf("%v", UsedMaxNumber),
			BlockChain:  v.BlockChainCode,
			TaskType:    2,
			TaskStatus:  0,
		}

		list = append(list, task)
		UsedMaxNumber++
	}

	recentNumber := UsedMaxNumber - 1

	//生成任务并保存
	err = s.dbService.AddNodeTask(list)
	if err != nil {
		return err
	}
	//更新最新下发的区块高度
	err = s.dbService.UpdateRecentNumber(v.BlockChainCode, recentNumber)
	if err != nil {
		return err
	}

	return nil
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/create_task", 24*time.Hour)
	db := db.NewMySQLTaskCreateService(config, xg)
	nodeId, err := util.GetLocalNodeId()
	if err != nil {
		panic(err)
	}
	return &Service{
		config:    config,
		log:       xg,
		dbService: db,
		nodeId:    nodeId,
	}
}
