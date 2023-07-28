package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/common/util"
	"github.com/segmentio/kafka-go"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../../cmd/collect/config_tron.json")
	nodeId, _ := util.GetLocalNodeId(cfg.KeyPath)
	return NewService(&cfg, nodeId)
}

func TestService_Start(t *testing.T) {
	s := Init()
	s.Start(context.Background())
}

func TestService_CheckErrTx(t *testing.T) {
	s := Init()
	for blockchain, store := range s.taskStore {

		list, err := store.GetAllKeyForErrTx(blockchain)
		if err != nil {
			continue
		}

		tempList := make([]*kafka.Message, 0, 10)

		for _, hash := range list {
			count, data, err := store.GetErrTxNodeTask(blockchain, hash)
			if err != nil || count >= 5 {
				continue
			}

			//todo 重发交易任务
			var v collect.NodeTask
			_ = json.Unmarshal([]byte(data), &v)

			//清理 已经重试成功的交易
			if count < 5 && time.Since(v.LogTime) >= 24*time.Hour {
				_, _ = store.DelErrTxNodeTask(blockchain, hash)
				continue
			}

			v.CreateTime = time.Now()
			v.LogTime = time.Now()
			if v.BlockChain < 1 {
				v.BlockChain = int(blockchain)
			}
			v.Id = time.Now().UnixNano()
			v.TaskStatus = 0
			bs, _ := json.Marshal(v)
			msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: 0, Key: []byte(v.NodeId), Value: bs}
			tempList = append(tempList, msg)

			if len(tempList) > 10 {
				//s.kafkaSender[blockchain] <- tempList
				tempList = tempList[len(tempList):]
			}
		}

		//s.kafkaSender[blockchain] <- tempList
	}
}

func TestService_CheckNodeTask(t *testing.T) {
	s := Init()
	for blockchain, store := range s.taskStore {

		list, err := store.GetAllKeyForNodeTask(blockchain)
		if err != nil {
			continue
		}

		tempList := make([]*kafka.Message, 0, 10)

		for _, hash := range list {
			count, task, err := store.GetNodeTask(blockchain, hash)
			if err != nil || count >= 5 || task == nil {
				continue
			}

			//清理 已经重试成功的交易
			if count < 5 && time.Since(task.LogTime) >= 24*time.Hour {
				_, _, _ = store.DelNodeTask(blockchain, hash)
				continue
			}

			task.CreateTime = time.Now()
			task.LogTime = time.Now()
			if task.BlockChain < 1 {
				task.BlockChain = int(blockchain)
			}
			task.Id = time.Now().UnixNano()
			task.TaskStatus = 0
			bs, _ := json.Marshal(task)
			msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", blockchain), Partition: 0, Key: []byte(task.NodeId), Value: bs}
			tempList = append(tempList, msg)

			if len(tempList) > 10 {
				//s.kafkaSender[blockchain] <- tempList
				tempList = tempList[len(tempList):]
			}
		}

		//s.kafkaSender[blockchain] <- tempList
	}
}

func TestService_CheckContract(t *testing.T) {

}
