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
	chDb *gorm.DB
	cfg  *config.Config
}

func (m *ClickhouseDb) NewTx(txs []*service.Tx) error {
	return m.chDb.Table(m.cfg.ClickhouseDb.TxTable).Create(txs).Error
}

func (m *ClickhouseDb) NewBlock(blocks []*service.Block) error {

	for _, v := range blocks {
		bs, _ := json.Marshal(v.Transactions)
		v.TransactionList = string(bs)
	}

	return m.chDb.Table(m.cfg.ClickhouseDb.BlockTable).Create(blocks).Error
}

func (m *ClickhouseDb) NewReceipt(receipts []*service.Receipt) error {
	for _, v := range receipts {
		bs, _ := json.Marshal(v.Logs)
		v.LogList = string(bs)
	}
	return m.chDb.Table(m.cfg.ClickhouseDb.BlockTable).Create(receipts).Error
}

func NewChService(cfg *config.Config, log *xlog.XLog) service.DbMonitorAddressInterface {
	c, err := driver.OpenCK(cfg.ClickhouseDb.User, cfg.ClickhouseDb.Password, cfg.ClickhouseDb.Addr, cfg.ClickhouseDb.DbName, cfg.ClickhouseDb.Port, log)
	if err != nil {
		panic(err)
	}
	m := &ClickhouseDb{
		chDb: c,
		cfg:  cfg,
	}

	return m
}

func (m *ClickhouseDb) AddMonitorAddress(blockchain int64, address *service.MonitorAddress) error {
	return m.chDb.Table(m.cfg.ClickhouseDb.AddressTable).Create(address).Error
}

func (m *ClickhouseDb) GetAddressByToken(blockchain int64, token string) ([]*service.MonitorAddress, error) {
	var list []*service.MonitorAddress
	err := m.chDb.Table(m.cfg.ClickhouseDb.AddressTable).Where("token=? and block_chain=?", token, blockchain).Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}
