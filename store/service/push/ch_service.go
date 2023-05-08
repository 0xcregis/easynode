package push

import (
	"encoding/json"
	"errors"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/common/driver"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"gorm.io/gorm"
)

type ClickhouseDb struct {
	chDb map[int64]*gorm.DB
	cfg  map[int64]*config.Chain
}

func (m *ClickhouseDb) NewTx(blockchain int64, txs []*service.Tx) error {
	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.TxTable).Create(txs).Error
}

func (m *ClickhouseDb) NewBlock(blockchain int64, blocks []*service.Block) error {

	for _, v := range blocks {
		bs, _ := json.Marshal(v.Transactions)
		v.TransactionList = string(bs)
	}

	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.BlockTable).Create(blocks).Error
}

func (m *ClickhouseDb) NewReceipt(blockchain int64, receipts []*service.Receipt) error {
	for _, v := range receipts {
		bs, _ := json.Marshal(v.Logs)
		v.LogList = string(bs)
	}
	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.BlockTable).Create(receipts).Error
}

func NewChService(cfg *config.Config, log *xlog.XLog) service.DbMonitorAddressInterface {
	dbs := make(map[int64]*gorm.DB, 2)
	cs := make(map[int64]*config.Chain, 2)
	for _, v := range cfg.Chains {
		c, err := driver.OpenCK(v.ClickhouseDb.User, v.ClickhouseDb.Password, v.ClickhouseDb.Addr, v.ClickhouseDb.DbName, v.ClickhouseDb.Port, log)
		if err != nil {
			panic(err)
		}
		dbs[v.BlockChain] = c
		cs[v.BlockChain] = v
	}

	m := &ClickhouseDb{
		chDb: dbs,
		cfg:  cs,
	}

	return m
}

func (m *ClickhouseDb) AddMonitorAddress(blockchain int64, address *service.MonitorAddress) error {
	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.AddressTable).Create(address).Error
}

func (m *ClickhouseDb) GetAddressByToken(blockchain int64, token string) ([]*service.MonitorAddress, error) {
	var list []*service.MonitorAddress
	err := m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.AddressTable).Where("token=? and block_chain=?", token, blockchain).Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}
