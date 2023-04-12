package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/taskapi/common"
	"github.com/uduncloud/easynode/taskapi/config"
	"gorm.io/gorm"
	"time"
)

type ClickhouseDb struct {
	chDb        map[int64]*gorm.DB
	cfg         *config.Config
	kafkaClient *kafkaClient.EasyKafka
	sendCh      chan []*kafka.Message
	receiverCh  chan []*kafka.Message
}

func NewChService(cfg *config.Config, log *xlog.XLog) DbApiInterface {
	kf := kafkaClient.NewEasyKafka(log)
	sendCh := make(chan []*kafka.Message, 10)
	receiverCh := make(chan []*kafka.Message, 5)

	mp := make(map[int64]*gorm.DB, 2)
	if len(cfg.ClickhouseDb) > 0 {
		for k, v := range cfg.ClickhouseDb {
			c, err := common.OpenCK(v.User, v.Password, v.Addr, v.DbName, v.Port, log)
			if err != nil {
				panic(err)
			}
			mp[k] = c
		}
	} else {
		log.Error("some function does not work for clickhouse`s config is null")
		panic("some function does not work for clickhouse`s config is null")
	}

	m := &ClickhouseDb{
		chDb:        mp,
		cfg:         cfg,
		kafkaClient: kf,
		sendCh:      sendCh,
		receiverCh:  receiverCh,
	}

	go func() {
		m.startKafka()
	}()
	return m
}

func (m *ClickhouseDb) startKafka() {
	broker := fmt.Sprintf("%v:%v", m.cfg.Kafka.Host, m.cfg.Kafka.Port)
	m.kafkaClient.Write(kafkaClient.Config{Brokers: []string{broker}}, m.sendCh, nil)
}

func (m *ClickhouseDb) AddNodeTask(task *NodeTask) error {
	task.CreateTime = time.Now()
	task.LogTime = time.Now()
	task.Id = time.Now().UnixNano()
	bs, err := json.Marshal(task)
	if err != nil {
		return err
	}
	msg := &kafka.Message{Topic: fmt.Sprintf("task_%v", task.BlockChain), Partition: 0, Key: []byte(task.NodeId), Value: bs}
	m.sendCh <- []*kafka.Message{msg}
	return nil
}

func (m *ClickhouseDb) QueryTxFromCh(blockChain int64, txHash string) (*Tx, error) {
	//clickhouse 非必须配置项，因此 可能不存此次连接
	if _, ok := m.chDb[blockChain]; !ok {
		return nil, errors.New("not found db source ,please check config file")
	}

	var tx Tx
	err := m.chDb[blockChain].Table(m.cfg.ClickhouseDb[blockChain].TxTable).Where("hash=?", txHash).Scan(&tx).Error
	if err != nil || tx.Id < 1 {
		return nil, errors.New("no record")
	}
	return &tx, nil
}

func (m *ClickhouseDb) AddMonitorAddress(blockchain int64, address *MonitorAddress) error {
	return m.chDb[blockchain].Table(m.cfg.ClickhouseDb[blockchain].AddressTable).Create(address).Error
}
