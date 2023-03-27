package main

import (
	"context"
	"flag"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service/cmd/task"
	"github.com/uduncloud/easynode/collect/service/monitor"
	"github.com/uduncloud/easynode/collect/service/nodeinfo"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "config", "./cmd/collect/config.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	if cfg.LogConfig == nil {
		cfg.LogConfig.Delay = 2
		cfg.LogConfig.Path = "./log/collect"
	}

	log.Printf("%+v\n", cfg)

	//上传节点信息 服务
	nodeinfo.NewService(cfg.NodeInfoDb, cfg.Chains, cfg.LogConfig).Start()

	//启动处理日志服务
	monitor.NewService(cfg.LogConfig).Start()

	//启动公链服务
	for _, v := range cfg.Chains {
		//GetBlockChainService(v, cfg.TaskDb, cfg.SourceDb).Start()
		task.NewService(v, cfg.TaskDb, cfg.SourceDb, cfg.LogConfig).Start()
	}

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
