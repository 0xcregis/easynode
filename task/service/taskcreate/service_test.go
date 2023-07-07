package taskcreate

import (
	"github.com/uduncloud/easynode/task/config"
	"log"
	"testing"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../../cmd/task/config.json")
	return NewService(&cfg)
}

func TestService_GetRecentNumber(t *testing.T) {
	s := Init()
	log.Println(s.store.GetRecentNumber(205))
}

func TestService_UpdateRecentNumber(t *testing.T) {
	s := Init()
	log.Println(s.store.UpdateRecentNumber(205, 15986827))
}

func TestService_UpdateLastNumber(t *testing.T) {
	s := Init()
	log.Println(s.store.UpdateLastNumber(205, 16072907))
}

func TestService_NewBlockTask(t *testing.T) {
	s := Init()
	log.Println(s.NewBlockTask(*s.config.BlockConfigs[1], nil))
}

func TestEther_GetLastBlockNumber(t *testing.T) {

}
