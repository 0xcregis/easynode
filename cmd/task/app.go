package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/0xcregis/easynode/task/config"
	"github.com/0xcregis/easynode/task/service/taskcreate"
)

func main() {
	var configPath string
	flag.StringVar(&configPath, "task", "./cmd/task/config_btc.json", "The system file of config")
	flag.Parse()
	if len(configPath) < 1 {
		panic("can not find config file")
	}
	cfg := config.LoadConfig(configPath)

	log.Printf("%+v\n", cfg)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//生产任务 服务
	if cfg.AutoCreateBlockTask {
		taskcreate.NewService(&cfg).Start(ctx)
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
	cancel()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	t, c := context.WithTimeout(context.Background(), 2*time.Second)
	defer c()
	<-t.Done()
}
