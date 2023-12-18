package cmd

import (
	"context"
	"log"
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/common/util"
	"github.com/sunjiangjun/xlog"
)

func Init() *Cmd {
	cfg := config.LoadConfig("./../../../cmd/collect/collect_config_tron.json")

	if cfg.LogConfig == nil {
		cfg.LogConfig.Delay = 2
		cfg.LogConfig.Path = "./log/collect"
	}

	log.Printf("%+v\n", cfg)

	nodeId, err := util.GetLocalNodeId(cfg.KeyPath)
	if err != nil {
		panic(err.Error())
	}

	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	return NewService(cfg.Chains[0], cfg.LogConfig, nodeId)
}

func TestCmd_execMultiTx(t *testing.T) {
	c := Init()
	task := &collect.NodeTask{Id: 1702870213358851508, NodeId: "f5492b85-c531-456b-bf21-61d3155fbe60", BlockNumber: "39825842", TaskType: 4, BlockChain: 198, TaskStatus: 0}
	e := xlog.NewXLogger().WithField("id", "1")
	msg := c.execMultiTx(task, e)
	if msg == nil {
		t.Error("empty")
	} else {
		t.Log("ok")
	}
}
