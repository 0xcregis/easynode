package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"strings"
	"time"

	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/common/chain"
	easyKafka "github.com/0xcregis/easynode/common/kafka"
	"github.com/0xcregis/easynode/common/util"
	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

type HttpHandler struct {
	log                 *logrus.Entry
	nodeCluster         map[int64][]*config.NodeCluster
	blockChainClients   map[int64]blockchain.API
	nftClients          map[int64]blockchain.NftApi
	exBlockChainClients map[int64]blockchain.ExApi
	kafkaClient         *easyKafka.EasyKafka
	kafkaCfg            *config.Kafka
	sendCh              chan []*kafka.Message
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
		log:                 x,
		nodeCluster:         cluster,
		kafkaCfg:            kafkaCfg,
		kafkaClient:         kafkaClient,
		sendCh:              sendCh,
		blockChainClients:   NewApis(cluster, xlog),
		nftClients:          NewNftApis(cluster, xlog),
		exBlockChainClients: NewExApi(cluster, xlog),
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
	res, err := h.blockChainClients[blockChainCode].TokenBalance(blockChainCode, addr, codeHash, abi)
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

// GetAccountResource ,for only tron
func (h *HttpHandler) GetAccountResource(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	addr := gjson.ParseBytes(b).Get("address").String()

	if _, ok := h.exBlockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.exBlockChainClients[blockChainCode].GetAccountResourceForTron(blockChainCode, addr)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), gjson.Parse(res).Value(), ctx.Request.RequestURI)
}

func (h *HttpHandler) EstimateGasForTron(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	if _, ok := h.exBlockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	fromAddr := gjson.ParseBytes(b).Get("from").String()
	toAddr := gjson.ParseBytes(b).Get("to").String()
	functionSelector := gjson.ParseBytes(b).Get("functionSelector").String()
	parameter := gjson.ParseBytes(b).Get("parameter").String()
	res, err := h.exBlockChainClients[blockChainCode].EstimateGasForTron(blockChainCode, fromAddr, toAddr, functionSelector, parameter)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	root := gjson.Parse(res)
	if !root.Get("result.result").Bool() {
		h.Error(ctx, string(b), ctx.Request.RequestURI, res)
		return
	}

	energy := root.Get("energy_required").Int()

	h.Success(ctx, string(b), energy, ctx.Request.RequestURI)
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
		msg := &kafka.Message{Topic: fmt.Sprintf("%v-%v", h.kafkaCfg.Topic, chainCode), Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
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

func (h *HttpHandler) SendRawTx1(ctx *gin.Context) {
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
		msg := &kafka.Message{Topic: fmt.Sprintf("%v-%v", h.kafkaCfg.Topic, chainCode), Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
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

	m := make(map[string]any)
	//{
	//  "id": 1,
	//  "jsonrpc": "2.0",
	//  "result": "0x2a7e11bcb80ea248e09975c48da02b7d0c29d42521d6e9e65e112358132134"
	//}
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		hash := gjson.Parse(res).Get("result").String()
		m["status"] = "SUCCESS"
		m["hash"] = hash
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		//{
		//  "result": {
		//    "code": "SUCCESS",
		//    "message": "SUCCESS",
		//    "transaction": {
		//      "txID": "a6a4c7d5b7f2c3871f4a8271f8b690c324f8c350d0e4b9abaddfe7c2294a1737",
		//      "raw_data": {
		//        // 其他原始数据
		//      },
		//      "signature": [
		//        // 签名信息
		//      ],
		//      // 其他交易信息
		//    }
		//  }
		//}

		//{"code":"SIGERROR","txid":"77ddfa7093cc5f745c0d3a54abb89ef070f983343c05e0f89e5a52f3e5401299","message":"56616c6964617465207369676e6174757265206572726f723a206d69737320736967206f7220636f6e7472616374"}

		root := gjson.Parse(res)
		if root.Get("result.code").Exists() {
			status := root.Get("result.code").String()
			if status == "SUCCESS" {
				m["status"] = status
				id := root.Get("result.transaction.txID").String()
				m["hash"] = id
			} else {
				h.Error(ctx, string(b), ctx.Request.RequestURI, res)
				return
			}

		} else {
			h.Error(ctx, string(b), ctx.Request.RequestURI, res)
			return
		}

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), m, ctx.Request.RequestURI)
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

func (h *HttpHandler) GetBlockByHash1(ctx *gin.Context) {
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

	m := make(map[string]any)
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("hash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("number").String()
		m["blockNumber"], _ = util.HexToInt(blockNumber)
		timestamp := root.Get("timestamp").String()
		timestamp, _ = util.HexToInt(timestamp)
		timestamp = fmt.Sprintf("%v000", timestamp)
		m["timestamp"] = timestamp

	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {

		root := gjson.Parse(res)
		id := root.Get("blockID").String()
		m["blockHash"] = id
		blockNumber := root.Get("block_header.raw_data.number").String()
		m["blockNumber"] = blockNumber
		timestamp := root.Get("block_header.raw_data.timestamp").String()
		m["timestamp"] = timestamp

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), m, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBlockByNumber1(ctx *gin.Context) {
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

	m := make(map[string]any)
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("hash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("number").String()
		m["blockNumber"], _ = util.HexToInt(blockNumber)
		timestamp := root.Get("timestamp").String()
		timestamp, _ = util.HexToInt(timestamp)
		timestamp = fmt.Sprintf("%v000", timestamp)
		m["timestamp"] = timestamp

	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {

		root := gjson.Parse(res)
		id := root.Get("blockID").String()
		m["blockHash"] = id
		blockNumber := root.Get("block_header.raw_data.number").String()
		m["blockNumber"] = blockNumber
		timestamp := root.Get("block_header.raw_data.timestamp").String()
		m["timestamp"] = timestamp

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), m, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetTxByHash1(ctx *gin.Context) {
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

	m := make(map[string]any)
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("blockHash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("blockNumber").String()
		m["blockNumber"], _ = util.HexToInt(blockNumber)
		transactionHash := root.Get("transactionHash").String()
		m["txHash"] = transactionHash
		from := root.Get("from").String()
		m["from"] = from
		to := root.Get("to").String()
		m["to"] = to

		status := root.Get("status").String()
		if status == "0x1" {
			m["status"] = 1
		} else {
			m["status"] = 0
		}

		gasUsed := root.Get("gasUsed").String()

		txResp, err := h.blockChainClients[blockChainCode].GetTxByHash(blockChainCode, hash)
		if err == nil && len(txResp) > 1 {
			gasPrice := gjson.Parse(txResp).Get("result.gasPrice").String()
			bigPrice, b := new(big.Int).SetString(gasPrice, 0)
			bigGas, b2 := new(big.Int).SetString(gasUsed, 0)
			if b && b2 {
				fee := bigPrice.Mul(bigPrice, bigGas)
				m["fee"] = fee.String()
			}
		}

		blockResp, err := h.blockChainClients[blockChainCode].GetBlockByNumber(blockChainCode, blockNumber, false)
		if err == nil && len(blockResp) > 1 {
			timestamp := gjson.Parse(blockResp).Get("result.timestamp").String()
			blockTime, err := util.HexToInt(timestamp)
			if err == nil {
				m["blockTime"] = fmt.Sprintf("%v000", blockTime)
			}
		}

		m["energyUsageTotal"] = "0"

	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {

		root := gjson.Parse(res)
		id := root.Get("id").String()
		m["txHash"] = id
		blockNumber := root.Get("blockNumber").String()
		m["blockNumber"] = blockNumber

		t1 := root.Get("blockTimeStamp").String()
		m["blockTime"] = t1

		fee := root.Get("fee").String()
		m["fee"] = fee

		energyUsageTotal := root.Get("receipt.energy_usage_total").String()
		if len(energyUsageTotal) < 1 {
			energyUsageTotal = "0"
		}
		m["energyUsageTotal"] = energyUsageTotal

		status := root.Get("receipt.result").String()
		if status == "SUCCESS" {
			m["status"] = 1
		} else {
			m["status"] = 0
		}

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), m, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBalance1(ctx *gin.Context) {
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

	r := make(map[string]any)
	r["utxo"] = ""
	r["address"] = addr
	var nonce string
	//{"jsonrpc":"2.0","id":1,"result":"0x20ef7755b4d96faf5"}
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		balance := gjson.Parse(res).Get("result").String()
		balance, _ = util.HexToInt(balance)
		r["balance"] = balance
		resNonce, err := h.blockChainClients[blockChainCode].Nonce(blockChainCode, addr, "latest")
		if err == nil && len(resNonce) > 1 {
			nonce = gjson.Parse(resNonce).Get("result").String()
			nonce, _ = util.HexToInt(nonce)
		}
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		// {"balance":309739}
		balance := gjson.Parse(res).Get("balance").String()
		r["balance"] = balance
		nonce = "0"
	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	//r["totalAmount"] = r["balance"]
	r["nonce"] = nonce

	h.Success(ctx, string(b), r, ctx.Request.RequestURI)
}

// GetTokenBalance1  ERC20协议代币余额，后期补充
func (h *HttpHandler) GetTokenBalance1(ctx *gin.Context) {
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
	res, err := h.blockChainClients[blockChainCode].TokenBalance(blockChainCode, addr, codeHash, abi)
	if err != nil {
		h.Error(ctx, r.String(), ctx.Request.RequestURI, err.Error())
		return
	}

	m := make(map[string]any)
	m["utxo"] = ""
	m["address"] = addr
	var nonce string

	//  {"balance":"1233764293093","decimals":6,"name":"Tether USD","symbol":"USDT"}
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		balance := gjson.Parse(res).Get("result").String()
		balance, _ = util.HexToInt(balance)
		m["balance"] = balance
		resNonce, err := h.blockChainClients[blockChainCode].Nonce(blockChainCode, addr, "latest")
		if err == nil && len(resNonce) > 1 {
			nonce = gjson.Parse(resNonce).Get("result").String()
			nonce, _ = util.HexToInt(nonce)
		}
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		//{"balance":"74305000412789","decimals":"6"}
		balance := gjson.Parse(res).Get("balance").String()
		m["balance"] = balance
		nonce = "0"
	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	//m["totalAmount"] = m["balance"]
	m["nonce"] = nonce

	h.Success(ctx, r.String(), m, ctx.Request.RequestURI)
}

// GetNonce1 todo 仅适用于 ether,tron 暂不支持
func (h *HttpHandler) GetNonce1(ctx *gin.Context) {
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

	var nonce string

	//{"jsonrpc":"2.0","id":1,"result":"0x1a48c9"}
	if chain.GetChainCode(blockChainCode, "ETH", nil) {
		nonce = gjson.Parse(res).Get("result").String()
		nonce, _ = util.HexToInt(nonce)
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		nonce = "0"
	} else if chain.GetChainCode(blockChainCode, "BSC", nil) {
		nonce = gjson.Parse(res).Get("result").String()
		nonce, _ = util.HexToInt(nonce)
	} else if chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		nonce = gjson.Parse(res).Get("result").String()
		nonce, _ = util.HexToInt(nonce)
	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), nonce, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetLatestBlock1(ctx *gin.Context) {
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

	var number string
	if chain.GetChainCode(blockChainCode, "ETH", nil) {
		//{"jsonrpc":"2.0","result":"0x11d9e01","id":1}
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		//{"blockId":"00000000036660c483762cf17eb14ddb705228e0f616fb973a3ee3b7a9062910","number":57041092}
		number = gjson.Parse(res).Get("number").String()

	} else if chain.GetChainCode(blockChainCode, "BSC", nil) {
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, "", number, ctx.Request.RequestURI)
}

// GasPrice1 Returns the current price per gas in wei.
func (h *HttpHandler) GasPrice1(ctx *gin.Context) {
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
	//{
	//    "jsonrpc": "2.0",
	//    "id": 1,
	//    "result": "0xadc3e16bc"
	//}
	gas := gjson.Parse(res).Get("result").String()
	gas, err = util.HexToInt(gas)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), gas, ctx.Request.RequestURI)
}

// EstimateGas1 Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain.
func (h *HttpHandler) EstimateGas1(ctx *gin.Context) {
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

	//{
	//    "jsonrpc": "2.0",
	//    "id": 1,
	//    "result": "0x5208"
	//}
	gas := gjson.Parse(res).Get("result").String()

	gas, err = util.HexToInt(gas)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), gas, ctx.Request.RequestURI)
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
	mp["message"] = "ok"
	mp["data"] = resp
	c.JSON(200, mp)
}

func (h *HttpHandler) Error(c *gin.Context, req string, path string, err string) {
	req = strings.Replace(req, "\t", "", -1)
	req = strings.Replace(req, "\n", "", -1)
	h.log.Errorf("path=%v,req=%v,err=%v\n", path, req, err)
	mp := make(map[string]interface{})
	mp["code"] = FAIL
	mp["message"] = err
	mp["data"] = ""
	c.JSON(200, mp)
}
