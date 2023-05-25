package db

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	bigcacheStore "github.com/eko/gocache/store/bigcache/v4"
	redisStore "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"sync"
	"time"
)

type Service struct {
	log         *xlog.XLog
	lock        *sync.RWMutex
	cacheClient *cache.Cache[string]
	kafkaCh     chan []*kafka.Message
}

func (s *Service) StoreErrTxNodeTask(blockchain int64, key string, data any) error {
	bs, _ := json.Marshal(data)
	err := s.cacheClient.Set(context.Background(), fmt.Sprintf("errTx_%v_%v", blockchain, key), string(bs), func(o *store.Options) {
		o.Expiration = 0
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetErrTxNodeTask(blockchain int64, key string) (string, error) {
	task, err := s.cacheClient.Get(context.Background(), fmt.Sprintf("errTx_%v_%v", blockchain, key))
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) StoreContract(blockchain int64, contract string, data string) error {
	err := s.cacheClient.Set(context.Background(), fmt.Sprintf("contract_%v_%v", blockchain, contract), data, func(o *store.Options) {
		o.Expiration = 0
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetContract(blockchain int64, contract string) (string, error) {
	task, err := s.cacheClient.Get(context.Background(), fmt.Sprintf("contract_%v_%v", blockchain, contract))
	if err != nil || task == "" {
		return "", errors.New("no record")
	}

	return task, nil
}

func (s *Service) GetNodeTask(key string) (*service.NodeTask, error) {
	task, err := s.cacheClient.Get(context.Background(), key)
	if err == nil && task != "" {
		t := service.NodeTask{}
		json.Unmarshal([]byte(task), &t)
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
	_ = s.cacheClient.Delete(context.Background(), oldKey)
	s.StoreNodeTask(key, task)
	return nil
}

// StoreExecTask key which tx:tx_txHash,receipt:receipt_txHash, block: block_number_blockHash
func (s *Service) StoreNodeTask(key string, task *service.NodeTask) {
	bs, _ := json.Marshal(task)
	err := s.cacheClient.Set(context.Background(), key, string(bs))
	if err != nil {
		s.log.Errorf("StoreNodeTask|err=%v", err.Error())
	}
}

func (s *Service) SendNodeTask(list []*service.NodeTask) error {
	resultList := make([]*kafka.Message, 0)
	for _, v := range list {
		v.CreateTime = time.Now()
		v.LogTime = time.Now()
		v.Id = time.Now().UnixNano()
		bs, _ := json.Marshal(v)
		msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", v.BlockChain), Partition: 0, Key: []byte(v.NodeId), Value: bs}
		resultList = append(resultList, msg)
	}
	s.kafkaCh <- resultList
	return nil
}

func (s *Service) UpdateNodeTaskStatus(key string, status int) error {

	if status == 1 {
		err := s.cacheClient.Delete(context.Background(), key)
		if err != nil {
			s.log.Errorf("UpdateNodeTaskStatus|err=%v", err.Error())
			return err
		}
	} else {
		task, err := s.cacheClient.Get(context.Background(), key)
		if err == nil && task != "" {

			t := service.NodeTask{}
			json.Unmarshal([]byte(task), &t)
			t.TaskStatus = status
			t.LogTime = time.Now()
			bs, _ := json.Marshal(t)
			s.cacheClient.Set(context.Background(), key, string(bs))
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

func NewTaskCacheService(cfg *config.Chain, kafkaCh chan []*kafka.Message, x *xlog.XLog) service.StoreTaskInterface {
	opt := func(o *store.Options) {
		o.Expiration = 1 * time.Hour
	}

	var store store.StoreInterface
	if cfg.Redis == nil {
		c, _ := bigcache.New(context.Background(), bigcache.DefaultConfig(30*time.Minute))
		store = bigcacheStore.NewBigcache(c, opt)
	} else {
		store = redisStore.NewRedis(redis.NewClient(&redis.Options{
			Addr:         fmt.Sprintf("%v:%v", cfg.Redis.Addr, cfg.Redis.Port),
			DB:           cfg.Redis.DB,
			DialTimeout:  time.Second,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second}), opt)
	}

	cacheClient := cache.New[string](store)
	return &Service{
		log:         x,
		cacheClient: cacheClient,
		kafkaCh:     kafkaCh,
		lock:        &sync.RWMutex{},
	}
}
