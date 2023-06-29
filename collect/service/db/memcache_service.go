package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ContractKey = "contract_%v"
	ErrTxKey    = "errTx_%v"
	NodeTaskKey = "nodeTask_%v"
	MonitorKey  = "monitorAddress_%v"
)

type Service struct {
	log         *xlog.XLog
	lock        *sync.RWMutex
	cacheClient *redis.Client
}

func (s *Service) GetMonitorAddress(blockChain int64) ([]string, error) {

	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf(MonitorKey, blockChain)).Result()
	if err != nil {
		s.log.Warnf("GetMonitorAddress|err=%v", err.Error())
		return nil, errors.New("no record")
	}

	if len(list) < 1 {
		s.log.Warnf("GetMonitorAddress|err=no data")
		return nil, errors.New("no record")
	}
	return list, nil
}

func (s *Service) DelNodeTask(blockchain int64, key string) (int64, *service.NodeTask, error) {
	c, task, err := s.GetNodeTask(blockchain, key)
	if err != nil {
		return c, nil, err
	}
	_ = s.cacheClient.HDel(context.Background(), fmt.Sprintf(NodeTaskKey, task.BlockChain), key)
	return c, task, nil
}

func (s *Service) DelErrTxNodeTask(blockchain int64, key string) (string, error) {
	_, task, err := s.GetErrTxNodeTask(blockchain, key)
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	err = s.cacheClient.HDel(context.Background(), fmt.Sprintf(ErrTxKey, blockchain), key).Err()
	if err != nil {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) GetAllKeyForNodeTask(blockchain int64) ([]string, error) {
	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf(NodeTaskKey, blockchain)).Result()
	if err != nil {
		s.log.Warnf("GetAllKeyForContract|err=%v", err.Error())
		return nil, errors.New("no record")
	}

	if len(list) < 1 {
		s.log.Warnf("GetAllKeyForContract|err=no data")
		return nil, errors.New("no record")
	}

	return list, nil
}

func (s *Service) GetAllKeyForContract(blockchain int64) ([]string, error) {
	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf(ContractKey, blockchain)).Result()
	if err != nil {
		s.log.Warnf("GetAllKeyForContract|err=%v", err.Error())
		return nil, errors.New("no record")
	}

	if len(list) < 1 {
		s.log.Warnf("GetAllKeyForContract|err=no data")
		return nil, errors.New("no record")
	}

	return list, nil
}

func (s *Service) GetAllKeyForErrTx(blockchain int64) ([]string, error) {
	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf(ErrTxKey, blockchain)).Result()
	if err != nil {
		s.log.Warnf("GetAllKeyForErrTx|err=%v", err.Error())
		return nil, errors.New("no record")
	}

	if len(list) < 1 {
		s.log.Warnf("GetAllKeyForErrTx|err=no data")
		return nil, errors.New("no record")
	}

	return list, nil
}

func (s *Service) StoreErrTxNodeTask(blockchain int64, key string, data any) error {
	c, d, err := s.GetErrTxNodeTask(blockchain, key)
	//已存在且错误超过5次 忽略
	if c > 5 {
		return nil
	}

	if c > 0 && len(d) > 10 {
		var task service.NodeTask
		_ = json.Unmarshal([]byte(d), &task)
		task.LogTime = time.Now()
		data = task
	}

	mp := make(map[string]any, 2)
	mp["count"] = c + 1
	mp["data"] = data

	bs, _ := json.Marshal(mp)
	err = s.cacheClient.HSet(context.Background(), fmt.Sprintf(ErrTxKey, blockchain), key, string(bs)).Err()
	if err != nil {
		s.log.Warnf("StoreErrTxNodeTask|err=%v", err.Error())
		return err
	}
	return nil
}

func (s *Service) GetErrTxNodeTask(blockchain int64, key string) (int64, string, error) {

	has, err := s.cacheClient.HExists(context.Background(), fmt.Sprintf(ErrTxKey, blockchain), key).Result()
	if err != nil {
		return 0, "", err
	}

	if !has {
		return 0, "", errors.New("no record")
	}

	task, err := s.cacheClient.HGet(context.Background(), fmt.Sprintf(ErrTxKey, blockchain), key).Result()
	if err != nil || task == "" {
		return 0, "", errors.New("no record")
	}
	r := gjson.Parse(task)

	count := r.Get("count").Int()
	data := r.Get("data").String()

	return count, data, nil
}

func (s *Service) StoreContract(blockchain int64, contract string, data string) error {
	err := s.cacheClient.HSet(context.Background(), fmt.Sprintf(ContractKey, blockchain), contract, data).Err()
	if err != nil {
		s.log.Warnf("StoreContract|err=%v", err.Error())
		return err
	}
	return nil
}

func (s *Service) GetContract(blockchain int64, contract string) (string, error) {
	task, err := s.cacheClient.HGet(context.Background(), fmt.Sprintf(ContractKey, blockchain), contract).Result()
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) SendNodeTask(list []*service.NodeTask, partitions []int64) []*kafka.Message {
	resultList := make([]*kafka.Message, 0)
	for _, v := range list {
		v.CreateTime = time.Now()
		v.LogTime = time.Now()
		v.Id = time.Now().UnixNano()
		bs, _ := json.Marshal(v)

		var p int
		if len(partitions) > 0 {
			p = rand.Intn(len(partitions))
			p = int(partitions[p])
		}

		msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: p, Key: []byte(v.NodeId), Value: bs}
		resultList = append(resultList, msg)
	}

	return resultList
}

func (s *Service) GetNodeTask(blockchain int64, key string) (int64, *service.NodeTask, error) {
	has, err := s.cacheClient.HExists(context.Background(), fmt.Sprintf(NodeTaskKey, blockchain), key).Result()
	if err != nil {
		return 0, nil, err
	}

	if !has {
		return 0, nil, errors.New("no record")
	}

	task, err := s.cacheClient.HGet(context.Background(), fmt.Sprintf(NodeTaskKey, blockchain), key).Result()
	if err != nil || task == "" {
		return 0, nil, errors.New("no record")
	}
	r := gjson.Parse(task)

	count := r.Get("count").Int()
	data := r.Get("data").String()

	t := service.NodeTask{}
	_ = json.Unmarshal([]byte(data), &t)
	return count, &t, nil
}

// StoreNodeTask StoreExecTask key which tx:tx_txHash,receipt:receipt_txHash, block: block_number_blockHash
func (s *Service) StoreNodeTask(key string, task *service.NodeTask) {

	c, d, err := s.GetNodeTask(int64(task.BlockChain), key)
	//已存在且错误超过5次 忽略
	if c > 5 {
		return
	}

	if d != nil {
		d.LogTime = task.LogTime
		d.TaskStatus = task.TaskStatus
	} else {
		d = task
	}

	mp := make(map[string]any, 2)
	mp["count"] = c + 1
	mp["data"] = d

	bs, _ := json.Marshal(mp)
	err = s.cacheClient.HSet(context.Background(), fmt.Sprintf(NodeTaskKey, task.BlockChain), key, string(bs)).Err()
	if err != nil {
		s.log.Errorf("StoreNodeTask|err=%v", err.Error())
	}
}

func (s *Service) ResetNodeTask(blockchain int64, oldKey, key string) error {
	_, task, err := s.GetNodeTask(blockchain, oldKey)
	if err != nil {
		return err
	}
	_ = s.cacheClient.HDel(context.Background(), fmt.Sprintf(NodeTaskKey, task.BlockChain), oldKey)
	s.StoreNodeTask(key, task)
	return nil
}

func (s *Service) UpdateNodeTaskStatus(key string, status int) error {

	strArray := strings.Split(key, "_")
	var blockchain int64
	if len(strArray) > 0 {
		chain := strArray[0]
		blockchain, _ = strconv.ParseInt(chain, 0, 64)
	} else {
		return errors.New(fmt.Sprintf("key is error,key:%v", key))
	}

	if status == 1 {
		err := s.cacheClient.HDel(context.Background(), fmt.Sprintf(NodeTaskKey, blockchain), key).Err()
		if err != nil {
			s.log.Errorf("UpdateNodeTaskStatus|err=%v", err.Error())
			return err
		}
	} else {

		_, task, err := s.GetNodeTask(blockchain, key)
		if err == nil && task != nil {
			task.TaskStatus = status
			task.LogTime = time.Now()
			s.StoreNodeTask(key, task)
		}

	}

	return nil
}

func (s *Service) UpdateNodeTaskStatusWithBatch(keys []string, status int) error {
	for _, v := range keys {
		_ = s.UpdateNodeTaskStatus(v, status)
	}
	return nil
}

func NewTaskCacheService(cfg *config.Chain, x *xlog.XLog) service.StoreTaskInterface {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.Redis.Addr, cfg.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Service{
		log:         x,
		cacheClient: client,
		lock:        &sync.RWMutex{},
	}
}
