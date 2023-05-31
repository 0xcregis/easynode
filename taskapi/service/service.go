package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/taskapi/config"
	"io/ioutil"
)

type Server struct {
	log        *xlog.XLog
	blockChain []int64
	nodeId     string
	cfg        *config.Config
	db         DbApiInterface
}

func NewServer(cfg *config.Config, blockChain []int64, log *xlog.XLog) *Server {
	db := NewChService(cfg, log)
	nodeId, err := util.GetLocalNodeId()
	if err != nil {
		panic(err)
	}
	return &Server{
		db:         db,
		cfg:        cfg,
		nodeId:     nodeId,
		log:        log,
		blockChain: blockChain,
	}
}

func (s *Server) PushBlockTask(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
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

	task := &NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 2, TaskStatus: 0, NodeId: s.nodeId}
	err = s.db.AddNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

//func (s *Server) PushSyncTxTask(c *gin.Context) {
//	bs, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		s.Error(c, c.Request.URL.Path, err.Error())
//		return
//	}
//
//	r := gjson.ParseBytes(bs)
//	blockChain := r.Get("blockChain").Int()
//	has := false
//	for _, c := range s.blockChain {
//		if c == blockChain {
//			has = true
//			break
//		}
//	}
//	if !has {
//		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
//		return
//	}
//
//	txHash := r.Get("txHash").String()
//
//	if len(txHash) < 10 {
//		s.Error(c, c.Request.URL.Path, errors.New("txHash is wrong").Error())
//		return
//	}
//
//	task := &NodeTask{BlockChain: blockChain, TxHash: txHash, TaskType: 1, TaskStatus: 0, NodeId: s.nodeId}
//	err = s.db.AddNodeTask(task)
//	if err != nil {
//		s.Error(c, c.Request.URL.Path, err.Error())
//		return
//	}
//
//	//开启超时定时器
//	ch := make(chan *Tx)
//	con, cancel := context.WithCancel(context.Background())
//
//	go func(ctx context.Context, respCh chan *Tx) {
//
//		has := true
//		for has {
//			//轮询的查询clickhouse
//
//			tx, err := s.db.QueryTxFromCh(blockChain, txHash)
//			if err == nil && tx != nil {
//				has = false
//				respCh <- tx
//				break
//			}
//
//			if err != nil {
//				s.log.Errorf("QueryTxFromCh|error=%v", err)
//			}
//			//设定结束条件：超时、已查询到结果
//			select {
//			case <-con.Done():
//				has = false
//			default:
//				has = true
//				<-time.After(2 * time.Second)
//			}
//
//		}
//	}(con, ch)
//
//	//返回
//	select {
//	case <-time.After(1 * time.Minute):
//		cancel()
//		s.Error(c, c.Request.URL.Path, "request is timeout")
//	case tx := <-ch:
//		cancel()
//		s.Success(c, tx, c.Request.URL.Path)
//	}
//}

func (s *Server) PushTxTask(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
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

	task := &NodeTask{BlockChain: blockChain, TxHash: txHash, TaskType: 1, TaskStatus: 0, NodeId: s.nodeId}
	err = s.db.AddNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushTxsTask(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
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

	task := &NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 4, TaskStatus: 0, NodeId: s.nodeId}
	err = s.db.AddNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushReceiptTask(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
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

	task := &NodeTask{BlockChain: blockChain, TxHash: txHash, TaskType: 3, TaskStatus: 0, NodeId: s.nodeId}
	err = s.db.AddNodeTask(task)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) PushReceiptsTask(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
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

	task := &NodeTask{BlockChain: blockChain, BlockHash: blockHash, BlockNumber: blockNumber, TaskType: 5, TaskStatus: 0, NodeId: s.nodeId}
	err = s.db.AddNodeTask(task)
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
