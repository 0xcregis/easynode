package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/taskapi/config"
	"github.com/uduncloud/easynode/taskapi/service"
	"log"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(1).BuildFormatter(xlog.FORMAT_JSON).BuildFile("./log/task_api", 24*time.Hour)

	e := gin.Default()

	root := e.Group(cfg.RootPath)

	root.Use(gin.LoggerWithConfig(gin.LoggerConfig{Output: xLog.Out}))

	srv := service.NewServer(cfg.TaskDb, cfg.ClickhouseDb, cfg.BlockChain, xLog)

	root.GET("/node", srv.GetActiveNodes)
	root.POST("/block", srv.PushBlockTask)

	root.POST("/tx", srv.PushTxTask)
	root.POST("/syncTx", srv.PushSyncTxTask)
	root.POST("/txs", srv.PushTxsTask)

	root.POST("/receipt", srv.PushReceiptTask)
	root.POST("/receipts", srv.PushReceiptsTask)

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
