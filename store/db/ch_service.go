package db

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/0xcregis/easynode/common/driver"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/config"
	"github.com/sunjiangjun/xlog"
	"gorm.io/gorm"
)

type ClickhouseDb struct {
	chDb map[int64]*gorm.DB
	cfg  map[int64]*config.Chain

	baseDb     *gorm.DB
	baseConfig *config.Config
}

func (m *ClickhouseDb) NewSubFilter(filters []*store.SubFilter) (err error) {
	for _, v := range filters {
		list, _ := m.GetSubFilter(v.Token, v.BlockChain, v.TxCode)
		if len(list) > 0 {
			continue
		}
		v.Id = time.Now().UnixNano()
		err = m.baseDb.Table(m.baseConfig.BaseDb.FilterTable).Create(v).Error
	}

	return err
}

func (m *ClickhouseDb) DelSubFilter(id int64) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.FilterTable).Delete(store.SubFilter{Id: id}).Error
}

func (m *ClickhouseDb) DelSubFilter2(filter *store.SubFilter) error {

	whereSql := make([]string, 0, 2)
	valueSql := make([]any, 0, 2)
	if len(filter.Token) > 0 {
		whereSql = append(whereSql, "token=?")
		valueSql = append(valueSql, filter.Token)
	}

	if filter.BlockChain > 0 {
		whereSql = append(whereSql, "block_chain=?")
		valueSql = append(valueSql, filter.BlockChain)
	}

	if len(filter.TxCode) > 0 {
		whereSql = append(whereSql, "tx_code=?")
		valueSql = append(valueSql, filter.TxCode)
	}

	where := strings.Join(whereSql, " and ")

	return m.baseDb.Table(m.baseConfig.BaseDb.FilterTable).Where(where, valueSql...).Delete(store.SubFilter{}).Error
}

func (m *ClickhouseDb) GetSubFilter(token string, blockChain int64, txCode string) ([]*store.SubFilter, error) {

	whereSql := make([]string, 0, 2)
	valueSql := make([]any, 0, 2)
	if len(token) > 0 {
		whereSql = append(whereSql, "token=?")
		valueSql = append(valueSql, token)
	}

	if blockChain > 0 {
		whereSql = append(whereSql, "block_chain=?")
		valueSql = append(valueSql, blockChain)
	}

	if len(txCode) > 0 {
		whereSql = append(whereSql, "tx_code=?")
		valueSql = append(valueSql, txCode)
	}

	where := strings.Join(whereSql, " and ")

	var list []*store.SubFilter
	err := m.baseDb.Table(m.baseConfig.BaseDb.FilterTable).Where(where, valueSql...).Scan(&list).Error

	if err != nil {
		return nil, err
	}
	return list, nil
}

func (m *ClickhouseDb) GetAddressByToken2(token string) ([]*store.MonitorAddress, error) {
	var list []*store.MonitorAddress
	err := m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Select("address,block_chain").Where("token=?", token).Group("address,block_chain").Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}

func (m *ClickhouseDb) DelMonitorAddress(blockchain int64, token string, address string) error {
	err := m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Where("token=? and address=? and block_chain=?", token, address, blockchain).Delete(store.MonitorAddress{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (m *ClickhouseDb) NewToken(token *store.NodeToken) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Create(token).Error
}

func (m *ClickhouseDb) UpdateToken(token string, nodeToken *store.NodeToken) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Where("token=?", token).UpdateColumn("email", nodeToken.Email).Error
}

func (m *ClickhouseDb) GetNodeTokenByEmail(email string) (error, *store.NodeToken) {
	var r store.NodeToken
	err := m.baseDb.Table(m.baseConfig.BaseDb.TokenTable).Where("email=?", email).First(&r).Error
	if err != nil {
		return err, nil
	}
	if r.Id < 1 {
		return errors.New("no record"), nil
	}
	return nil, &r
}
func (m *ClickhouseDb) NewSubTx(blockchain int64, txs []*store.SubTx) error {
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

func (m *ClickhouseDb) NewTx(blockchain int64, txs []*store.Tx) error {
	if _, ok := m.chDb[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	if _, ok := m.cfg[blockchain]; !ok {
		return errors.New("the server has not support this blockchain")
	}
	return m.chDb[blockchain].Table(m.cfg[blockchain].ClickhouseDb.TxTable).Create(txs).Error
}

func (m *ClickhouseDb) NewBlock(blockchain int64, blocks []*store.Block) error {
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

func (m *ClickhouseDb) NewReceipt(blockchain int64, receipts []*store.Receipt) error {
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

func NewChService(cfg *config.Config, log *xlog.XLog) store.DbStoreInterface {
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

func (m *ClickhouseDb) AddMonitorAddress(blockchain int64, address *store.MonitorAddress) error {
	return m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Create(address).Error
}

func (m *ClickhouseDb) GetAddressByToken(blockchain int64, token string) ([]*store.MonitorAddress, error) {
	var list []*store.MonitorAddress
	err := m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Select("address").Where("token=? and block_chain in (?)", token, []int64{blockchain, 0}).Group("address").Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}

func (m *ClickhouseDb) GetAddressByToken3(blockchain int64) ([]*store.MonitorAddress, error) {
	var list []*store.MonitorAddress
	err := m.baseDb.Table(m.baseConfig.BaseDb.AddressTable).Select("address").Where("block_chain in (?)", []int64{blockchain, 0}).Group("address").Scan(&list).Error
	if err != nil || len(list) < 1 {
		return nil, errors.New("no record")
	}
	return list, nil
}
