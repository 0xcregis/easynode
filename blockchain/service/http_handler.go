package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	easyKafka "github.com/0xcregis/easynode/common/kafka"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type HttpHandler struct {
	log               *logrus.Entry
	nodeCluster       map[int64][]*config.NodeCluster
	blockChainClients map[int64]blockchain.API
	kafkaClient       *easyKafka.EasyKafka
	kafkaCfg          *config.Kafka
	sendCh            chan []*kafka.Message
}

func (h *HttpHandler) StartKafka(ctx context.Context) {
	go func(ctx context.Context) {
		broker := fmt.Sprintf("%v:%v", h.kafkaCfg.Host, h.kafkaCfg.Port)
		kafkaCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		h.kafkaClient.Write(easyKafka.Config{Brokers: []string{broker}, Topic: h.kafkaCfg.Topic, Partition: h.kafkaCfg.Partition}, h.sendCh, nil, kafkaCtx)
	}(ctx)
}

func NewHttpHandler(cluster map[int64][]*config.NodeCluster, kafkaCfg *config.Kafka, xlog *xlog.XLog) *HttpHandler {
	kafkaClient := easyKafka.NewEasyKafka(xlog)
	sendCh := make(chan []*kafka.Message)
	return &HttpHandler{
		log:               xlog.WithField("model", "httpSrv"),
		nodeCluster:       cluster,
		kafkaCfg:          kafkaCfg,
		kafkaClient:       kafkaClient,
		sendCh:            sendCh,
		blockChainClients: NewApis(cluster, xlog),
	}
}

func (h *HttpHandler) GetBlockByHash(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	hash := gjson.ParseBytes(b).Get("hash").String()
	res, err := h.blockChainClients[blockChainCode].GetBlockByHash(blockChainCode, hash, true)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBlockByNumber(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	number := gjson.ParseBytes(b).Get("number").String()
	res, err := h.blockChainClients[blockChainCode].GetBlockByNumber(blockChainCode, number, true)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetTxByHash(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	hash := gjson.ParseBytes(b).Get("hash").String()
	res, err := h.blockChainClients[blockChainCode].GetTxByHash(blockChainCode, hash)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetTxReceiptByHash(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	hash := gjson.ParseBytes(b).Get("hash").String()
	res, err := h.blockChainClients[blockChainCode].GetTransactionReceiptByHash(blockChainCode, hash)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBalance(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	addr := gjson.ParseBytes(b).Get("address").String()
	tag := gjson.ParseBytes(b).Get("tag").String()
	res, err := h.blockChainClients[blockChainCode].Balance(blockChainCode, addr, tag)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

// GetTokenBalance ERC20协议代币余额，后期补充
func (h *HttpHandler) GetTokenBalance(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	r := gjson.ParseBytes(b)
	blockChainCode := r.Get("chain").Int()
	addr := r.Get("address").String()
	codeHash := r.Get("codeHash").String()
	abi := r.Get("abi").String()

	res, err := h.blockChainClients[blockChainCode].TokenBalance(blockChainCode, codeHash, addr, abi)
	if err != nil {
		h.Error(ctx, r.String(), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, r.String(), res, ctx.Request.RequestURI)
}

// GetNonce todo 仅适用于 ether,tron 暂不支持
func (h *HttpHandler) GetNonce(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	addr := gjson.ParseBytes(b).Get("address").String()
	tag := gjson.ParseBytes(b).Get("tag").String() //pending,latest
	res, err := h.blockChainClients[blockChainCode].Nonce(blockChainCode, addr, tag)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetLatestBlock(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	res, err := h.blockChainClients[blockChainCode].LatestBlock(blockChainCode)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, "", res, ctx.Request.RequestURI)
}

func (h *HttpHandler) SendRawTx(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	backup := make(map[string]any, 5)
	defer func(backup map[string]any) {
		bs, _ := json.Marshal(backup)
		msg := &kafka.Message{Topic: h.kafkaCfg.Topic, Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
		h.sendCh <- []*kafka.Message{msg}
	}(backup)

	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	backup["chainCode"] = blockChainCode
	signedTx := root.Get("signed_tx").String()
	backup["signed"] = signedTx
	backup["id"] = time.Now().UnixMicro()
	from := root.Get("from").String()
	backup["from"] = from
	to := root.Get("to").String()
	backup["to"] = to
	extra := root.Get("extra").String()
	backup["extra"] = extra
	res, err := h.blockChainClients[blockChainCode].SendRawTransaction(blockChainCode, signedTx)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		backup["status"] = 0
		return
	}
	backup["status"] = 1
	backup["response"] = res
	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

// HandlerReq  有用户自定义请求内容，然后直接发送到节点 ，和eth_call 函数无关
func (h *HttpHandler) HandlerReq(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	data := gjson.ParseBytes(b).Get("data").String()
	res, err := h.blockChainClients[blockChainCode].SendJsonRpc(blockChainCode, data)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

const (
	SUCCESS = 0
	FAIL    = 1
)

func (h *HttpHandler) Success(c *gin.Context, req string, resp interface{}, path string) {
	req = strings.Replace(req, "\t", "", -1)
	req = strings.Replace(req, "\n", "", -1)
	h.log.Printf("path=%v,req=%v,resp=%v\n", path, req, resp)
	mp := make(map[string]interface{})
	mp["code"] = SUCCESS
	mp["data"] = resp
	c.JSON(200, mp)
}

func (h *HttpHandler) Error(c *gin.Context, req string, path string, err string) {
	req = strings.Replace(req, "\t", "", -1)
	req = strings.Replace(req, "\n", "", -1)
	h.log.Errorf("path=%v,req=%v,err=%v\n", path, req, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["data"] = err
	c.JSON(200, mp)
}
