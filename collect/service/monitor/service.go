package monitor

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	chainConfig "github.com/uduncloud/easynode/blockchain/config"
	chainService "github.com/uduncloud/easynode/blockchain/service"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/db"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/common/util"
	"path"
	"time"
)

type Service struct {
	cfg         config.Config
	log         *logrus.Entry
	taskStore   map[int64]service.StoreTaskInterface
	kafka       *kafkaClient.EasyKafka
	apis        map[int64]chainService.API
	kafkaSender map[int64]chan []*kafka.Message
}

// 监控服务
func (s *Service) Start() {

	go func() {
		for _, v := range s.cfg.Chains {
			broker := fmt.Sprintf("%v:%v", v.TaskKafka.Host, v.TaskKafka.Port)
			s.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, s.kafkaSender[int64(v.BlockChainCode)], nil)
		}
	}()

	//清理日志
	s.clearLog()

	//错误交易重试
	s.CheckErrTx()

	//构建合约数据
	s.CheckContract()
}

func (s *Service) CheckContract() {

	go func() {

		for true {
			<-time.After(30 * time.Minute)

			for blockchain, store := range s.taskStore {

				list, err := store.GetAllKeyForContract(blockchain, "")
				if err != nil {
					continue
				}

				for _, contract := range list {
					data, _ := store.GetContract(blockchain, contract)
					if len(data) < 1 {
						//todo 合约无效，需要刷新
						s.getToken(blockchain, contract, contract)
					}

				}

			}

		}

	}()

}

func (s *Service) CheckErrTx() {
	go func() {

		for true {
			<-time.After(20 * time.Minute)

			for blockchain, store := range s.taskStore {

				list, err := store.GetAllKeyForErrTx(blockchain, "")
				if err != nil {
					continue
				}

				tempList := make([]*kafka.Message, 0, 10)

				for _, hash := range list {
					data, err := store.DelErrTxNodeTask(blockchain, hash)
					if err != nil {
						continue
					}

					//todo 重发交易任务
					var v service.NodeTask
					_ = json.Unmarshal([]byte(data), &v)
					v.CreateTime = time.Now()
					v.LogTime = time.Now()
					if v.BlockChain < 1 {
						v.BlockChain = int(blockchain)
					}
					v.Id = time.Now().UnixNano()
					bs, _ := json.Marshal(v)
					msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: 0, Key: []byte(v.NodeId), Value: bs}
					tempList = append(tempList, msg)

					if len(tempList) > 10 {
						s.kafkaSender[blockchain] <- tempList
						tempList = tempList[len(tempList):]
					}
				}

				s.kafkaSender[blockchain] <- tempList
			}

		}

	}()
}

func (s *Service) clearLog() {
	go func() {
		for true {
			<-time.After(7 * time.Hour)
			p := s.cfg.LogConfig.Path
			d := s.cfg.LogConfig.Delay

			h := time.Duration(d*24) * time.Hour
			t := time.Now().Add(-h)

			for i := 0; i < 5; i++ {
				datePath := t.Format(service.DateFormat)

				datePath = fmt.Sprintf("%v%v", datePath, "0000")
				cmdLog := fmt.Sprintf("%v_%v", "cmd_log", datePath)
				_ = util.DeleteFile(path.Join(p, cmdLog))

				chainInfoLog := fmt.Sprintf("%v_%v", "chain_info_log", datePath)
				_ = util.DeleteFile(path.Join(p, chainInfoLog))

				t = t.Add(-24 * time.Hour)
			}
		}

	}()
}

func (s *Service) Stop() {
	panic("implement me")
}

func (s *Service) getToken(blockChain int64, from string, contract string) {

	token, err := s.apis[blockChain].TokenBalance(blockChain, from, contract, "")
	if err != nil {
		s.log.Errorf("TokenBalance fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
		return
	}

	if len(token) > 0 {
		err = s.taskStore[blockChain].StoreContract(blockChain, contract, token)
		if err != nil {
			s.log.Warnf("StoreContract fail: blockchain:%v,contract:%v,err:%v", blockChain, contract, err.Error())
		}
		return
	}
}

func NewService(config *config.Config) *Service {
	log := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/monitor", config.LogConfig.Path), 24*time.Hour)
	mp := make(map[int64]service.StoreTaskInterface, 2)
	sender := make(map[int64]chan []*kafka.Message, 2)
	for _, v := range config.Chains {
		store := db.NewTaskCacheService(v, log)
		mp[int64(v.BlockChainCode)] = store
		sender[int64(v.BlockChainCode)] = make(chan []*kafka.Message)
	}
	kf := kafkaClient.NewEasyKafka(log)

	apis := make(map[int64]chainService.API, 2)

	for _, v := range config.Chains {
		if v.TxTask != nil {
			list := make([]*chainConfig.NodeCluster, 0, 4)
			for _, t := range v.TxTask.FromCluster {
				temp := &chainConfig.NodeCluster{
					NodeUrl:   t.Host,
					NodeToken: t.Key,
					Weight:    t.Weight,
				}
				list = append(list, temp)
			}

			api := chainService.NewApi(int64(v.BlockChainCode), list, log)
			apis[int64(v.BlockChainCode)] = api
		}
	}

	return &Service{
		cfg:         *config,
		kafka:       kf,
		taskStore:   mp,
		kafkaSender: sender,
		apis:        apis,
		log:         log.WithField("model", "monitor"),
	}
}
