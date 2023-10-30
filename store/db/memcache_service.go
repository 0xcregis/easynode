package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/config"
	"github.com/redis/go-redis/v9"
	"github.com/sunjiangjun/xlog"
)

var (
	MonitorKey = "monitorAddress_%v"
)

type CacheService struct {
	log         *xlog.XLog
	cacheClient map[int64]*redis.Client
}

func (s *CacheService) GetMonitorAddress(blockChain int64) ([]string, error) {

	if _, ok := s.cacheClient[blockChain]; !ok {
		return nil, fmt.Errorf("blockChain:%v does not support", blockChain)
	}
	list, err := s.cacheClient[blockChain].HKeys(context.Background(), fmt.Sprintf(MonitorKey, blockChain)).Result()
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

func (s *CacheService) SetMonitorAddress(blockChain int64, addrList []*store.MonitorAddress) error {
	if _, ok := s.cacheClient[blockChain]; !ok {
		return fmt.Errorf("blockChain:%v does not support", blockChain)
	}
	_, err := s.cacheClient[blockChain].Del(context.Background(), fmt.Sprintf(MonitorKey, blockChain)).Result()
	if err != nil {
		s.log.Warnf("GetMonitorAddress|err=%v", err.Error())
		return err
	}

	if len(addrList) < 1 {
		return fmt.Errorf("blockChain:%v ,no monitor address to save", blockChain)
	}

	mp := make(map[string]string, 20)
	for _, v := range addrList {
		mp[v.Address] = v.Address
		if len(mp) > 100 {
			_, err = s.cacheClient[blockChain].HSet(context.Background(), fmt.Sprintf(MonitorKey, blockChain), mp).Result()
			if err != nil {
				s.log.Warnf("GetMonitorAddress|err=%v", err.Error())
				//return err
			}
			mp = make(map[string]string, 20)
		}
	}

	if len(mp) < 1 {
		return fmt.Errorf("blockChain:%v ,no monitor address to save", blockChain)
	}
	_, err = s.cacheClient[blockChain].HSet(context.Background(), fmt.Sprintf(MonitorKey, blockChain), mp).Result()
	if err != nil {
		s.log.Warnf("GetMonitorAddress|err=%v", err.Error())
		return err
	}

	return nil
}

func NewCacheService(cfgs []*config.Chain, x *xlog.XLog) *CacheService {

	mp := make(map[int64]*redis.Client)
	for _, v := range cfgs {
		if v.Redis != nil {
			client := redis.NewClient(&redis.Options{
				Addr:        fmt.Sprintf("%v:%v", v.Redis.Addr, v.Redis.Port),
				Password:    "", // no password set
				DB:          0,  // use default DB
				DialTimeout: 20 * time.Second,
			})
			mp[v.BlockChain] = client
		}
	}
	return &CacheService{
		log:         x,
		cacheClient: mp,
	}
}
