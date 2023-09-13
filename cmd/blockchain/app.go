package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/blockchain/service"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "blockchain", "./cmd/blockchain/config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain/chain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := service.NewHttpHandler(cfg.Cluster, cfg.Kafka, xLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//start kafka listener
	srv.StartKafka(ctx)

	//支持JSON-RPC协议的公链
	root.POST("/jsonrpc", srv.HandlerReq)

	//自定义或不支持JSON-RPC协议的公链
	root.POST("/block/hash", srv.GetBlockByHash)
	root.POST("/block/number", srv.GetBlockByNumber)
	root.POST("/tx/hash", srv.GetTxByHash)
	root.POST("/receipts/hash", srv.GetTxReceiptByHash)
	root.POST("/account/balance", srv.GetBalance)
	root.POST("/account/tokenBalance", srv.GetTokenBalance)
	root.POST("/account/nonce", srv.GetNonce)
	root.POST("/block/latest", srv.GetLatestBlock)
	root.POST("/tx/sendRawTransaction", srv.SendRawTx)
	root.POST("/nft/tokenUri", srv.TokenUri)
	root.POST("/nft/balanceOf", srv.BalanceOf)
	root.POST("/nft/owner", srv.OwnerOf)
	root.POST("/nft/totalSupply", srv.TotalSupply)

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
