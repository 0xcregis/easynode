package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/0xcregis/easynode/store/config"
	"github.com/0xcregis/easynode/store/service"
	"github.com/gin-gonic/gin"
	"github.com/sunjiangjun/xlog"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "store", "./cmd/store/config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	if cfg.LogLevel == 0 {
		cfg.LogLevel = 4
	}
	log.Printf("%+v\n", cfg)

	xLog := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildLevel(xlog.Level(cfg.LogLevel)).BuildFile("./log/store/store", 24*time.Hour)

	kafkaCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//是否存盘
	service.NewStoreHandler(&cfg, xLog).Start(kafkaCtx)

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
	wsServer.Start(kafkaCtx)
	root.Handle("GET", "/ws/:token", func(ctx *gin.Context) {
		wsServer.Sub2(ctx, ctx.Writer, ctx.Request)
	})

	err := e.Run(fmt.Sprintf(":%v", cfg.Port))
	if err != nil {
		panic(err)
	}
}
