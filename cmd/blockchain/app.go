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
	flag.StringVar(&configPath, "blockchain", "./cmd/blockchain/config_ether.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain/chain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := service.NewHttpHandler(cfg.Cluster, cfg.Kafka, xLog)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//start kafka listener
	srv.StartKafka(ctx)

	origin := root.Group("origin")
	//支持JSON-RPC协议的公链
	origin.POST("/jsonrpc", srv.HandlerReq)

	//自定义或不支持JSON-RPC协议的公链
	origin.POST("/block/hash", srv.GetBlockByHash)
	origin.POST("/block/number", srv.GetBlockByNumber)
	origin.POST("/tx/hash", srv.GetTxByHash)
	origin.POST("/tx/trace", srv.GetTraceTransaction)
	origin.POST("/receipts/hash", srv.GetTxReceiptByHash)
	origin.POST("/account/balance", srv.GetBalance)
	origin.POST("/account/token", srv.GetToken)
	origin.POST("/account/tokenBalance", srv.GetTokenBalance)
	origin.POST("/account/nonce", srv.GetNonce)
	origin.POST("/block/latest", srv.GetLatestBlock)
	origin.POST("/tx/sendRawTransaction", srv.SendRawTx)
	origin.POST("/gas/price", srv.GasPrice)
	origin.POST("/gas/estimateGas", srv.EstimateGas)
	origin.POST("/nft/tokenUri", srv.TokenUri)
	origin.POST("/nft/balanceOf", srv.BalanceOf)
	origin.POST("/nft/owner", srv.OwnerOf)
	origin.POST("/nft/totalSupply", srv.TotalSupply)

	myRoot := root.Group("easynode")
	myRoot.POST("/block/hash", srv.GetBlockByHash1)
	myRoot.POST("/block/number", srv.GetBlockByNumber1)
	myRoot.POST("/tx/hash", srv.GetTxByHash1)
	myRoot.POST("/account/balance", srv.GetBalance1)
	myRoot.POST("/account/token", srv.GetToken1)
	myRoot.POST("/account/tokenBalance", srv.GetTokenBalance1)
	myRoot.POST("/account/nonce", srv.GetNonce1)
	myRoot.POST("/block/latest", srv.GetLatestBlock1)
	myRoot.POST("/gas/price", srv.GasPrice1)
	myRoot.POST("/gas/estimateGas", srv.EstimateGas1)
	myRoot.POST("/gas/estimateGasForTron", srv.EstimateGasForTron)
	myRoot.POST("/tx/sendRawTransaction", srv.SendRawTx1)
	myRoot.POST("/account/getAccountResource", srv.GetAccountResource)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
