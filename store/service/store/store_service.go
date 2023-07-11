package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"github.com/uduncloud/easynode/store/service/db"
	"sync"
	"time"
)

type Service struct {
	core   service.DbMonitorAddressInterface
	log    *logrus.Entry
	config *config.Config
	kafka  *kafkaClient.EasyKafka
}

func NewStoreService(config *config.Config, log *xlog.XLog) *Service {
	ch := db.NewChService(config, log)
	kfk := kafkaClient.NewEasyKafka(log)

	return &Service{
		config: config,
		core:   ch,
		log:    log.WithField("model", "store"),
		kafka:  kfk,
	}
}

func (s *Service) Start(ctx context.Context) {
	for _, v := range s.config.Chains {
		if v.TxStore {
			go s.ReadTxFromKafka(v.BlockChain, v.KafkaCfg, ctx)
		}

		if v.BlockStore {
			go s.ReadBlockFromKafka(v.BlockChain, v.KafkaCfg, ctx)
		}

		if v.ReceiptStore {
			go s.ReadReceiptFromKafka(v.BlockChain, v.KafkaCfg, ctx)
		}

		if v.SubStore {
			go s.ReadSubTxFromKafka(v.BlockChain, v.KafkaCfg, ctx)
		}
	}

}

func (s *Service) ReadSubTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["SubTx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_subtx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*service.SubTx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewSubTx(blockChain, list)
				if err != nil {
					s.log.Errorf("ReadSubTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:

			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}
			var tx service.SubTx
			err := json.Unmarshal(msg.Value, &tx)
			if err != nil {
				s.log.Errorf("ReadSubTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				continue
			}

			lock.Lock()
			list = append(list, &tx)
			lock.Unlock()
		}
	}
}

func (s *Service) ReadTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Tx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_tx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*service.Tx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-ctx.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewTx(blockChain, list)
				if err != nil {
					s.log.Errorf("ReadTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:

			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}

			tx, err := GetTxFromKafka(msg.Value, blockChain)
			if err != nil {
				s.log.Errorf("ReadTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				continue
			}

			lock.Lock()
			list = append(list, tx)
			lock.Unlock()
		}
	}
}

func (s *Service) ReadBlockFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Block"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_block_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*service.Block, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}
	for true {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewBlock(blockChain, list)
				if err != nil {
					s.log.Errorf("ReadBlockFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()
		case msg := <-receiver:
			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}

			block, err := GetBlockFromKafka(msg.Value, blockChain)
			if err != nil {
				s.log.Errorf("ReadBlockFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				continue
			}

			lock.Lock()
			list = append(list, block)
			lock.Unlock()
		}
	}
}

func (s *Service) ReadReceiptFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Receipt"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_receipt_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*service.Receipt, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewReceipt(blockChain, list)
				if err != nil {
					s.log.Errorf("ReadReceiptFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:
			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}

			receipt, err := GetReceiptFromKafka(msg.Value, blockChain)
			if err != nil {
				s.log.Errorf("ReadReceiptFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				continue
			}

			lock.Lock()
			list = append(list, receipt)
			lock.Unlock()
		}
	}
}
