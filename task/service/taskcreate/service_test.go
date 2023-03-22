package taskcreate

import (
	"github.com/uduncloud/easynode/task/config"
	"log"
	"testing"
)

func Init() *Service {
	cfg := config.LoadConfig("./cmd/task/config.json")
	return NewService(&cfg)
}

func TestService_GetRecentNumber(t *testing.T) {
	s := Init()
	log.Println(s.GetRecentNumber(200))
}

func TestService_UpdateRecentNumber(t *testing.T) {
	s := Init()
	log.Println(s.UpdateRecentNumber(200, 15986827))
}

func TestService_UpdateLastNumber(t *testing.T) {
	s := Init()
	log.Println(s.UpdateLastNumber(200, 16072907))
}

func TestService_NewBlockTask(t *testing.T) {
	s := Init()
	log.Println(s.NewBlockTask(*s.config.BlockConfigs[1], nil))
}

func TestService_CreateBlockTask(t *testing.T) {
	s := Init()
	s.CreateBlockTask()
}

func TestEther_GetLastBlockNumber(t *testing.T) {
	s := Init()
	log.Println(s.GetLastBlockNumberForEther(s.config.BlockConfigs[0]))
}
