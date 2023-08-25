package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	kafkaClient "github.com/0xcregis/easynode/common/kafka"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/chain"
	"github.com/0xcregis/easynode/store/config"
	"github.com/0xcregis/easynode/store/db"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

type StoreHandler struct {
	store  store.DbStoreInterface
	log    *logrus.Entry
	config *config.Config
	kafka  *kafkaClient.EasyKafka
}

func NewStoreHandler(config *config.Config, log *xlog.XLog) *StoreHandler {
	ch := db.NewChService(config, log)
	kfk := kafkaClient.NewEasyKafka(log)

	return &StoreHandler{
		config: config,
		store:  ch,
		log:    log.WithField("model", "store"),
		kafka:  kfk,
	}
}

func (s *StoreHandler) Start(ctx context.Context) {
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

		if v.BackupTxStore {
			go s.ReadBackupTxFromKafka(v.BlockChain, v.KafkaCfg, ctx)
		}
	}

}

func (s *StoreHandler) ReadBackupTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["BackupTx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_backuptx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: fmt.Sprintf("%v_%v", Kafka.Topic, blockChain), Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*store.BackupTx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.store.NewBackupTx(blockChain, list)
				if err != nil {
					s.log.Errorf("ReadBackupTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:

			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}
			var tx store.BackupTx
			err := json.Unmarshal(msg.Value, &tx)
			if err != nil {
				s.log.Errorf("ReadBackupTxFromKafka|blockChain=%v,error=%v", blockChain, err.Error())
				continue
			}

			lock.Lock()
			list = append(list, &tx)
			lock.Unlock()
		}
	}
}

func (s *StoreHandler) ReadSubTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["SubTx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_subtx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*store.SubTx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.store.NewSubTx(blockChain, list)
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
			var tx store.SubTx
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

func (s *StoreHandler) ReadTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Tx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_tx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*store.Tx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for {
		select {
		case <-ctx.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.store.NewTx(blockChain, list)
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

			tx, err := chain.GetTxFromKafka(msg.Value, blockChain)
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

func (s *StoreHandler) ReadBlockFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Block"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_block_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*store.Block, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}
	for {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.store.NewBlock(blockChain, list)
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

			block, err := chain.GetBlockFromKafka(msg.Value, blockChain)
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

func (s *StoreHandler) ReadReceiptFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Receipt"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("gr_store_receipt_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx2)
	}()

	list := make([]*store.Receipt, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for {
		select {
		case <-ctx2.Done():
			tk.Stop()
			return
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.store.NewReceipt(blockChain, list)
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

			receipt, err := chain.GetReceiptFromKafka(msg.Value, blockChain)
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
