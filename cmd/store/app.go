package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service/network"
	"github.com/uduncloud/easynode/store/service/store"
	"log"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "store", "./cmd/store/config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

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
	root.POST("/monitor/address/get", srv.GetMonitorAddress)
	root.POST("/monitor/address/delete", srv.DelMonitorAddress)

	//ws 协议
	wsServer := network.NewWsHandler(&cfg, xLog)
	wsServer.Start()
	root.Handle("GET", "/ws/:token", func(ctx *gin.Context) {
		wsServer.Sub2(ctx, ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
