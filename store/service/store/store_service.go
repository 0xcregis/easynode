package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"github.com/uduncloud/easynode/store/service/push"
	"sync"
	"time"
)

type StoreService struct {
	core   service.DbMonitorAddressInterface
	log    *logrus.Entry
	config *config.Config
	kafka  *kafkaClient.EasyKafka
}

func NewStoreService(config *config.Config, log *xlog.XLog) *StoreService {
	ch := push.NewChService(config, log)
	kfk := kafkaClient.NewEasyKafka(log)

	return &StoreService{
		config: config,
		core:   ch,
		log:    log.WithField("model", "store"),
		kafka:  kfk,
	}
}

func (s *StoreService) Start() {
	for _, v := range s.config.Chains {
		if v.TxStore {
			go s.readTxFromKafka(v.BlockChain, v.KafkaCfg)
		}

		if v.BlockStore {
			go s.readBlockFromKafka(v.BlockChain, v.KafkaCfg)
		}

		if v.ReceiptStore {
			go s.readReceiptFromKafka(v.BlockChain, v.KafkaCfg)
		}

		if v.SubStore {
			go s.readSubTxFromKafka(v.BlockChain, v.KafkaCfg)
		}
	}

}

func (s *StoreService) readSubTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := kafkaCfg["SubTx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_subtx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	list := make([]*service.SubTx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewSubTx(blockChain, list)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
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
				s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				continue
			}

			lock.Lock()
			list = append(list, &tx)
			lock.Unlock()
		}
	}
}

func (s *StoreService) readTxFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Tx"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_tx_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	list := make([]*service.Tx, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewTx(blockChain, list)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:

			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}
			var tx service.Tx
			if blockChain == 200 {
				err := json.Unmarshal(msg.Value, &tx)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
					continue
				}
			} else if blockChain == 205 {

				txBody := gjson.ParseBytes(msg.Value).Get("tx").String()
				if len(txBody) < 5 {
					continue
				}
				txRoot := gjson.Parse(txBody)

				status := txRoot.Get("ret.0.contractRet").String()
				hash := txRoot.Get("txID").String()
				blockHash := txRoot.Get("raw_data.ref_block_hash").String()
				txTime := txRoot.Get("raw_data.timestamp").Uint()
				limit := txRoot.Get("raw_data.fee_limit").Uint()
				txType := txRoot.Get("raw_data.contract.0.type").String()
				v := txRoot.Get("raw_data.contract.0.parameter.value")
				from := v.Get("owner_address").String()
				var to string
				if v.Get("receiver_address").Exists() {
					to = v.Get("receiver_address").String()
				}

				if v.Get("to_address").Exists() {
					to = v.Get("to_address").String()
				}

				if v.Get("contract_address").Exists() {
					to = v.Get("contract_address").String()
				}

				var input string
				if v.Get("data").Exists() {
					input = v.Get("data").String()
				}

				txValue := v.String()
				tx.Id = uint64(time.Now().UnixNano())
				tx.Value = txValue
				tx.BlockHash = blockHash
				tx.TxHash = hash
				tx.TxStatus = status
				tx.TxTime = fmt.Sprintf("%v", txTime)
				tx.ToAddr = to
				tx.FromAddr = from
				tx.Type = txType
				tx.InputData = input
				tx.GasLimit = fmt.Sprintf("%v", limit)
				receiptBody := gjson.ParseBytes(msg.Value).Get("receipt").String()
				if len(receiptBody) > 5 {
					receiptRoot := gjson.Parse(receiptBody)
					fee := receiptRoot.Get("fee").Uint()
					number := receiptRoot.Get("blockNumber").Uint()
					tx.Fee = fmt.Sprintf("%v", fee)
					tx.BlockNumber = fmt.Sprintf("%v", number)
				}
			}

			lock.Lock()
			list = append(list, &tx)
			lock.Unlock()
		}
	}
}

func (s *StoreService) readBlockFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Block"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_block_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	list := make([]*service.Block, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}
	for true {
		select {
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewBlock(blockChain, list)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()
		case msg := <-receiver:
			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}
			var block service.Block
			if blockChain == 200 {
				err := json.Unmarshal(msg.Value, &block)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
					continue
				}
			} else if blockChain == 205 {
				r := gjson.ParseBytes(msg.Value)
				hash := r.Get("blockID").String()
				number := r.Get("block_header.raw_data").Get("number").Int()
				txRoot := r.Get("block_header.raw_data").Get("txTrieRoot").String()
				parentHash := r.Get("block_header.raw_data").Get("parentHash").String()
				coinAddr := r.Get("block_header.raw_data").Get("witness_address").String()
				blockTime := r.Get("block_header.raw_data").Get("timestamp").String()

				block.Id = uint64(time.Now().UnixNano())
				block.BlockHash = hash
				block.BlockNumber = fmt.Sprintf("%v", number)
				block.BlockTime = blockTime
				block.ParentHash = parentHash
				block.TxRoot = txRoot
				block.Coinbase = coinAddr

				array := r.Get("transactions.#.txID").Array()
				list := make([]string, 0, len(array))
				for _, v := range array {
					list = append(list, v.String())
				}
				block.Transactions = list

			}

			lock.Lock()
			list = append(list, &block)
			lock.Unlock()
		}
	}
}

func (s *StoreService) readReceiptFromKafka(blockChain int64, kafkaCfg map[string]*config.KafkaConfig) {
	receiver := make(chan *kafka.Message)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		Kafka := kafkaCfg["Receipt"]
		broker := fmt.Sprintf("%v:%v", Kafka.Host, Kafka.Port)
		group := fmt.Sprintf("group_store_receipt_%v", Kafka.Group)
		s.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: Kafka.Topic, Group: group, Partition: Kafka.Partition, StartOffset: Kafka.StartOffset}, receiver, ctx)
	}()

	list := make([]*service.Receipt, 0, 20)
	tk := time.NewTicker(5 * time.Second)
	lock := sync.RWMutex{}

	for true {
		select {
		case <-tk.C:
			lock.Lock()
			if len(list) > 0 {
				err := s.core.NewReceipt(blockChain, list)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
				}
				list = list[:0]
			}
			lock.Unlock()

		case msg := <-receiver:
			if msg.Value == nil || len(msg.Value) < 5 {
				continue
			}

			var receipt service.Receipt
			if blockChain == 200 {
				err := json.Unmarshal(msg.Value, &receipt)
				if err != nil {
					s.log.Errorf("readTxFromKafka|error=%v", err.Error())
					continue
				}
			} else if blockChain == 205 {

				root := gjson.ParseBytes(msg.Value)
				hash := root.Get("id").String()
				fee := root.Get("fee").Uint()
				number := root.Get("blockNumber").Uint()
				blockTime := root.Get("blockTimeStamp").Uint()
				contractAddress := root.Get("contract_address").String()
				status := root.Get("receipt.result").String()

				logs := root.Get("log").String()
				var list *service.Logs
				_ = json.Unmarshal([]byte(logs), &list)

				receipt.Id = uint64(time.Now().UnixNano())
				receipt.TransactionHash = hash
				receipt.BlockNumber = fmt.Sprintf("%v", number)
				receipt.GasUsed = fmt.Sprintf("%v", fee)
				receipt.Status = status
				receipt.ContractAddress = contractAddress
				receipt.CreateTime = fmt.Sprintf("%v", blockTime)
				receipt.Logs = list
				receipt.CumulativeGasUsed = root.Get("receipt").String()
			}

			lock.Lock()
			list = append(list, &receipt)
			lock.Unlock()
		}
	}
}
