package push

import (
	"context"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"log"
	"testing"
	"time"
)

func Init() service.DbMonitorAddressInterface {
	cfg := config.LoadConfig("./../../../cmd/store/config.json")
	return NewChService(&cfg, xlog.NewXLogger())
}

func TestClickhouseDb_AddMonitorAddress(t *testing.T) {
	s := Init()
	log.Println(s.AddMonitorAddress(200, &service.MonitorAddress{Address: "0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8", Token: "5fe5f231-7051-4caf-9b52-108db92edbb4", BlockChain: 200, TxType: "0x2", Id: time.Now().UnixMilli()}))
}

func TestClickhouseDb_GetAddressByToken(t *testing.T) {
	s := Init()
	log.Println(s.GetAddressByToken(200, "5fe5f231-7051-4caf-9b52-108db92edbb4"))
}

func TestName(t *testing.T) {

	ctx, cancel := context.WithCancel(context.Background())

	//c, _ := context.WithCancel(ctx)

	go func(ctx2 context.Context) {
		for {
			select {
			case <-ctx2.Done():
				log.Println("2 done")
				return
			}
		}
	}(ctx)

	go func(ctx2 context.Context) {
		for {
			select {
			case <-ctx2.Done():
				log.Println("1 done")
				return
			}
		}
	}(ctx)

	time.Sleep(5 * time.Second)
	cancel()

	time.Sleep(5 * time.Second)
}
