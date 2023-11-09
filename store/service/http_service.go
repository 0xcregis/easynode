package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/0xcregis/easynode/common/util"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/config"
	"github.com/0xcregis/easynode/store/db"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type HttpHandler struct {
	log    *logrus.Entry
	chains map[int64]*config.Chain
	store  store.DbStoreInterface
}

func NewHttpHandler(cfg *config.Config, log *xlog.XLog) *HttpHandler {
	s := db.NewChService(cfg, log)
	mp := make(map[int64]*config.Chain, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
	}
	return &HttpHandler{
		store:  s,
		log:    log.WithField("model", "httpSrv"),
		chains: mp,
	}
}

func (s *HttpHandler) AddSubFilter(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	var list []*store.SubFilter

	err = json.Unmarshal(bs, &list)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	err = s.store.NewSubFilter(list)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *HttpHandler) DelSubFilter(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	root := gjson.ParseBytes(bs)

	if root.Get("id").Exists() {
		id := root.Get("id").Int()
		err = s.store.DelSubFilter(id)
	} else {
		chainCode := root.Get("blockChain").Int()
		token := root.Get("token").String()
		txCode := root.Get("txCode").String()
		sf := store.SubFilter{BlockChain: chainCode, Token: token, TxCode: txCode}
		err = s.store.DelSubFilter2(&sf)
	}

	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *HttpHandler) QuerySubFilter(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	root := gjson.ParseBytes(bs)

	chainCode := root.Get("blockChain").Int()
	token := root.Get("token").String()
	txCode := root.Get("txCode").String()
	list, err := s.store.GetSubFilter(token, chainCode, txCode)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, list, c.Request.URL.Path)
}

func (s *HttpHandler) NewToken(c *gin.Context) {
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

	err, _ = s.store.GetNodeTokenByEmail(email)

	if err == nil {
		//已存在
		s.Error(c, c.Request.URL.Path, errors.New("the email already has token").Error())
		return
	}

	//保存token
	err = s.store.NewToken(&store.NodeToken{Token: token.String(), Email: email, Id: time.Now().UnixMicro()})
	if err != nil {
		s.Error(c, c.Request.URL.Path, errors.New("create token failure").Error())
		return
	}

	s.Success(c, token.String(), c.Request.URL.Path)
}

func (s *HttpHandler) DelMonitorAddress(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	token := r.Get("token").String()
	address := r.Get("address").String()
	blockchain := r.Get("blockChain").Int()
	err = s.store.DelMonitorAddress(blockchain, token, address)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, nil, c.Request.URL.Path)
}

func (s *HttpHandler) GetMonitorAddress(c *gin.Context) {
	bs, err := io.ReadAll(c.Request.Body)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	r := gjson.ParseBytes(bs)
	token := r.Get("token").String()

	list, err := s.store.GetAddressByToken2(token)
	if err != nil {
		s.Error(c, c.Request.URL.Path, err.Error())
		return
	}

	s.Success(c, list, c.Request.URL.Path)
}

func (s *HttpHandler) MonitorAddress(c *gin.Context) {
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
	if blockChain == 205 && !strings.HasPrefix(addr, "0x") && !strings.HasPrefix(addr, "41") && !strings.HasPrefix(addr, "0x41") {
		base58Addr, err := util.Base58ToAddress(addr)
		if err != nil {
			s.Error(c, c.Request.URL.Path, err.Error())
			return
		}
		addr = base58Addr.Hex()
	}

	//address:hex string
	addressTask := &store.MonitorAddress{BlockChain: blockChain, Address: addr, Token: token, TxType: fmt.Sprintf("%v", 0), Id: time.Now().UnixNano()}
	err = s.store.AddMonitorAddress(blockChain, addressTask)
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

func (s *HttpHandler) Success(c *gin.Context, resp interface{}, path string) {
	s.log.Printf("path=%v,body=%v\n", path, resp)
	mp := make(map[string]interface{})
	mp["code"] = SUCCESS
	mp["data"] = resp
	c.JSON(200, mp)
}

func (s *HttpHandler) Error(c *gin.Context, path string, err string) {
	s.log.Printf("path=%v,err=%v\n", path, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["data"] = err
	c.JSON(200, mp)
}
