package taskcreate

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/taskcreate/db"
	"math/rand"
	"time"
)

type Service struct {
	config      *config.Config
	store       service.StoreTaskInterface
	log         *xlog.XLog
	kafkaClient *kafkaClient.EasyKafka
	sendCh      chan []*kafka.Message
	api         map[int64]service.BlockChainInterface
}

func (s *Service) Start(ctx context.Context) {

	log := s.log.WithFields(logrus.Fields{
		"model": "CreateBlockProc",
		"id":    time.Now().UnixMilli(),
	})
	blockConfigs := s.config.BlockConfigs
	if len(blockConfigs) < 1 {
		log.Warnf("config.BlockConfigs|info=%v", "chain config is null")
		return
	}

	go s.startKafka(ctx)

	for _, v := range blockConfigs {
		notify := make(chan struct{})
		go func(ctx context.Context, cfg *config.BlockConfig, log *logrus.Entry, notify chan struct{}) {
			s.updateLatestBlock(ctx, cfg, log, notify)
		}(ctx, v, log, notify)

		go func(ctx context.Context, cfg *config.BlockConfig, log *logrus.Entry, notify chan struct{}) {
			s.startCreateBlockProc(ctx, cfg, log, notify)
		}(ctx, v, log, notify)
	}

}

func (s *Service) startKafka(ctx context.Context) {
	broker := fmt.Sprintf("%v:%v", s.config.TaskKafka.Host, s.config.TaskKafka.Port)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	s.kafkaClient.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, s.sendCh, nil, ctx, 2)
}

func (s *Service) updateLatestBlock(ctx context.Context, cfg *config.BlockConfig, log *logrus.Entry, notify chan struct{}) {
	interrupt := true
	for interrupt {
		var lastNumber int64
		var err error
		if api, ok := s.api[cfg.BlockChainCode]; ok {
			lastNumber, err = api.GetLatestBlockNumber()
		}
		if err != nil {
			log.Errorf("GetLastBlockNumber|err=%v", err.Error())
			time.Sleep(3 * time.Second)
			continue
		}
		if lastNumber > 1 {
			err = s.store.UpdateLastNumber(cfg.BlockChainCode, lastNumber)
			if err == nil { //通知分配区块任务
				notify <- struct{}{}
			}
		}

		select {
		case <-ctx.Done():
			interrupt = false
			break
		default:
			time.Sleep(10 * time.Second)
			continue
		}
	}
}
func (s *Service) startCreateBlockProc(ctx context.Context, cfg *config.BlockConfig, log *logrus.Entry, notify chan struct{}) {
	interrupt := true
	for interrupt {
		//处理自产生区块任务，包括：区块
		err := s.NewBlockTask(*cfg, log)
		if err != nil {
			log.Errorf("NewBlockTask|err=%v", err.Error())
		}
		//<-time.After(20 * time.Second)
		select {
		case <-ctx.Done():
			interrupt = false
			break
		case <-notify:
			continue
		}
	}
}

func (s *Service) NewBlockTask(v config.BlockConfig, log *logrus.Entry) error {
	if v.BlockMin < 1 {
		panic("Min blockNumber is not less 1")
	}

	//已经下发的最新区块高度
	UsedMaxNumber, lastBlockNumber, err := s.store.GetRecentNumber(v.BlockChainCode)
	if err != nil {
		log.Errorf("GetRecentNumber|err=%v", err)
		return err
	}

	log.Printf("NewBlockTask:blockchain:%v,UsedMaxNumber=%v,lastBlockNumber=%v", v.BlockChainCode, UsedMaxNumber, lastBlockNumber)

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
		return errors.New(fmt.Sprintf("UsedMaxNumber more than BlockMax,UsedMaxNumber:%v,BlockMax:%v", UsedMaxNumber, v.BlockMax))
	}
	list := make([]*service.NodeTask, 0)

	UsedMaxNumber++

	nodeIdList, err := s.store.GetNodeId(v.BlockChainCode)
	if err != nil {
		log.Errorf("GetNodeId|err=%v", err)
		return err
	}

	l := len(nodeIdList)

	for UsedMaxNumber <= v.BlockMax {
		index := rand.Intn(l)
		task := &service.NodeTask{
			NodeId:      nodeIdList[index],
			BlockNumber: fmt.Sprintf("%v", UsedMaxNumber),
			BlockChain:  v.BlockChainCode,
			TaskType:    2,
			TaskStatus:  0,
			CreateTime:  time.Now(),
			LogTime:     time.Now(),
			Id:          time.Now().UnixNano(),
		}

		list = append(list, task)
		UsedMaxNumber++
	}

	recentNumber := UsedMaxNumber - 1

	//生成任务并保存
	msgList, err := s.store.AddNodeTask(list)
	if err != nil {
		return err
	}

	//更新最新下发的区块高度
	err = s.store.UpdateRecentNumber(v.BlockChainCode, recentNumber)
	if err != nil {
		return err
	}

	if len(msgList) > 0 {
		s.sendCh <- msgList
	}

	return nil
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/create_task", 24*time.Hour)
	kf := kafkaClient.NewEasyKafka(xg)
	sendCh := make(chan []*kafka.Message, 5)
	//receiverCh := make(chan []*kafka.Message, 5)
	store := db.NewFileTaskCreateService(config, xg)

	blockClient := make(map[int64]service.BlockChainInterface, 2)

	for _, v := range config.BlockConfigs {
		api := NewApi(v.BlockChainCode, xg, v)
		if api != nil {
			blockClient[v.BlockChainCode] = api
		} else {
			xg.Warnf("new api client is error by config. config=%+v", v)
			panic("the system does not support the chain")
		}
	}

	return &Service{
		config:      config,
		log:         xg,
		store:       store,
		kafkaClient: kf,
		sendCh:      sendCh,
		api:         blockClient,
	}
}
