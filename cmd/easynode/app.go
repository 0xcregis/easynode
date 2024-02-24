package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	blockchainConfig "github.com/0xcregis/easynode/blockchain/config"
	blockchainService "github.com/0xcregis/easynode/blockchain/service"
	collectConfig "github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/cmd"
	collectMonitor "github.com/0xcregis/easynode/collect/service/monitor"
	"github.com/0xcregis/easynode/common/util"
	storeConfig "github.com/0xcregis/easynode/store/config"
	"github.com/0xcregis/easynode/store/service"
	taskConfig "github.com/0xcregis/easynode/task/config"
	"github.com/0xcregis/easynode/task/service/taskcreate"
	taskapiConfig "github.com/0xcregis/easynode/taskapi/config"
	taskapiService "github.com/0xcregis/easynode/taskapi/service"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
)

// 将 collect,blockchain,task,task_api 集成为一个 app 并运行
func main() {
	var collectPath string
	var taskPath string
	var blockchainPath string
	var taskapiPath string
	var storePath string
	flag.StringVar(&collectPath, "collect", "./collect_config.json", "The system file of config")
	flag.StringVar(&taskPath, "task", "./task_config.json", "The system file of config")
	flag.StringVar(&blockchainPath, "blockchain", "./blockchain_config.json", "The system file of config")
	flag.StringVar(&taskapiPath, "taskapi", "./taskapi_config.json", "The system file of config")
	flag.StringVar(&storePath, "store", "./store_config.json", "The system file of config")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//start collect service
	go startCollect(collectPath, ctx)

	//start task service
	go startTask(taskPath, ctx)

	//start blockchain service
	go startBlockchain(blockchainPath, ctx)

	//start taskapi service
	go startTaskApi(taskapiPath, ctx)

	//start store service
	go startStore(storePath, ctx)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	cancel()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	t, c := context.WithTimeout(ctx, 2*time.Second)
	defer c()
	<-t.Done()
}

func startStore(configPath string, ctx context.Context) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := storeConfig.LoadConfig(configPath)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFile("./log/store/store", 24*time.Hour)

	//是否存盘
	service.NewStoreHandler(&cfg, xLog).Start(ctx)

	//http 协议
	e := gin.Default()
	root := e.Group(cfg.RootPath)
	srv := service.NewHttpHandler(&cfg, xLog)
	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))
	root.POST("/monitor/token", srv.NewToken)
	root.POST("/monitor/address", srv.MonitorAddress)
	root.POST("/monitor/address/get", srv.GetMonitorAddress)
	root.POST("/monitor/address/delete", srv.DelMonitorAddress)
	root.POST("/filter/new", srv.AddSubFilter)
	root.POST("/filter/get", srv.QuerySubFilter)
	root.POST("/filter/delete", srv.DelSubFilter)

	//ws 协议
	wsServer := service.NewWsHandler(&cfg, xLog)
	wsServer.Start(ctx)
	root.Handle("GET", "/ws/:token", func(ctx *gin.Context) {
		wsServer.Sub2(ctx, ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}

func startTaskApi(configPath string, ctx context.Context) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := taskapiConfig.LoadConfig(configPath)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFile("./log/taskapi/task_api", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := taskapiService.NewServer(&cfg, cfg.BlockChain, xLog)

	//root.GET("/node", srv.GetActiveNodes)
	root.POST("/block", srv.PushBlockTask)

	root.POST("/tx", srv.PushTxTask)
	//root.POST("/syncTx", srv.PushSyncTxTask)
	root.POST("/txs", srv.PushTxsTask)

	root.POST("/receipt", srv.PushReceiptTask)
	root.POST("/receipts", srv.PushReceiptsTask)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}

func startBlockchain(configPath string, ctx context.Context) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := blockchainConfig.LoadConfig(configPath)
	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain/chain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := blockchainService.NewHttpHandler(cfg.Cluster, cfg.Kafka, xLog)

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
	myRoot.POST("/tx/sendRawTransaction", srv.SendRawTx1)
	myRoot.POST("/account/getAccountResource", srv.GetAccountResource)
	myRoot.POST("/gas/estimateGasForTron", srv.EstimateGasForTron)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}

func startTask(configPath string, ctx context.Context) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := taskConfig.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	//生产任务 服务
	if cfg.AutoCreateBlockTask {
		taskcreate.NewService(&cfg).Start(ctx)
	}
}

func startCollect(configPath string, ctx context.Context) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := collectConfig.LoadConfig(configPath)

	if cfg.LogConfig == nil {
		cfg.LogConfig.Delay = 2
		cfg.LogConfig.Path = "./log/collect"
	}

	log.Printf("%+v\n", cfg)

	nodeId, err := util.GetLocalNodeId(cfg.KeyPath)
	if err != nil {
		panic(err.Error())
	}

	//启动监控服务
	collectMonitor.NewService(&cfg, nodeId).Start(ctx)

	//启动公链服务
	for _, v := range cfg.Chains {
		go cmd.NewService(v, cfg.LogConfig, nodeId).Start(ctx)
	}
}
