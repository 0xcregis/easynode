package db

import (
	"context"
	"encoding/json"
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

// StoreExecTask key which tx:tx_txHash,receipt:receipt_txHash, block: block_number_blockHash
func (s *Service) StoreExecTask(key string, task *service.NodeTask) {
	bs, _ := json.Marshal(task)
	err := s.cacheClient.Set(context.Background(), key, string(bs))
	if err != nil {
		s.log.Errorf("StoreExecTask|err=%v", err.Error())
	}
}

func (s *Service) AddNodeTask(list []*service.NodeTask) error {
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

			bs, _ := json.Marshal(t)
			s.cacheClient.Set(context.Background(), key, string(bs))
		}
	}

	return nil
}

func (s *Service) UpdateNodeTaskStatusWithBatch(keys []string, status int) error {
	for _, v := range keys {
		s.UpdateNodeTaskStatus(v, status)
	}
	return nil
}

func (s *Service) GetTaskWithReceipt(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	return nil, nil
}

func (s *Service) GetTaskWithTx(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	return nil, nil
}

func (s *Service) GetTaskWithBlock(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	return nil, nil
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
