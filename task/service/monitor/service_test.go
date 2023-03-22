package monitor

import (
	"github.com/uduncloud/easynode/task/config"
	"testing"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../config.json")
	return NewService(&cfg)
}

func TestService_HandlerDeadTask(t *testing.T) {
	s := Init()
	s.HandlerDeadTask()

	//log.Println(time.Now().Add(-1 * time.Hour).UTC())
}

func TestService_HandlerManyFailTask(t *testing.T) {
	s := Init()
	s.HandlerManyFailTask()
	//s.createNodeTaskTable()
}
