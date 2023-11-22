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
	nftClients        map[int64]blockchain.NftApi
	kafkaClient       *easyKafka.EasyKafka
	kafkaCfg          *config.Kafka
	sendCh            chan []*kafka.Message
}

func (h *HttpHandler) StartKafka(ctx context.Context) {
	go func(ctx context.Context) {
		broker := fmt.Sprintf("%v:%v", h.kafkaCfg.Host, h.kafkaCfg.Port)
		kafkaCtx, cancel := context.WithCancel(ctx)
		defer cancel()
		h.kafkaClient.Write(easyKafka.Config{Brokers: []string{broker}}, h.sendCh, nil, kafkaCtx)
	}(ctx)
}

func NewHttpHandler(cluster map[int64][]*config.NodeCluster, kafkaCfg *config.Kafka, xlog *xlog.XLog) *HttpHandler {
	x := xlog.WithField("model", "httpSrv")
	kafkaClient := easyKafka.NewEasyKafka2(x)
	sendCh := make(chan []*kafka.Message)
	return &HttpHandler{
		log:               x,
		nodeCluster:       cluster,
		kafkaCfg:          kafkaCfg,
		kafkaClient:       kafkaClient,
		sendCh:            sendCh,
		blockChainClients: NewApis(cluster, xlog),
		nftClients:        NewNftApis(cluster, xlog),
	}
}

func (h *HttpHandler) TokenUri(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	contract := root.Get("contract").String()
	tokenId := root.Get("tokenId").String()
	eip := root.Get("eip").Int()
	if _, ok := h.nftClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.nftClients[blockChainCode].TokenURI(blockChainCode, contract, tokenId, eip)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) BalanceOf(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	contract := root.Get("contract").String()
	addr := root.Get("address").String()
	tokenId := root.Get("tokenId").String()
	eip := root.Get("eip").Int()
	if _, ok := h.nftClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.nftClients[blockChainCode].BalanceOf(blockChainCode, contract, addr, tokenId, eip)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) OwnerOf(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	contract := root.Get("contract").String()
	tokenId := root.Get("tokenId").String()
	eip := root.Get("eip").Int()
	if _, ok := h.nftClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.nftClients[blockChainCode].OwnerOf(blockChainCode, contract, tokenId, eip)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) TotalSupply(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	contract := root.Get("contract").String()
	eip := root.Get("eip").Int()
	if _, ok := h.nftClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.nftClients[blockChainCode].TotalSupply(blockChainCode, contract, eip)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}
	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBlockByHash(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	hash := gjson.ParseBytes(b).Get("hash").String()

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	codeHash := r.Get("contract").String()
	abi := r.Get("abi").String()

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.blockChainClients[blockChainCode].LatestBlock(blockChainCode)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, "", res, ctx.Request.RequestURI)
}

// GasPrice Returns the current price per gas in wei.
func (h *HttpHandler) GasPrice(ctx *gin.Context) {
	req := `
		{
		"id": 1,
		"jsonrpc": "2.0",
		"method": "eth_gasPrice"
		}
		`

	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.blockChainClients[blockChainCode].SendJsonRpc(blockChainCode, req)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

// EstimateGas Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain.
func (h *HttpHandler) EstimateGas(ctx *gin.Context) {
	req := `
	{
		  "id": 1,
		  "jsonrpc": "2.0",
		  "method": "eth_estimateGas",
		  "params": [
			{
			  "from":"%v",
			  "to": "%v",
			  "data": "%v"
			}
		  ]
	}`

	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	root := gjson.ParseBytes(b)
	blockChainCode := root.Get("chain").Int()
	from := root.Get("from").String()
	to := root.Get("to").String()
	data := root.Get("data").String()
	if len(data) < 1 {
		data = "0x"
	}

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	req = fmt.Sprintf(req, from, to, data)
	res, err := h.blockChainClients[blockChainCode].SendJsonRpc(blockChainCode, req)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) SendRawTx(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	backup := make(map[string]any, 5)
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

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	defer func(backup map[string]any, chainCode int64) {
		bs, _ := json.Marshal(backup)
		msg := &kafka.Message{Topic: fmt.Sprintf("%v_%v", h.kafkaCfg.Topic, chainCode), Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
		h.sendCh <- []*kafka.Message{msg}
	}(backup, blockChainCode)

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
	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
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
	if v, ok := resp.(string); ok {
		resp = strings.Replace(v, "\n", "", -1)
	}
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
