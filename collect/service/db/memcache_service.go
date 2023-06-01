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
	"sync"
	"time"
)

type Service struct {
	log         *xlog.XLog
	lock        *sync.RWMutex
	cacheClient *redis.Client
}

func (s *Service) DelErrTxNodeTask(blockchain int64, key string) (string, error) {
	_, task, err := s.GetErrTxNodeTask(blockchain, key)
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	err = s.cacheClient.HDel(context.Background(), fmt.Sprintf("errTx_%v", blockchain), key).Err()
	if err != nil {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) GetAllKeyForContract(blockchain int64, key string) ([]string, error) {
	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf("contract_%v", blockchain)).Result()
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

func (s *Service) GetAllKeyForErrTx(blockchain int64, key string) ([]string, error) {
	list, err := s.cacheClient.HKeys(context.Background(), fmt.Sprintf("errTx_%v", blockchain)).Result()
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

	if err != nil && len(d) > 0 {
		var task service.NodeTask
		_ = json.Unmarshal([]byte(d), &task)
		task.LogTime = time.Now()
		data = task
	}

	mp := make(map[string]any, 2)
	mp["count"] = c + 1
	mp["data"] = data

	bs, _ := json.Marshal(mp)
	err = s.cacheClient.HSet(context.Background(), fmt.Sprintf("errTx_%v", blockchain), key, string(bs)).Err()
	if err != nil {
		s.log.Warnf("StoreErrTxNodeTask|err=%v", err.Error())
		return err
	}
	return nil
}

func (s *Service) GetErrTxNodeTask(blockchain int64, key string) (int64, string, error) {

	has, err := s.cacheClient.HExists(context.Background(), fmt.Sprintf("errTx_%v", blockchain), key).Result()
	if err != nil {
		return 0, "", err
	}

	if !has {
		return 0, "", errors.New("no record")
	}

	task, err := s.cacheClient.HGet(context.Background(), fmt.Sprintf("errTx_%v", blockchain), key).Result()
	if err != nil || task == "" {
		return 0, "", errors.New("no record")
	}
	r := gjson.Parse(task)

	count := r.Get("count").Int()
	data := r.Get("data").String()

	return count, data, nil
}

func (s *Service) StoreContract(blockchain int64, contract string, data string) error {
	err := s.cacheClient.HSet(context.Background(), fmt.Sprintf("contract_%v", blockchain), contract, data).Err()
	if err != nil {
		s.log.Warnf("StoreContract|err=%v", err.Error())
		return err
	}
	return nil
}

func (s *Service) GetContract(blockchain int64, contract string) (string, error) {
	task, err := s.cacheClient.HGet(context.Background(), fmt.Sprintf("contract_%v", blockchain), contract).Result()
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) GetNodeTask(key string) (*service.NodeTask, error) {
	task, err := s.cacheClient.Get(context.Background(), key).Result()
	if err == nil && task != "" {
		t := service.NodeTask{}
		_ = json.Unmarshal([]byte(task), &t)
		return &t, nil
	} else {
		return nil, errors.New("no record")
	}
}

func (s *Service) ResetNodeTask(oldKey, key string) error {
	task, err := s.GetNodeTask(oldKey)
	if err != nil {
		return err
	}
	_ = s.cacheClient.Del(context.Background(), oldKey)
	s.StoreNodeTask(key, task)
	return nil
}

// StoreExecTask key which tx:tx_txHash,receipt:receipt_txHash, block: block_number_blockHash
func (s *Service) StoreNodeTask(key string, task *service.NodeTask) {
	bs, _ := json.Marshal(task)
	err := s.cacheClient.Set(context.Background(), key, string(bs), 24*time.Hour).Err()
	if err != nil {
		s.log.Errorf("StoreNodeTask|err=%v", err.Error())
	}
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

func (s *Service) UpdateNodeTaskStatus(key string, status int) error {

	if status == 1 {
		err := s.cacheClient.Del(context.Background(), key).Err()
		if err != nil {
			s.log.Errorf("UpdateNodeTaskStatus|err=%v", err.Error())
			return err
		}
	} else {
		task, err := s.cacheClient.Get(context.Background(), key).Result()
		if err == nil && task != "" {
			t := service.NodeTask{}
			_ = json.Unmarshal([]byte(task), &t)
			t.TaskStatus = status
			t.LogTime = time.Now()
			bs, _ := json.Marshal(t)
			s.cacheClient.Set(context.Background(), key, string(bs), 24*time.Hour)
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
	//opt := func(o *store.Options) {
	//	o.Expiration = 1 * time.Hour
	//}
	//
	//var store store.StoreInterface

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.Redis.Addr, cfg.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	//if cfg.Redis == nil {
	//	c, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Minute))
	//	store = bigcacheStore.NewBigcache(c, opt)
	//} else {
	//	store = redisStore.NewRedis(redis.NewClient(&redis.Options{
	//		Addr:         fmt.Sprintf("%v:%v", cfg.Redis.Addr, cfg.Redis.Port),
	//		DB:           cfg.Redis.DB,
	//		DialTimeout:  time.Second,
	//		ReadTimeout:  time.Second,
	//		WriteTimeout: time.Second}), opt)
	//}

	//cacheClient := cache.New[string](store)
	return &Service{
		log:         x,
		cacheClient: client,
		lock:        &sync.RWMutex{},
	}
}
