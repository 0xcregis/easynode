package push

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"io/ioutil"
	"time"
)

type Server struct {
	log    *xlog.XLog
	chains map[int64]*config.Chain
	db     service.DbMonitorAddressInterface
}

func NewServer(cfg *config.Config, log *xlog.XLog) *Server {
	db := NewChService(cfg, log)
	mp := make(map[int64]*config.Chain, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
	}
	return &Server{
		db:     db,
		log:    log,
		chains: mp,
	}
}

func (s *Server) NewToken(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	if _, ok := s.chains[blockChain]; !ok {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	token, err := uuid.NewV4()
	if err != nil {
		s.Error(c, c.Request.URL.Path, errors.New("create token failure").Error())
		return
	}

	//保存token
	//todo 是否需要保存

	s.Success(c, token.String(), c.Request.URL.Path)
}

func (s *Server) MonitorAddress(c *gin.Context) {

	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	blockChain := r.Get("blockChain").Int()
	if _, ok := s.chains[blockChain]; !ok {
		s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
		return
	}

	addr := r.Get("address").String()
	token := r.Get("token").String()
	txType := r.Get("txType").Int()

	addressTask := &service.MonitorAddress{BlockChain: blockChain, Address: addr, Token: token, TxType: fmt.Sprintf("%v", txType), Id: time.Now().UnixMilli()}
	err = s.db.AddMonitorAddress(blockChain, addressTask)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, addressTask.Id, c.Request.URL.Path)
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
