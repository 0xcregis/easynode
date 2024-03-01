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
	"github.com/0xcregis/easynode/common/ethtypes"
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
		if h.kafkaCfg != nil {
			broker := fmt.Sprintf("%v:%v", h.kafkaCfg.Host, h.kafkaCfg.Port)
			kafkaCtx, cancel := context.WithCancel(ctx)
			defer cancel()
			h.kafkaClient.Write(easyKafka.Config{Brokers: []string{broker}}, h.sendCh, nil, kafkaCtx)
		}
	}(ctx)
}

func NewHttpHandler(cluster map[int64][]*config.NodeCluster, kafkaCfg *config.Kafka, xlog *xlog.XLog) *HttpHandler {
	x := xlog.WithField("model", "httpSrv")

	var kafkaClient *easyKafka.EasyKafka
	sendCh := make(chan []*kafka.Message)
	if kafkaCfg != nil {
		kafkaClient = easyKafka.NewEasyKafka2(x)
	}
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

func (h *HttpHandler) GetTraceTransaction(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	blockChainCode := gjson.ParseBytes(b).Get("chain").Int()
	hash := gjson.ParseBytes(b).Get("hash").String()
	if _, ok := h.exBlockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.exBlockChainClients[blockChainCode].TraceTransaction(blockChainCode, hash)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	r := make([]*blockchain.InterTx, 0, 10)
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		list := gjson.Parse(res).Get("result").Array()
		for _, v := range list {
			value := v.Get("action.value").String()
			if value != "0x0" {
				var t blockchain.InterTx
				_ = json.Unmarshal([]byte(v.String()), &t)
				hexValue, _ := util.HexToInt(t.Action.Value)
				t.Action.Value = util.Div(hexValue, 18)
				t.Action.Gas, _ = util.HexToInt(t.Action.Gas)
				r = append(r, &t)
			}
		}

	}

	h.Success(ctx, string(b), r, ctx.Request.RequestURI)
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

// GetToken get token info
func (h *HttpHandler) GetToken(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	r := gjson.ParseBytes(b)
	blockChainCode := r.Get("chain").Int()
	codeHash := r.Get("contract").String()
	abi := r.Get("abi").String()
	eip := r.Get("eip").String()

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.blockChainClients[blockChainCode].Token(blockChainCode, codeHash, abi, eip)
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

	r := make(map[string]any, 4)
	gjson.Parse(res).ForEach(func(key, value gjson.Result) bool {
		k := key.String()
		k = strings.ToLower(k[:1]) + k[1:]
		r[k] = value.Value()
		return true
	})

	h.Success(ctx, string(b), r, ctx.Request.RequestURI)
}

// EstimateGasForTron ,for only tron
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

	energy := root.Get("energy_used").Int()

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
	res, err := h.exBlockChainClients[blockChainCode].GasPrice(blockChainCode)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

// EstimateGas Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain.
func (h *HttpHandler) EstimateGas(ctx *gin.Context) {
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

	if _, ok := h.exBlockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	res, err := h.exBlockChainClients[blockChainCode].EstimateGas(blockChainCode, from, to, data)
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
		if h.kafkaClient != nil {
			bs, _ := json.Marshal(backup)
			msg := &kafka.Message{Topic: fmt.Sprintf("%v-%v", h.kafkaCfg.Topic, chainCode), Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
			h.sendCh <- []*kafka.Message{msg}
		}
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
		if h.kafkaClient != nil {
			bs, _ := json.Marshal(backup)
			msg := &kafka.Message{Topic: fmt.Sprintf("%v-%v", h.kafkaCfg.Topic, chainCode), Partition: h.kafkaCfg.Partition, Key: []byte(fmt.Sprintf("%v", time.Now().UnixNano())), Value: bs}
			h.sendCh <- []*kafka.Message{msg}
		}
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
	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		//res = `{
		// "id": 1,
		// "jsonrpc": "2.0",
		// "result": "0x2a7e11bcb80ea248e09975c48da02b7d0c29d42521d6e9e65e112358132134"
		//}`
		hash := gjson.Parse(res).Get("result").String()
		m["status"] = "SUCCESS"
		m["hash"] = hash
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {

		//		res := `
		//{
		//  "result": true,
		//  "code": "SUCCESS",
		//  "txid": "ce5f03b89a735353bebf7214d24cac222d35f98c92d690ed3d3bc4f81450654b",
		//  "message": "",
		//  "transaction": {
		//    "raw_data": {
		//      "ref_block_bytes": "7fc6",
		//      "ref_block_hash": "11d6bde1afbed1f3",
		//      "expiration": 1702633179923,
		//      "contract": [
		//        {
		//          "type": "TransferContract",
		//          "parameter": {
		//            "type_url": "type.googleapis.com/protocol.TransferContract",
		//            "value": "0a1541cb3725cf4bb8acbe2758bf49da12dfc9b522458f12154184963abb43debc2d129d1905eae92e95a4aa0a531880897a"
		//          }
		//        }
		//      ],
		//      "timestamp": 1702632879923
		//    },
		//    "signature": [
		//      "91fe34c5142236098167f726da602a269393a8277cfdffa467b5c8ee659f887e0341fb012b1390af2d3e74900d8b07cf25081471cd0669206740d86d06d7d2d300"
		//    ]
		//  }
		//}
		//`

		//{"code":"SIGERROR","txid":"77ddfa7093cc5f745c0d3a54abb89ef070f983343c05e0f89e5a52f3e5401299","message":"56616c6964617465207369676e6174757265206572726f723a206d69737320736967206f7220636f6e7472616374"}

		root := gjson.Parse(res)
		if root.Get("code").Exists() {
			status := root.Get("code").String()
			if status == "SUCCESS" {
				m["status"] = status
				id := root.Get("txid").String()
				m["hash"] = id
			} else {
				h.Error(ctx, string(b), ctx.Request.RequestURI, res)
				return
			}

		} else {
			h.Error(ctx, string(b), ctx.Request.RequestURI, res)
			return
		}

	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {
		/**
		{
		  "/": "bafy2bzacea3wsdh6y3a36tb3skempjoxqpuyompjbmfeyf34fi3uy6uue42v4"
		}
		*/
		hash := gjson.Parse(res).Get("/").String()
		m["status"] = "SUCCESS"
		m["hash"] = hash
	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {
		/**
		{
		    "result": {
		        "accepted": true,
		        "account_sequence_available": 393,
		        "account_sequence_next": 393,
		        "applied": false,
		        "broadcast": false,
		        "engine_result": "tefPAST_SEQ",
		        "engine_result_code": -190,
		        "engine_result_message": "This sequence number has already passed.",
		        "kept": true,
		        "open_ledger_cost": "10",
		        "queued": false,
		        "status": "success",
		        "tx_blob": "1200002280000000240000000361D4838D7EA4C6800000000000000000000000000055534400000000004B4E9C06F24296074F7BC48F92A97916C6DC5EA968400000000000000A732103AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB74473045022100D184EB4AE5956FF600E7536EE459345C7BBCF097A84CC61A93B9AF7197EDB98702201CEA8009B7BEEBAA2AACC0359B41C427C1C5B550A4CA4B80CF2174AF2D6D5DCE81144B4E9C06F24296074F7BC48F92A97916C6DC5EA983143E9D4A2B8AA0780F682D136F7A56D6724EF53754",
		        "tx_json": {
		            "Account": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
		            "Amount": {
		                "currency": "USD",
		                "issuer": "rf1BiGeXwwQoi8Z2ueFYTEXSwuJYfV2Jpn",
		                "value": "1"
		            },
		            "Destination": "ra5nK24KXen9AHvsdFTKHSANinZseWnPcX",
		            "Fee": "10",
		            "Flags": 2147483648,
		            "Sequence": 3,
		            "SigningPubKey": "03AB40A0490F9B7ED8DF29D246BF2D6269820A0EE7742ACDD457BEA7C7D0931EDB",
		            "TransactionType": "Payment",
		            "TxnSignature": "3045022100D184EB4AE5956FF600E7536EE459345C7BBCF097A84CC61A93B9AF7197EDB98702201CEA8009B7BEEBAA2AACC0359B41C427C1C5B550A4CA4B80CF2174AF2D6D5DCE",
		            "hash": "82230B9D489370504B39BC2CE46216176CAC9E752E5C1774A8CBEC9FBB819208"
		        },
		        "validated_ledger_index": 85447184
		    }
		}

		*/

		root := gjson.Parse(res)
		status := root.Get("result.status").String()
		if status == "success" {
			m["status"] = "SUCCESS"
		} else {
			m["status"] = "FAILURE"
		}

		hash := root.Get("result.tx_json.hash").String()
		m["hash"] = hash
	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), m, ctx.Request.RequestURI)
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

	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

		root := gjson.Parse(res)
		blockHash := root.Get("Cid").String()
		m["blockHash"] = blockHash
		head := root.Get("BlockHead").String()
		timestamp := gjson.Parse(head).Get("Timestamp").String()
		timestamp = fmt.Sprintf("%v000", timestamp)
		m["timestamp"] = timestamp
		Height := gjson.Parse(head).Get("Height").String()
		m["blockNumber"] = Height

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {
		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("ledger_hash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("ledger_index").String()
		m["blockNumber"] = blockNumber
		timestamp := root.Get("ledger.close_time").String()
		timestamp = fmt.Sprintf("%v000", timestamp)
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

	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

		m["blockNumber"] = number
		list := gjson.Parse(res).Array()
		blocks := make([]any, 0, 2)
		for _, v := range list {
			mp := make(map[string]any, 2)
			blockHash := v.Get("Cid").String()
			mp["blockHash"] = blockHash
			head := v.Get("BlockHead").String()
			timestamp := gjson.Parse(head).Get("Timestamp").String()
			timestamp = fmt.Sprintf("%v000", timestamp)
			mp["timestamp"] = timestamp
			blocks = append(blocks, mp)
		}
		m["blocks"] = blocks

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {

		/**
		{
		    "result":{
		        "ledger":{
		            "account_hash":"B3D7C8884773E0EA71EFD7637188DB0E284D1F17ABB15CA93B830F2B51FC3511",
		            "close_flags":0,
		            "close_time":758969601,
		            "close_time_human":"2024-Jan-19 08:53:21.000000000 UTC",
		            "close_time_iso":"2024-01-19T08:53:21Z",
		            "close_time_resolution":10,
		            "closed":true,
		            "ledger_hash":"8AC92BF3FC88166143F5039154B900288189682E2FA28EA1373FB497BAC155A5",
		            "ledger_index":"85379901",
		            "parent_close_time":758969600,
		            "parent_hash":"B60522E77C4ECCD0E936E091A020C813B8D32A159CDE52E5EB772777025CA98F",
		            "total_coins":"99987964714422794",
		            "transaction_hash":"A705F5EAB7DD7C610C239F86F16C6E721CA8C2422783A7329911B72E36C17000",
		            "transactions":Array[170]
		        },
		        "ledger_hash":"8AC92BF3FC88166143F5039154B900288189682E2FA28EA1373FB497BAC155A5",
		        "ledger_index":85379901,
		        "status":"success",
		        "validated":true
		    }
		}
		*/

		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("ledger_hash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("ledger_index").String()
		m["blockNumber"] = blockNumber
		timestamp := root.Get("ledger.close_time").String()
		timestamp = fmt.Sprintf("%v000", timestamp)
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

		if len(root.Get("contract_address").String()) > 5 {
			// contract tx
			status := root.Get("receipt.result").String()
			if status == "SUCCESS" {
				m["status"] = 1
			} else {
				m["status"] = 0
			}
		} else {
			// normal tx
			m["status"] = 1
		}

	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

		root := gjson.Parse(res).Get("result")
		blockHash := root.Get("blockHash").String()
		eth, _ := ethtypes.ParseEthHash(blockHash)
		m["blockHash"] = eth.ToCid().String()
		blockNumber := root.Get("blockNumber").String()
		m["blockNumber"], _ = util.HexToInt(blockNumber)
		transactionHash := root.Get("transactionHash").String()
		eth, _ = ethtypes.ParseEthHash(transactionHash)
		txHash := eth.ToCid().String()
		m["txHash"] = txHash
		from := root.Get("from").String()
		addr, _ := ethtypes.ParseEthAddress(from)
		fAddr, _ := addr.ToFilecoinAddress()
		m["from"] = fAddr.String()
		to := root.Get("to").String()
		addr, _ = ethtypes.ParseEthAddress(to)
		fAddr, _ = addr.ToFilecoinAddress()
		m["to"] = fAddr.String()

		status := root.Get("status").String()
		if status == "0x1" {
			m["status"] = 1
		} else {
			m["status"] = 0
		}

		gasUsed := root.Get("gasUsed").String()
		gasPrice := root.Get("effectiveGasPrice").String()
		bigPrice, b := new(big.Int).SetString(gasPrice, 0)
		bigGas, b2 := new(big.Int).SetString(gasUsed, 0)
		if b && b2 {
			fee := bigPrice.Mul(bigPrice, bigGas)
			m["fee"] = fee.String()
		}

		blockResp, err := h.blockChainClients[blockChainCode].GetBlockByHash(blockChainCode, hash, false)
		if err == nil && len(blockResp) > 1 {
			root := gjson.Parse(blockResp)
			head := root.Get("BlockHead").String()
			timestamp := gjson.Parse(head).Get("Timestamp").String()
			timestamp = fmt.Sprintf("%v000", timestamp)
			m["blockTime"] = timestamp
			Height := gjson.Parse(head).Get("Height").String()
			m["blockNumber"] = Height
		}

		m["energyUsageTotal"] = "0"

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {
		/**

		{
		    "account":"rMZH9Lqv3ASSBAnjLFEb9PRDGbbXRqpZsj",
		    "date":758969601,
		    "destination":"rxRpSNb1VktvzBz8JF2oJC6qaww6RZ7Lw",
		    "fee":12,
		    "hash":"A7E17BB97A75A2FEF8F47179EF9329F669C0A31019F950CBEF504277444A8C40",
		    "ledgerIndex":85379901,
		    "status":"success",
		    "transactionIndex":0,
		    "transactionResult":"tesSUCCESS"
		}
		*/
		root := gjson.Parse(res)
		blockHash := root.Get("blockHash").String()
		m["blockHash"] = blockHash
		blockNumber := root.Get("ledgerIndex").String()
		m["blockNumber"], _ = util.HexToInt(blockNumber)
		transactionHash := root.Get("hash").String()
		m["txHash"] = transactionHash
		from := root.Get("account").String()
		m["from"] = from
		to := root.Get("destination").String()
		m["to"] = to

		status := root.Get("status").String()
		if status == "success" {
			m["status"] = 1
		} else {
			m["status"] = 0
		}

		fee := root.Get("fee").String()
		m["fee"] = fee
		blockTime := root.Get("date").String()
		m["blockTime"] = fmt.Sprintf("%v000", blockTime)
		m["energyUsageTotal"] = "0"
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
	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {
		//\"jsonrpc\":\"2.0\",\"result\":\"1050000000000000000\",\"id\":1}
		balance := gjson.Parse(res).Get("result").String()
		balance, _ = util.HexToInt(balance)
		r["balance"] = balance
		resNonce, err := h.blockChainClients[blockChainCode].Nonce(blockChainCode, addr, "latest")
		if err == nil && len(resNonce) > 1 {
			nonce = gjson.Parse(resNonce).Get("result").String()
			nonce, _ = util.HexToInt(nonce)
		}

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {

		/**
		{\"Account\":\"rKag74JVKyscABvLKLEvpQifHvuPVcx8xY\",\"Balance\":\"10598167\",\"Flags\":0,\"LedgerEntryType\":\"AccountRoot\",\"OwnerCount\":0,\"PreviousTxnID\":\"EE68D4341564502454166EF553D295E245375CFC948851CF5161AD1F81AA7ABE\",\"PreviousTxnLgrSeq\":85455870,\"Sequence\":84842028,\"index\":\"A78EFB04464466F56EC538BC5F4BCE4F928FD55EB4E0B62DB0FA51EAECC10CD0\"}
		*/
		balance := gjson.Parse(res).Get("Balance").String()
		r["balance"] = balance

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	//r["totalAmount"] = r["balance"]
	r["nonce"] = nonce

	h.Success(ctx, string(b), r, ctx.Request.RequestURI)
}

// GetToken1 get token info
func (h *HttpHandler) GetToken1(ctx *gin.Context) {
	b, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	r := gjson.ParseBytes(b)
	blockChainCode := r.Get("chain").Int()
	codeHash := r.Get("contract").String()
	abi := r.Get("abi").String()
	eip := r.Get("eip").String()

	if _, ok := h.blockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}
	res, err := h.blockChainClients[blockChainCode].Token(blockChainCode, codeHash, abi, eip)
	if err != nil {
		h.Error(ctx, r.String(), ctx.Request.RequestURI, err.Error())
		return
	}
	h.Success(ctx, r.String(), gjson.Parse(res).Value(), ctx.Request.RequestURI)
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
		balance := gjson.Parse(res).Get("balance").String()
		//balance, _ = util.HexToInt(balance)
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
	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {

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
	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {
		//{\"jsonrpc\":\"2.0\",\"result\":38817,\"id\":1}
		nonce = gjson.Parse(res).Get("result").String()
		nonce, _ = util.HexToInt(nonce)
	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {
		nonce = "0"
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

	var number, blockHash string
	if chain.GetChainCode(blockChainCode, "ETH", nil) {
		//{"jsonrpc":"2.0","result":"0x11d9e01","id":1}
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		//{"blockId":"00000000036660c483762cf17eb14ddb705228e0f616fb973a3ee3b7a9062910","number":57041092}
		root := gjson.Parse(res)
		number = root.Get("number").String()
		blockHash = root.Get("blockId").String()

	} else if chain.GetChainCode(blockChainCode, "BSC", nil) {
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

		//{\"jsonrpc\":\"2.0\",\"result\":\"0x36c62f\",\"id\":1}
		number = gjson.Parse(res).Get("result").String()
		number, err = util.HexToInt(number)

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {
		/**
		{\"result\":{\"ledger_hash\":\"A3B94654D20C6BE048DCBF77E90E90518BE193D45D6CA1C8487CA17BE0C06DD1\",\"ledger_index\":85467362,\"status\":\"success\"}}
		*/
		root := gjson.Parse(res).Get("result")
		blockHash = root.Get("ledger_hash").String()
		number = root.Get("ledger_index").String()

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	r := make(map[string]any, 0)
	r["blockNumber"] = number
	r["blockHash"] = blockHash
	h.Success(ctx, string(b), r, ctx.Request.RequestURI)
}

// GasPrice1 Returns the current price per gas in wei.
func (h *HttpHandler) GasPrice1(ctx *gin.Context) {
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
	res, err := h.exBlockChainClients[blockChainCode].GasPrice(blockChainCode)
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

	if chain.GetChainCode(blockChainCode, "ETH", nil) || chain.GetChainCode(blockChainCode, "BSC", nil) || chain.GetChainCode(blockChainCode, "POLYGON", nil) {
		gas = util.Div(gas, 9) //gwei
	} else if chain.GetChainCode(blockChainCode, "TRON", nil) {
		//sun
	} else if chain.GetChainCode(blockChainCode, "FIL", nil) {

	} else if chain.GetChainCode(blockChainCode, "XRP", nil) {

	} else {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	h.Success(ctx, string(b), gas, ctx.Request.RequestURI)
}

// EstimateGas1 Generates and returns an estimate of how much gas is necessary to allow the transaction to complete. The transaction will not be added to the blockchain.
func (h *HttpHandler) EstimateGas1(ctx *gin.Context) {
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

	if _, ok := h.exBlockChainClients[blockChainCode]; !ok {
		h.Error(ctx, string(b), ctx.Request.RequestURI, fmt.Sprintf("blockchain:%v is not supported", blockChainCode))
		return
	}

	res, err := h.exBlockChainClients[blockChainCode].EstimateGas(blockChainCode, from, to, data)
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
