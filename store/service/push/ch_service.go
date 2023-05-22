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

	baseDb     *gorm.DB
	baseConfig *config.Config
}

func (m *ClickhouseDb) NewToken(token *service.NodeToken) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Create(token).Error
}

func (m *ClickhouseDb) UpdateToken(token string, nodeToken *service.NodeToken) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Where("token=?", token).UpdateColumn("email", nodeToken.Email).Error
}

func (m *ClickhouseDb) GetNodeTokenByEmail(email string) (error, *service.NodeToken) {
	var r service.NodeToken
	err := m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Where("email=?", email).First(&r).Error
	if err != nil {
		return err, nil
	}
	if r.Id < 1 {
		return errors.New("no record"), nil
	}
	return nil, &r
}
func (m *ClickhouseDb) NewSubTx(blockchain int64, txs []*service.SubTx) error {
	if _, ok := m.chDb[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	if _, ok := m.cfg[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}

	for _, v := range txs {
		bs, err := json.Marshal(v.ContractTx)
		if err == nil {
			v.ContractTxs = string(bs)
		}

		bs, err = json.Marshal(v.FeeDetail)
		if err == nil {
			v.FeeDetails = string(bs)
		}
	}

	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.SubTxTable).Create(txs).Error
}

func (m *ClickhouseDb) NewTx(blockchain int64, txs []*service.Tx) error {
	if _, ok := m.chDb[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	if _, ok := m.cfg[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.TxTable).Create(txs).Error
}

func (m *ClickhouseDb) NewBlock(blockchain int64, blocks []*service.Block) error {
	if _, ok := m.chDb[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	if _, ok := m.cfg[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	for _, v := range blocks {
		bs, _ := json.Marshal(v.Transactions)
		v.TransactionList = string(bs)
	}

	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.BlockTable).Create(blocks).Error
}

func (m *ClickhouseDb) NewReceipt(blockchain int64, receipts []*service.Receipt) error {
	if _, ok := m.chDb[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	if _, ok := m.cfg[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
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

	c, err := driver.OpenCK(cfg.BaseDb.User, cfg.BaseDb.Password, cfg.BaseDb.Addr, cfg.BaseDb.DbName, cfg.BaseDb.Port, log)
	if err != nil {
		panic(err)
	}

	m := &ClickhouseDb{
		chDb:       dbs,
		cfg:        cs,
		baseDb:     c,
		baseConfig: cfg,
	}

	return m
}

func (m *ClickhouseDb) AddMonitorAddress(blockchain int64, address *service.MonitorAddress) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Create(address).Error
}

func (m *ClickhouseDb) GetAddressByToken(blockchain int64, token string) ([]*service.MonitorAddress, error) {
	var list []*service.MonitorAddress
	err := m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Where("token=? and block_chain=?", token, blockchain).Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}
