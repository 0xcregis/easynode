package monitor

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"testing"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../../cmd/collect/config.json")
	return NewService(&cfg, cfg.LogConfig, xlog.NewXLogger())
}

func TestService_Start(t *testing.T) {
	s := Init()
	s.Start()
}

func TestService_CheckTable(t *testing.T) {
	s := Init()
	s.CheckTable()
}
