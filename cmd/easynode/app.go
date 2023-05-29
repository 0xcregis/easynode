package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	blockchainConfig "github.com/uduncloud/easynode/blockchain/config"
	blockchainService "github.com/uduncloud/easynode/blockchain/service"
	collectConfig "github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service/cmd"
	collectMonitor "github.com/uduncloud/easynode/collect/service/monitor"
	storeConfig "github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service/network"
	"github.com/uduncloud/easynode/store/service/store"
	taskConfig "github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service/taskcreate"
	taskapiConfig "github.com/uduncloud/easynode/taskapi/config"
	taskapiService "github.com/uduncloud/easynode/taskapi/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	//start collect service
	go startCollect(collectPath)

	//start task service
	go startTask(taskPath)

	//start blockchain service
	go startBlockchain(blockchainPath)

	//start taskapi service
	go startTaskApi(taskapiPath)

	//start store service
	go startStore(storePath)

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
}

func startStore(configPath string) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := storeConfig.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/store/store", 24*time.Hour)

	//是否存盘
	store.NewStoreService(&cfg, xLog).Start()

	//http 协议
	e := gin.Default()
	root := e.Group(cfg.RootPath)
	srv := network.NewServer(&cfg, xLog)
	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))
	root.POST("/monitor/token", srv.NewToken)
	root.POST("/monitor/address", srv.MonitorAddress)

	//ws 协议
	wsServer := network.NewWsHandler(&cfg, xLog)
	root.Handle("GET", "/ws/:token", func(ctx *gin.Context) {
		wsServer.Start(ctx, ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}

func startTaskApi(configPath string) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := taskapiConfig.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/taskapi/task_api", 24*time.Hour)

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

func startBlockchain(configPath string) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := blockchainConfig.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/blockchain/chain", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := blockchainService.NewHttpHandler(cfg.Cluster, xLog)
	//支持JSON-RPC协议的公链
	root.POST("/jsonrpc", srv.HandlerReq)

	//自定义或不支持JSON-RPC协议的公链
	root.POST("/block/hash", srv.GetBlockByHash)
	root.POST("/block/number", srv.GetBlockByNumber)
	root.POST("/tx/hash", srv.GetTxByHash)
	root.POST("/tx/receipts", srv.GetTxReceiptByHash)
	root.POST("/account/balance", srv.GetBalance)
	root.POST("/account/tokenBalance", srv.GetTokenBalance)
	root.POST("/account/nonce", srv.GetNonce)
	root.POST("/block/latest", srv.GetLatestBlock)
	root.POST("/tx/sendRawTransaction", srv.SendRawTx)

	//ws 协议
	wsServer := blockchainService.NewWsHandler(cfg.Cluster, xLog)
	root.Handle("GET", "/ws", func(ctx *gin.Context) {
		wsServer.Start(ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}

func startTask(configPath string) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := taskConfig.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	//生产任务 服务
	if cfg.AutoCreateBlockTask {
		taskcreate.NewService(&cfg).Start()
	}
}

func startCollect(configPath string) {

	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := collectConfig.LoadConfig(configPath)

	if cfg.LogConfig == nil {
		cfg.LogConfig.Delay = 2
		cfg.LogConfig.Path = "./log/collect"
	}

	log.Printf("%+v\n", cfg)

	//启动监控服务
	collectMonitor.NewService(&cfg).Start()

	//启动公链服务
	for _, v := range cfg.Chains {
		cmd.NewService(v, cfg.LogConfig).Start()
	}
}
