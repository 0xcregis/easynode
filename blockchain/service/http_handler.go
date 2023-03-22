package service

import (
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/blockchain/config"
	"io/ioutil"
	"strconv"
	"strings"
)

type HttpHandler struct {
	log               *xlog.XLog
	nodeCluster       map[int64][]*config.NodeCluster
	blockChainClients map[int64]API
}

func NewHttpHandler(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) *HttpHandler {

	blockChainClients := make(map[int64]API, 0)

	for k, _ := range cluster {
		if k == 200 {
			bc := NewEth(cluster, xlog)
			blockChainClients[k] = bc
		}
		if k == 205 {
			bc := NewTron(cluster, xlog)
			blockChainClients[k] = bc
		}
	}

	return &HttpHandler{
		log:               xlog,
		nodeCluster:       cluster,
		blockChainClients: blockChainClients,
	}
}

func (h *HttpHandler) GetBlockByHash(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	hash := gjson.ParseBytes(b).Get("hash").String()
	res, err := h.blockChainClients[blockChainCode].GetBlockByHash(blockChainCode, hash)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBlockByNumber(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	number := gjson.ParseBytes(b).Get("number").String()
	res, err := h.blockChainClients[blockChainCode].GetBlockByNumber(blockChainCode, number)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetTxByHash(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	hash := gjson.ParseBytes(b).Get("hash").String()
	res, err := h.blockChainClients[blockChainCode].GetTxByHash(blockChainCode, hash)

	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

func (h *HttpHandler) GetBalance(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

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
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	r := gjson.ParseBytes(b)
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
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

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
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	res, err := h.blockChainClients[blockChainCode].LatestBlock(blockChainCode)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	h.Success(ctx, "", res, ctx.Request.RequestURI)
}

func (h *HttpHandler) SendRawTx(ctx *gin.Context) {
	code := ctx.Param("chain")
	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	signedTx := gjson.ParseBytes(b).Get("signed_tx").String()
	res, err := h.blockChainClients[blockChainCode].SendRawTransaction(blockChainCode, signedTx)
	if err != nil {
		h.Error(ctx, string(b), ctx.Request.RequestURI, err.Error())
		return
	}
	h.Success(ctx, string(b), res, ctx.Request.RequestURI)
}

// HandlerReq  有用户自定义请求内容，然后直接发送到节点 ，和eth_call 函数无关
func (h *HttpHandler) HandlerReq(ctx *gin.Context) {
	code := ctx.Param("chain")

	blockChainCode, err := strconv.ParseInt(code, 0, 64)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}
	b, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.Error(ctx, "", ctx.Request.RequestURI, err.Error())
		return
	}

	res, err := h.blockChainClients[blockChainCode].SendJsonRpc(blockChainCode, string(b))
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
