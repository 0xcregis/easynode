package monitor

import (
	"github.com/uduncloud/easynode/collect/config"
	"testing"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../config.json")
	return NewService(cfg.LogConfig)
}

func TestService_Start(t *testing.T) {
	s := Init()
	s.Start()
}
