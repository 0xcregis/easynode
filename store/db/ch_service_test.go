package db

import (
	"log"
	"strings"
	"testing"
	"time"

	"github.com/0xcregis/easynode/common/util"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/config"
	"github.com/sunjiangjun/xlog"
)

func Init() store.DbStoreInterface {
	cfg := config.LoadConfig("./../../cmd/store/config_tron.json")
	return NewChService(&cfg, xlog.NewXLogger())
}

func TestClickhouseDb_AddMonitorAddress(t *testing.T) {
	s := Init()
	log.Println(s.AddMonitorAddress(0, &store.MonitorAddress{Address: "0xe5cB067E90D5Cd1F8052B83562Ae670bA4A211a8", Token: "5fe5f231-7051-4caf-9b52-108db92edbb4", BlockChain: 0, TxType: "0x2", Id: time.Now().UnixMilli()}))
}

func TestClickhouseDb_GetAddressByToken(t *testing.T) {
	s := Init()
	log.Println(s.GetAddressByToken(200, "5fe5f231-7051-4caf-9b52-108db92edbb4"))
}

func TestClickhouseDb_NewToken(t *testing.T) {
	s := Init()
	log.Println(s.NewToken(&store.NodeToken{
		Token: "1231token",
		Email: "123@qq.com",
		Id:    time.Now().UnixMilli(),
	}))
}

func TestClickhouseDb_UpdateToken(t *testing.T) {

	s := Init()
	log.Println(s.UpdateToken("1231token", &store.NodeToken{
		Token: "1231token",
		Email: "125@qq.com",
		Id:    time.Now().UnixMilli(),
	}))
}

func TestClickhouseDb_GetNodeTokenByEmail(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTokenByEmail("125@qq.com"))
}

func TestClickhouseDb_GetAddressByToken2(t *testing.T) {
	s := Init()
	//log.Println(s.GetAddressByToken2("36ee0ad5-f4bc-4bca-a1dc-c51db006e249"))

	list, err := s.GetAddressByToken4("90c232ef-ecea-4956-9132-325ebec46e86")
	if err != nil {
		t.Error(err)
	} else {
		for _, v := range list {
			if !strings.HasPrefix(v.Address, "0x") {
				base58Addr, err := util.Base58ToAddress(v.Address)
				if err != nil {
					t.Error(err)
					continue
				}
				v.Address = base58Addr.Hex()
				err = s.UpdateMonitorAddress(v.Id, v)
				if err != nil {
					t.Error(err)
				} else {
					t.Log("ok")
				}
			}

		}
	}

}

func TestClickhouseDb_UpdateMonitorAddress(t *testing.T) {
	s := Init()
	addr := &store.MonitorAddress{
		Address:    "0x4193465903d04f7890bb9b0b856ab31fc2929ab653",
		Id:         1702881171946676000,
		Token:      "73e89458-17ce-48c5-ad5a-58b9bab121c9",
		TxType:     "0",
		BlockChain: 198,
	}
	err := s.UpdateMonitorAddress(1702881171946676000, addr)
	if err != nil {
		t.Error(err)
	} else {
		t.Log("ok")
	}
}

func TestClickhouseDb_DelMonitorAddress(t *testing.T) {
	s := Init()
	log.Println(s.DelMonitorAddress(200, "36ee0ad5-f4bc-4bca-a1dc-c51db006e249", "0x28c6c06298d514db089934071355e5743bf21d61"))
}

func TestClickhouseDb_NewSubFilter(t *testing.T) {
	s := Init()
	t.Log(s.NewSubFilter([]*store.SubFilter{{Token: "token", BlockChain: 205, TxCode: "1"}}))
}

func TestClickhouseDb_NewBackupTx(t *testing.T) {
	s := Init()
	t.Log(s.NewBackupTx(200, []*store.BackupTx{{ChainCode: 200, ID: uint64(time.Now().UnixMilli()), From: "0x123", To: "0x456", Signed: "0xdddd", Status: 1, Response: "{\"id\":\"12312321\"}"}}))
}

func TestClickhouseDb_GetSubFilter(t *testing.T) {
	s := Init()
	t.Log(s.GetSubFilter("token", 205, "1"))
}

func TestClickhouseDb_DelSubFilter2(t *testing.T) {
	s := Init()
	t.Log(s.DelSubFilter2(&store.SubFilter{Token: "token"}))
}

func TestClickhouseDb_DelSubFilter(t *testing.T) {
	s := Init()
	t.Log(s.DelSubFilter(1691999490145550000))
}
