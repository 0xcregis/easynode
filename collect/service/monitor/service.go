package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	chainConfig "github.com/0xcregis/easynode/blockchain/config"
	chainService "github.com/0xcregis/easynode/blockchain/service"
	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	kafkaClient "github.com/0xcregis/easynode/common/kafka"
	"github.com/0xcregis/easynode/common/util"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

const (
	DateFormat    = "20060102"
	NodeTaskTopic = "task_%v"
)

type Service struct {
	cfg         config.Config
	log         *logrus.Entry
	nodeId      string
	taskStore   map[int64]collect.StoreTaskInterface
	kafka       *kafkaClient.EasyKafka
	apis        map[int64]blockchain.API
	kafkaSender map[int64]chan []*kafka.Message
}

// Start 监控服务
func (s *Service) Start(ctx context.Context) {
	//开启Kafka
	s.StartKafka()

	//更新nodeId
	s.UpdateNodeId()

	//公链节点健康检查
	s.CheckClusterHealth()

	//清理日志
	s.clearLog()

	//错误交易重试
	s.CheckErrTx()

	//任务执行失败，重试5次
	s.CheckNodeTask()

	//构建合约数据
	s.CheckContract()
}

func (s *Service) StartKafka() {
	go func() {
		for _, v := range s.cfg.Chains {
			broker := fmt.Sprintf("%v:%v", v.TaskKafka.Host, v.TaskKafka.Port)
			s.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, s.kafkaSender[int64(v.BlockChainCode)], nil, context.Background(), 1)
		}
	}()
}

func (s *Service) UpdateNodeId() {
	go func(nodeId string) {
		for {
			for b, s := range s.taskStore {
				_ = s.StoreNodeId(b, nodeId, 1)
			}
			<-time.After(5 * time.Second)
		}
	}(s.nodeId)
}

func (s *Service) CheckClusterHealth() {
	go func() {
		for {
			<-time.After(1 * time.Minute)
			for chainCode, t := range s.taskStore {
				mp := make(map[string]int64, 3)
				s.rebuildCluster(t, chainCode, "tx", mp)
				s.rebuildCluster(t, chainCode, "block", mp)
				s.rebuildCluster(t, chainCode, "receipt", mp)
				_ = t.StoreClusterHealthStatus(chainCode, mp)
			}
		}
	}()
}

func (s *Service) rebuildCluster(t collect.StoreTaskInterface, blockChain int64, prefix string, r map[string]int64) {
	m1, _ := t.GetClusterNode(blockChain, prefix)
	for k, v := range m1 {
		if o, ok := r[k]; ok {
			r[k] = v + o
		} else {
			r[k] = v
		}
	}
}

func (s *Service) CheckNodeTask() {
	go func() {
		for {

			var l int64
			if s.cfg.Retry == nil || s.cfg.Retry.NodeTask == 0 {
				l = 120
			} else {
				l = s.cfg.Retry.NodeTask
			}
			l = l * int64(time.Second)

			<-time.After(time.Duration(l))

			for chainCode, store := range s.taskStore {

				list, err := store.GetAllKeyForNodeTask(chainCode)
				if err != nil {
					continue
				}

				tempList := make([]*kafka.Message, 0, 10)

				for _, hash := range list {
					count, task, err := store.GetNodeTask(chainCode, hash)
					if err != nil || count >= 5 || task == nil {
						continue
					}

					//清理 已经重试成功的交易
					//if count < 5 && time.Since(task.LogTime) >= 24*time.Hour {
					//	_, _, _ = store.DelNodeTask(blockchain, hash)
					//	continue
					//}

					task.CreateTime = time.Now()
					task.LogTime = time.Now()
					if task.BlockChain < 1 {
						task.BlockChain = int(chainCode)
					}
					task.Id = time.Now().UnixNano()
					task.TaskStatus = 0
					task.NodeId = s.nodeId
					s.log.Printf("NodeTask|task:%+v", task)

					bs, _ := json.Marshal(task)
					msg := &kafka.Message{Topic: fmt.Sprintf(NodeTaskTopic, chainCode), Partition: 0, Key: []byte(task.NodeId), Value: bs}
					tempList = append(tempList, msg)

					if len(tempList) > 10 {
						s.kafkaSender[chainCode] <- tempList
						tempList = tempList[len(tempList):]
					}
				}

				s.kafkaSender[chainCode] <- tempList
			}

		}
	}()
}

func (s *Service) CheckContract() {

	go func() {

		for {
			//var l int64
			//if s.cfg.Retry == nil || s.cfg.Retry.Contract == 0 {
			//	l = 120
			//} else {
			//	l = s.cfg.Retry.Contract
			//}
			//l = l * int64(time.Second)

			<-time.After(30 * time.Minute)
			for chainCode, store := range s.taskStore {

				list, err := store.GetAllKeyForContract(chainCode)
				if err != nil {
					continue
				}

				for _, contract := range list {
					data, _ := store.GetContract(chainCode, contract)
					if len(data) < 1 {
						//todo 合约无效，需要刷新
						s.getToken(chainCode, contract, contract)
					}

				}

			}
		}

	}()

}

func (s *Service) CheckErrTx() {
	go func() {

		for {
			var l int64
			if s.cfg.Retry == nil || s.cfg.Retry.ErrTx == 0 {
				l = 120
			} else {
				l = s.cfg.Retry.ErrTx
			}
			l = l * int64(time.Second)

			<-time.After(time.Duration(l))

			for chainCode, store := range s.taskStore {

				list, err := store.GetAllKeyForErrTx(chainCode)
				if err != nil {
					continue
				}

				tempList := make([]*kafka.Message, 0, 10)

				for _, hash := range list {
					count, task, err := store.GetErrTxNodeTask(chainCode, hash)
					if err != nil || count >= 5 {
						continue
					}

					//重发交易任务
					//var v collect.NodeTask
					//_ = json.Unmarshal([]byte(data), &v)

					//清理 已经重试成功的交易
					//if count < 5 && time.Since(v.LogTime) >= 24*time.Hour {
					//	_, _ = store.DelErrTxNodeTask(blockchain, hash)
					//	continue
					//}

					task.CreateTime = time.Now()
					task.LogTime = time.Now()
					task.NodeId = s.nodeId
					if task.BlockChain < 1 {
						task.BlockChain = int(chainCode)
					}
					task.Id = time.Now().UnixNano()
					task.TaskStatus = 0
					s.log.Printf("ErrTx|task:%+v", task)

					bs, _ := json.Marshal(task)
					msg := &kafka.Message{Topic: fmt.Sprintf(NodeTaskTopic, chainCode), Partition: 0, Key: []byte(task.NodeId), Value: bs}
					tempList = append(tempList, msg)

					if len(tempList) > 10 {
						s.kafkaSender[chainCode] <- tempList
						tempList = tempList[len(tempList):]
					}
				}

				s.kafkaSender[chainCode] <- tempList
			}

		}

	}()
}

func (s *Service) clearLog() {
	go func() {
		for {
			<-time.After(7 * time.Hour)
			p := s.cfg.LogConfig.Path
			d := s.cfg.LogConfig.Delay

			h := time.Duration(d*24) * time.Hour
			t := time.Now().Add(-h)

			for i := 0; i < 5; i++ {
				datePath := t.Format(DateFormat)

				datePath = fmt.Sprintf("%v%v", datePath, "0000")
				cmdLog := fmt.Sprintf("%v_%v", "cmd_log", datePath)
				_ = util.DeleteFile(path.Join(p, cmdLog))

				chainInfoLog := fmt.Sprintf("%v_%v", "chain_info_log", datePath)
				_ = util.DeleteFile(path.Join(p, chainInfoLog))

				monitorLog := fmt.Sprintf("%v_%v", "monitor_log", datePath)
				_ = util.DeleteFile(path.Join(p, monitorLog))

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

func NewService(config *config.Config, nodeId string) *Service {
	log := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/monitor", config.LogConfig.Path), 24*time.Hour)
	x := log.WithField("model", "monitor")
	mp := make(map[int64]collect.StoreTaskInterface, 2)
	sender := make(map[int64]chan []*kafka.Message, 2)
	for _, v := range config.Chains {
		store := db.NewTaskCacheService2(v, x)
		mp[int64(v.BlockChainCode)] = store
		sender[int64(v.BlockChainCode)] = make(chan []*kafka.Message)
	}
	kf := kafkaClient.NewEasyKafka2(x)

	apis := make(map[int64]blockchain.API, 2)

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
		nodeId:      nodeId,
		log:         x,
	}
}
