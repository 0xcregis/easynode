package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
)

var (
	NodeKey    = "nodeKey_%v"
	BlockChain = "blockChain"
)

type TaskCreateFile struct {
	config *config.Config
	log    *xlog.XLog
	client *redis.Client
	//sendCh chan []*kafka.Message
}

func NewFileTaskCreateService(config *config.Config, xg *xlog.XLog) service.StoreTaskInterface {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", config.Redis.Addr, config.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &TaskCreateFile{
		config: config,
		log:    xg,
		client: client,
		//sendCh: sendCh,
	}
}

func (t *TaskCreateFile) GetNodeId(blockChainCode int64) ([]string, error) {
	list, err := t.client.HKeys(context.Background(), fmt.Sprintf(NodeKey, blockChainCode)).Result()
	if err != nil {
		return nil, err
	}
	if len(list) < 1 {
		return nil, errors.New("no record")
	}
	_ = t.client.Del(context.Background(), fmt.Sprintf(NodeKey, blockChainCode)).Err()
	return list, nil
}

func (t *TaskCreateFile) AddNodeTask(list []*service.NodeTask) ([]*kafka.Message, error) {
	resultList := make([]*kafka.Message, 0, 5)
	for _, v := range list {
		bs, _ := json.Marshal(v)
		msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: 0, Key: []byte(v.NodeId), Value: bs}
		resultList = append(resultList, msg)
	}
	//t.sendCh <- resultList
	return resultList, nil
}

func (t *TaskCreateFile) UpdateLastNumber(blockChainCode int64, latestNumber int64) error {

	err := t.client.HSet(context.Background(), BlockChain, fmt.Sprintf("latestNumber_%v", blockChainCode), latestNumber).Err()
	if err != nil {
		return err
	}
	return nil

	//mp := make(map[int64]*service.BlockNumber, 2)
	//
	//bs, err := util.ReadLatestBlock()
	//if err != nil {
	//	bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode, LogTime: time.Now()}
	//	mp[blockChainCode] = &bn
	//	t.log.Errorf("UpdateLastNumber|ReadLatestBlock err=%v", err)
	//	//return err
	//} else {
	//	_ = json.Unmarshal(bs, &mp)
	//	if _, ok := mp[blockChainCode]; ok {
	//		mp[blockChainCode].LatestNumber = latestNumber
	//		mp[blockChainCode].LogTime = time.Now()
	//	} else {
	//		bn := service.BlockNumber{LatestNumber: latestNumber, ChainCode: blockChainCode, LogTime: time.Now()}
	//		mp[blockChainCode] = &bn
	//	}
	//}
	//
	//bs, _ = json.Marshal(mp)
	//err = util.WriteLatestBlock(string(bs))
	//if err != nil {
	//	t.log.Errorf("UpdateLastNumber|WriteLatestBlock err=%v", err)
	//	return err
	//}
	//return nil
}

func (t *TaskCreateFile) UpdateRecentNumber(blockChainCode int64, recentNumber int64) error {

	err := t.client.HSet(context.Background(), BlockChain, fmt.Sprintf("recentNumber_%v", blockChainCode), recentNumber).Err()
	if err != nil {
		return err
	}

	return nil
	//
	//mp := make(map[int64]*service.BlockNumber, 2)
	//bs, err := util.ReadLatestBlock()
	//if err != nil {
	//	bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode, LogTime: time.Now()}
	//	mp[blockChainCode] = &bn
	//	t.log.Errorf("UpdateRecentNumber|ReadLatestBlock err=%v", err)
	//} else {
	//	_ = json.Unmarshal(bs, &mp)
	//	if _, ok := mp[blockChainCode]; ok {
	//		mp[blockChainCode].RecentNumber = recentNumber
	//		mp[blockChainCode].LogTime = time.Now()
	//	} else {
	//		bn := service.BlockNumber{RecentNumber: recentNumber, ChainCode: blockChainCode, LogTime: time.Now()}
	//		mp[blockChainCode] = &bn
	//	}
	//}
	//
	//bs, _ = json.Marshal(mp)
	//err = util.WriteLatestBlock(string(bs))
	//if err != nil {
	//	t.log.Errorf("UpdateRecentNumber|WriteLatestBlock err=%v", err)
	//	return err
	//}
	//
	//return nil
}

func (t *TaskCreateFile) GetRecentNumber(blockCode int64) (int64, int64, error) {

	//var recentNumber,LatestNumber int64

	recentNumber, err := t.client.HGet(context.Background(), BlockChain, fmt.Sprintf("recentNumber_%v", blockCode)).Int64()
	if err != nil {
		t.log.Warnf("GetRecentNumber,err=%v",err.Error())
	}

	latestNumber, err := t.client.HGet(context.Background(), BlockChain, fmt.Sprintf("latestNumber_%v", blockCode)).Int64()
	if err != nil {
		t.log.Warnf("GetRecentNumber,err=%v",err.Error())
	}

	return recentNumber, latestNumber, nil

	//bs, err := util.ReadLatestBlock()
	//if err != nil {
	//	return 0, 0, err
	//}
	//
	//mp := make(map[int64]*service.BlockNumber, 2)
	//_ = json.Unmarshal(bs, &mp)
	//
	//if v, ok := mp[blockCode]; ok {
	//	return v.RecentNumber, v.LatestNumber, nil
	//} else {
	//	return 0, 0, errors.New("no record")
	//}

}
