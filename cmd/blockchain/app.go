package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/blockchain/config"
	"github.com/uduncloud/easynode/blockchain/service"
	"log"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./cmd/blockchain/config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(1).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := service.NewHttpHandler(cfg.Cluster, xLog)
	//支持JSON-RPC协议的公链
	root.POST("/:chain/jsonrpc", srv.HandlerReq)

	//自定义或不支持JSON-RPC协议的公链
	root.POST("/:chain/block/hash", srv.GetBlockByHash)
	root.POST("/:chain/block/number", srv.GetBlockByNumber)
	root.POST("/:chain/tx/hash", srv.GetTxByHash)
	root.POST("/:chain/tx/receipts", srv.GetTxReceiptByHash)
	root.POST("/:chain/account/balance", srv.GetBalance)
	root.POST("/:chain/account/tokenBalance", srv.GetTokenBalance)
	root.POST("/:chain/account/nonce", srv.GetNonce)
	root.POST("/:chain/block/latest", srv.GetLatestBlock)
	root.POST("/:chain/tx/sendRawTransaction", srv.SendRawTx)

	//ws 协议
	wsServer := service.NewWsHandler(cfg.Cluster, xLog)
	root.Handle("GET", "/ws", func(ctx *gin.Context) {
		wsServer.Start(ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
