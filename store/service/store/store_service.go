package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"github.com/uduncloud/easynode/store/service/push"
)

type StoreService struct {
	core   service.DbMonitorAddressInterface
	log    *xlog.XLog
	config *config.Config
	kafka  *kafkaClient.EasyKafka
}

func NewStoreService(config *config.Config, log *xlog.XLog) *StoreService {
	ch := push.NewChService(config, log)
	kfk := kafkaClient.NewEasyKafka(log)

	return &StoreService{
		config: config,
		core:   ch,
		log:    log,
		kafka:  kfk,
	}
}

func (s *StoreService) Start() {
	if s.config.TxStore {
		go s.readTxFromKafka(s.config.BlockChain)
	}

	if s.config.BlockStore {
		go s.readBlockFromKafka(s.config.BlockChain)
	}

	if s.config.ReceiptStore {
		go s.readReceiptFromKafka(s.config.BlockChain)
	}
}

func (s *StoreService) readTxFromKafka(blockChain int64) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := s.config.KafkaCfg["Tx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_tx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	for true {
		select {
		case msg := <-receiver:
			var tx service.Tx
			err := json.Unmarshal(msg.Value, &tx)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				continue
			}
			err = s.core.NewTx(&tx)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
			}
		}
	}
}

func (s *StoreService) readBlockFromKafka(blockChain int64) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := s.config.KafkaCfg["Block"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_block_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	for true {
		select {
		case msg := <-receiver:
			var block service.Block
			err := json.Unmarshal(msg.Value, &block)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				continue
			}
			err = s.core.NewBlock(&block)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
			}
		}
	}
}

func (s *StoreService) readReceiptFromKafka(blockChain int64) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := s.config.KafkaCfg["Receipt"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_receipt_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	for true {
		select {
		case msg := <-receiver:
			var receipt service.Receipt
			err := json.Unmarshal(msg.Value, &receipt)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				continue
			}
			err = s.core.NewReceipt(&receipt)
			if err != nil {
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
			}
		}
	}
}
