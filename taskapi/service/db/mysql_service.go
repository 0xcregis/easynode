package db

import (
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/taskapi/common"
	"github.com/uduncloud/easynode/taskapi/config"
	"github.com/uduncloud/easynode/taskapi/service"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type MysqlDb struct {
	coreDd   *gorm.DB
	chDb     map[int64]*gorm.DB
	chConfig map[int64]*config.ClickhouseDb
}

func NewMysqlService(dbConfig *config.TaskDb, chConfig map[int64]*config.ClickhouseDb, log *xlog.XLog) service.DbApiInterface {
	g, err := common.Open(dbConfig.User, dbConfig.Password, dbConfig.Addr, dbConfig.DbName, dbConfig.Port, log)
	if err != nil {
		panic(err)
	}

	//clickhouse 配置非必须
	mp := make(map[int64]*gorm.DB, 2)
	if len(chConfig) > 0 {
		for k, v := range chConfig {
			c, err := common.OpenCK(v.User, v.Password, v.Addr, v.DbName, v.Port, log)
			if err != nil {
				panic(err)
			}
			mp[k] = c
		}
	} else {
		log.Warnf("some function does not work for clickhouse`s config is null")
	}

	return &MysqlDb{
		coreDd:   g,
		chDb:     mp,
		chConfig: chConfig,
	}
}

func getNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", service.NodeTaskTable, time.Now().Format(service.DayFormat))
	return table
}

func (m *MysqlDb) AddNodeTask(task *service.NodeTask) error {
	return m.coreDd.Table(getNodeTaskTable()).Omit("create_time,log_time,id").Clauses(clause.Insert{Modifier: "IGNORE"}).Create(task).Error

}

func (m *MysqlDb) GetActiveNodesFromDB() ([]string, error) {
	t := time.Now().Add(-5 * time.Minute).Format(service.TimeFormat)
	var nodeList []string
	err := m.coreDd.Table(getNodeTaskTable()).Select("node_id").Where("create_time>?", t).Pluck("node_id", &nodeList).Error

	if err != nil {
		return nil, err
	}

	return nodeList, nil
}

func (m *MysqlDb) QueryTxFromCh(blockChain int64, txHash string) (*service.Tx, error) {
	//clickhouse 非必须配置项，因此 可能不存此次连接
	if _, ok := m.chDb[blockChain]; !ok {
		return nil, errors.New("not found db source ,please check config file")
	}

	var tx service.Tx
	err := m.chDb[blockChain].Table(m.chConfig[blockChain].TxTable).Where("hash=?", txHash).Scan(&tx).Error
	if err != nil || tx.Id < 1 {
		return nil, errors.New("no record")
	}
	return &tx, nil
}
