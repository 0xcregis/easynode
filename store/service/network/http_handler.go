package network

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/common/util"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	db2 "github.com/uduncloud/easynode/store/service/db"
	"io"
	"strings"
	"time"
)

type Server struct {
	log    *logrus.Entry
	chains map[int64]*config.Chain
	db     service.DbMonitorAddressInterface
}

func NewServer(cfg *config.Config, log *xlog.XLog) *Server {
	db := db2.NewChService(cfg, log)
	mp := make(map[int64]*config.Chain, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
	}
	return &Server{
		db:     db,
		log:    log.WithField("model", "httpSrv"),
		chains: mp,
	}
}

func (s *Server) NewToken(c *gin.Context) {

	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)

	email := r.Get("email").String()
	token, err := uuid.NewV4()
	if err != nil {
		s.Error(c, c.Request.URL.Path, errors.New("create token failure").Error())
		return
	}

	err, _ = s.db.GetNodeTokenByEmail(email)

	if err == nil {
		//已存在
		s.Error(c, c.Request.URL.Path, errors.New("the email already has token").Error())
		return
	}

	//保存token
	err = s.db.NewToken(&service.NodeToken{Token: token.String(), Email: email, Id: time.Now().UnixMicro()})
	if err != nil {
		s.Error(c, c.Request.URL.Path, errors.New("create token failure").Error())
		return
	}

	s.Success(c, token.String(), c.Request.URL.Path)
}

func (s *Server) DelMonitorAddress(c *gin.Context) {

	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	token := r.Get("token").String()
	address := r.Get("address").String()
	blockchain := r.Get("blockChain").Int()
	err = s.db.DelMonitorAddress(blockchain, token, address)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *Server) GetMonitorAddress(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	token := r.Get("token").String()

	list, err := s.db.GetAddressByToken2(token)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, list, c.Request.URL.Path)
}

func (s *Server) MonitorAddress(c *gin.Context) {

	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)

	var blockChain int64
	if r.Get("blockChain").Exists() {
		blockChain = r.Get("blockChain").Int()
		if _, ok := s.chains[blockChain]; !ok {
			s.Error(c, c.Request.URL.Path, errors.New("blockchain is wrong").Error())
			return
		}
	}

	addr := r.Get("address").String()
	token := r.Get("token").String()
	//txType := r.Get("txType").Int()

	if len(addr) < 1 || len(token) < 1 {
		s.Error(c, c.Request.URL.Path, "params is not null")
	}

	//tron base58进制的地址处理
	if !strings.HasPrefix(addr, "0x") && !strings.HasPrefix(addr, "41") && !strings.HasPrefix(addr, "0x41") {
		base58Addr, err := util.Base58ToAddress(addr)
		if err != nil {
			s.Error(c, c.Request.URL.Path, err.Error())
			return
		}
		addr = base58Addr.Hex()
	}

	//address:hex string
	addressTask := &service.MonitorAddress{BlockChain: blockChain, Address: addr, Token: token, TxType: fmt.Sprintf("%v", 0), Id: time.Now().UnixNano()}
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
