package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"

	"github.com/0xcregis/easynode/taskapi"
	"github.com/0xcregis/easynode/taskapi/config"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

const (
	NodeKey = "nodeKey_%v"
)

type Server struct {
	log        *xlog.XLog
	blockChain []int64
	cfg        *config.Config
	handler    taskapi.TaskApiInterface
	client     *redis.Client
}

func NewServer(cfg *config.Config, blockChain []int64, log *xlog.XLog) *Server {
	h := NewTaskHandler(cfg, log)
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", cfg.Redis.Addr, cfg.Redis.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		panic(err.Error())
	}
	return &Server{
		handler:    h,
		cfg:        cfg,
		log:        log,
		client:     client,
		blockChain: blockChain,
	}
}

func (s *Server) getNodeId(blockChainCode int64) (string, error) {
	list, err := s.client.HKeys(context.Background(), fmt.Sprintf(NodeKey, blockChainCode)).Result()
	if err != nil {
		return "", err
	}
	if l := len(list); l > 0 {
		return list[rand.Intn(l)], nil
	} else {
		return "", errors.New("no record")
	}
}

func (s *Server) ExistNodeId(blockChainCode int64, nodeId string) bool {
	exist, err := s.client.HExists(context.Background(), fmt.Sprintf(NodeKey, blockChainCode), nodeId).Result()
	if err != nil {
		return false
	}
	return exist
}

func (s *Server) PushBlockTask(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	has := false
	for _, c := range s.blockChain {
		if c == blockChain {
			has = true
			break
		}
	}
	if !has {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	blockHash := r.Get("blockHash").String()
	blockNumber := r.Get("blockNumber").String()

	if len(blockHash) < 10 && len(blockNumber) < 2 {
		s.Error(c, c.Request.URL.Path, errors.New("blockHash or blockNumber is wrong").Error())
		return
	}

	nodeId := r.Get("nodeId").String()
	if len(nodeId) < 1 {
		nodeId, err = s.getNodeId(blockChain)
		if err != nil {
			s.Error(c, c.Request.URL.Path, errors.New("not found nodeId").Error())
			return
		}
	} else {
		if !s.ExistNodeId(blockChain, nodeId) {
			s.Error(c, c.Request.URL.Path, errors.New("the nodeId is unusable").Error())
			return
		}
	}

	task := &taskapi.NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 2, TaskStatus: 0, NodeId: nodeId}
	err = s.handler.SendNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushTxTask(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	has := false
	for _, c := range s.blockChain {
		if c == blockChain {
			has = true
			break
		}
	}
	if !has {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	txHash := r.Get("txHash").String()

	if len(txHash) < 10 {
		s.Error(c, c.Request.URL.Path, errors.New("txHash is wrong").Error())
		return
	}
	nodeId := r.Get("nodeId").String()
	if len(nodeId) < 1 {
		nodeId, err = s.getNodeId(blockChain)
		if err != nil {
			s.Error(c, c.Request.URL.Path, errors.New("not found nodeId").Error())
			return
		}
	} else {
		if !s.ExistNodeId(blockChain, nodeId) {
			s.Error(c, c.Request.URL.Path, errors.New("the nodeId is unusable").Error())
			return
		}
	}

	task := &taskapi.NodeTask{BlockChain: blockChain, TxHash: txHash, TaskType: 1, TaskStatus: 0, NodeId: nodeId}
	err = s.handler.SendNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushTxsTask(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	has := false
	for _, c := range s.blockChain {
		if c == blockChain {
			has = true
			break
		}
	}
	if !has {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	blockHash := r.Get("blockHash").String()
	blockNumber := r.Get("blockNumber").String()

	if len(blockHash) < 10 && len(blockNumber) < 2 {
		s.Error(c, c.Request.URL.Path, errors.New("blockHash or blockNumber is wrong").Error())
		return
	}
	nodeId := r.Get("nodeId").String()
	if len(nodeId) < 1 {
		nodeId, err = s.getNodeId(blockChain)
		if err != nil {
			s.Error(c, c.Request.URL.Path, errors.New("not found nodeId").Error())
			return
		}
	} else {
		if !s.ExistNodeId(blockChain, nodeId) {
			s.Error(c, c.Request.URL.Path, errors.New("the nodeId is unusable").Error())
			return
		}
	}

	task := &taskapi.NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 4, TaskStatus: 0, NodeId: nodeId}
	err = s.handler.SendNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushReceiptTask(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	has := false
	for _, c := range s.blockChain {
		if c == blockChain {
			has = true
			break
		}
	}
	if !has {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	txHash := r.Get("txHash").String()

	if len(txHash) < 10 {
		s.Error(c, c.Request.URL.Path, errors.New("txHash  is wrong").Error())
		return
	}

	nodeId := r.Get("nodeId").String()
	if len(nodeId) < 1 {
		nodeId, err = s.getNodeId(blockChain)
		if err != nil {
			s.Error(c, c.Request.URL.Path, errors.New("not found nodeId").Error())
			return
		}
	} else {
		if !s.ExistNodeId(blockChain, nodeId) {
			s.Error(c, c.Request.URL.Path, errors.New("the nodeId is unusable").Error())
			return
		}
	}

	task := &taskapi.NodeTask{BlockChain: blockChain, TxHash: txHash, TaskType: 3, TaskStatus: 0, NodeId: nodeId}
	err = s.handler.SendNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushReceiptsTask(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	has := false
	for _, c := range s.blockChain {
		if c == blockChain {
			has = true
			break
		}
	}
	if !has {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	blockHash := r.Get("blockHash").String()
	blockNumber := r.Get("blockNumber").String()

	if len(blockHash) < 10 && len(blockNumber) < 2 {
		s.Error(c, c.Request.URL.Path, errors.New("blockHash or blockNumber is wrong").Error())
		return
	}

	nodeId := r.Get("nodeId").String()
	if len(nodeId) < 1 {
		nodeId, err = s.getNodeId(blockChain)
		if err != nil {
			s.Error(c, c.Request.URL.Path, errors.New("not found nodeId").Error())
			return
		}
	} else {
		if !s.ExistNodeId(blockChain, nodeId) {
			s.Error(c, c.Request.URL.Path, errors.New("the nodeId is unusable").Error())
			return
		}
	}

	task := &taskapi.NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 5, TaskStatus: 0, NodeId: nodeId}
	err = s.handler.SendNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

const (
	SUCCESS = 0
	FAIL    = 1
)

func (s *Server) Success(c *gin.Context, resp interface{}, path string) {
	s.log.Printf("path=%v,body=%v\n", path, resp)
	mp := make(map[string]interface{})
	mp["code"] = SUCCESS
	mp["data"] = resp
	c.JSON(200, mp)
}

func (s *Server) Error(c *gin.Context, path string, err string) {
	s.log.Printf("path=%v,err=%v\n", path, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["data"] = err
	c.JSON(200, mp)
}
