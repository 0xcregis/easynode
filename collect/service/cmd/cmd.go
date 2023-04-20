package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/db"
	"github.com/uduncloud/easynode/collect/service/cmd/task/ether"
	"github.com/uduncloud/easynode/collect/service/cmd/task/tron"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/common/util"
	"math/rand"
	"strings"
	"time"
)

const (
	KeyBlock       = "%v_block_%v"
	KeyBlockById   = "%v_blockId_%v"
	KeyTx          = "%v_tx_%v"
	KeyTxById      = "%v_txId_%v"
	KeyReceipt     = "%v_receipt_%v"
	KeyReceiptById = "%v_receiptId_%v"
)

type Cmd struct {
	blockChain  service.BlockChainInterface
	log         *xlog.XLog
	taskStore   service.StoreTaskInterface
	chain       *config.Chain
	kafka       *kafkaClient.EasyKafka
	txChan      chan *service.NodeTask
	blockChan   chan *service.NodeTask
	receiptChan chan *service.NodeTask
	kafkaCh     chan []*kafka.Message
	kafkaRespCh chan []*kafka.Message

	taskKafkaCh chan []*kafka.Message
}

func NewService(cfg *config.Chain, logConfig *config.LogConfig) *Cmd {
	log := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/cmd", logConfig.Path), 24*time.Hour)
	txChan := make(chan *service.NodeTask, 10)
	receiptChan := make(chan *service.NodeTask, 10)
	blockChan := make(chan *service.NodeTask, 10)
	kafkaCh := make(chan []*kafka.Message, 100)
	kafkaRespCh := make(chan []*kafka.Message, 30)
	taskKafkaCh := make(chan []*kafka.Message, 10)
	store := db.NewTaskCacheService(cfg, taskKafkaCh, log)
	kf := kafkaClient.NewEasyKafka(log)
	chain := getBlockChainService(cfg, logConfig)

	return &Cmd{
		log:         log,
		taskStore:   store,
		chain:       cfg,
		kafka:       kf,
		blockChain:  chain,
		txChan:      txChan,
		receiptChan: receiptChan,
		blockChan:   blockChan,
		kafkaCh:     kafkaCh,
		kafkaRespCh: kafkaRespCh,
		taskKafkaCh: taskKafkaCh,
	}
}

func getBlockChainService(c *config.Chain, logConfig *config.LogConfig) service.BlockChainInterface {
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/chain_info", logConfig.Path), 24*time.Hour)
	var srv service.BlockChainInterface
	if c.BlockChainCode == 200 {
		srv = ether.NewService(c, x)
	} else if c.BlockChainCode == 205 {
		srv = tron.NewService(c, x)
	}
	//第三方节点监控
	srv.Monitor()
	return srv
}

func (c *Cmd) Start() {
	c.log.Printf("start %v service...", c.chain.BlockChainName)

	nodeId, err := util.GetLocalNodeId()
	if err != nil {
		panic(err)
	}

	//task kafka write
	go func() {
		broker := fmt.Sprintf("%v:%v", c.chain.TaskKafka.Host, c.chain.TaskKafka.Port)
		//taskKafka write
		go func() {
			c.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, c.taskKafkaCh, nil)
		}()
	}()

	go c.HandlerNodeTaskFromKafka(nodeId, c.chain.BlockChainCode, c.blockChan, c.txChan, c.receiptChan)

	//配置了区块执行计划
	if c.chain.BlockTask != nil {
		//执行 exec_block 任务
		go c.ExecBlockTask(c.blockChan, c.kafkaCh)
	}

	//配置了 交易执行计划
	if c.chain.TxTask != nil {
		//执行交易
		go c.ExecTxTask(nodeId, c.txChan, c.kafkaCh)
	}

	//配置了 收据执行计划
	if c.chain.ReceiptTask != nil {
		//执行收据
		go c.ExecReceiptTask(c.receiptChan, c.kafkaCh)
	}

	//kafka 消息回调
	go c.HandlerKafkaRespMessage(c.kafkaRespCh)

	//发送Kafka
	broker := fmt.Sprintf("%v:%v", c.chain.Kafka.Host, c.chain.Kafka.Port)
	c.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, c.kafkaCh, c.kafkaRespCh)
}

func (c *Cmd) HandlerNodeTaskFromKafka(nodeId string, blockChain int, blockCh chan *service.NodeTask, txCh chan *service.NodeTask, receiptChan chan *service.NodeTask) {

	receiver := make(chan *kafka.Message)
	broker := fmt.Sprintf("%v:%v", c.chain.TaskKafka.Host, c.chain.TaskKafka.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//taskKafka read
	go func() {
		group := fmt.Sprintf("group_%v_%v", blockChain, rand.Intn(1000))
		topic := fmt.Sprintf("task_%v", blockChain)
		c.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: topic, Group: group, Partition: 0, StartOffset: -1}, receiver, ctx)
	}()

	for true {
		select {
		case msg := <-receiver:
			c.log.Printf("Read NodeTask offset=%v,topic=%v,value=%v", msg.Offset, msg.Topic, string(msg.Value))
			task := service.NodeTask{}
			err := json.Unmarshal(msg.Value, &task)
			if err != nil {
				c.log.Errorf("Read NodeTask offset=%v,error=%v", msg.Offset, err.Error())
				continue
			}

			if task.Id > 0 && task.NodeId == nodeId && task.TaskType == 1 && task.TaskStatus == 0 {
				//交易任务
				if len(task.TxHash) > 0 {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyTx, task.BlockChain, task.TxHash), &task)
				} else {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyTx, task.BlockChain, task.BlockHash+task.BlockNumber), &task)
				}
				txCh <- &task
			}

			if task.Id > 0 && task.NodeId == nodeId && task.TaskType == 2 && task.TaskStatus == 0 {
				//区块任务
				if len(task.BlockNumber) > 0 {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyBlock, task.BlockChain, task.BlockNumber), &task)
				} else if len(task.BlockHash) > 0 {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyBlock, task.BlockChain, task.BlockHash), &task)
				}
				blockCh <- &task
			}

			if task.Id > 0 && task.NodeId == nodeId && task.TaskType == 3 && task.TaskStatus == 0 {
				//收据任务
				if len(task.TxHash) > 0 {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyReceipt, task.BlockChain, task.TxHash), &task)
				} else {
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyReceipt, task.BlockChain, task.BlockHash+task.BlockNumber), &task)
				}
				receiptChan <- &task
			}
		}
	}

}

func (c *Cmd) HandlerKafkaRespMessageEx(kafkaRespCh chan []*kafka.Message) {
	buff := make(chan int64, 10)
	for true {
		buff <- 1
		go func(kafkaRespCh chan []*kafka.Message) {
			c.HandlerKafkaRespMessage(kafkaRespCh)
			<-buff
		}(kafkaRespCh)
	}
}

func (c *Cmd) HandlerKafkaRespMessage(kafkaRespCh chan []*kafka.Message) {

	for true {

		msList := <-kafkaRespCh
		ids := make([]string, 0, 20)
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "HandlerKafkaRespMessage",
		})

		for _, msg := range msList {
			topic := msg.Topic

			if strings.HasPrefix(topic, "task") {
				//任务消息回调，暂不做处理
				continue
			}

			//交易消息回调，如下处理
			r := gjson.ParseBytes(msg.Value)
			//log.Printf("topic=%v,msg.value=%v", topic, string(msg.Value))

			//交易
			if c.chain.TxTask != nil && topic == c.chain.TxTask.Kafka.Topic {
				txHash := r.Get("hash").String()
				ids = append(ids, fmt.Sprintf(KeyTxById, c.chain.BlockChainCode, txHash))
			}

			//收据
			if c.chain.ReceiptTask != nil && topic == c.chain.ReceiptTask.Kafka.Topic {
				txHash := r.Get("transactionHash").String()
				ids = append(ids, fmt.Sprintf(KeyReceiptById, c.chain.BlockChainCode, txHash))
			}

			//区块
			if c.chain.BlockTask != nil && topic == c.chain.BlockTask.Kafka.Topic {
				blockHash := r.Get("hash").String()
				ids = append(ids, fmt.Sprintf(KeyBlock, c.chain.BlockChainCode, blockHash))
			}

		}

		if len(ids) > 0 {
			err := c.taskStore.UpdateNodeTaskStatusWithBatch(ids, 1)
			if err != nil {
				log.Errorf("UpdateNodeTaskStatusWithBatch|err=%v", err)
			}
		}
	}

}

func (c *Cmd) ExecReceiptTask(receiptChan chan *service.NodeTask, kf chan []*kafka.Message) {
	buffCh := make(chan struct{}, 5)
	for true {
		//控制协程数量
		buffCh <- struct{}{}
		//执行任务
		receiptTask := <-receiptChan
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "ExecReceiptTask",
		})
		//log.Printf("receiptTask=%+v", receiptTask)

		go func(taskReceipt *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
			defer func() {
				<-buffCh
			}()

			//receiptList := make([]*service.Receipt, 0, 10)
			if len(taskReceipt.TxHash) > 1 {
				key := fmt.Sprintf(KeyReceipt, taskReceipt.BlockChain, taskReceipt.TxHash)
				err := c.taskStore.UpdateNodeTaskStatus(key, 3)

				//根据交易
				receipt := c.blockChain.GetReceipt(taskReceipt.TxHash, c.chain.ReceiptTask, log)
				if receipt == nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("GetReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "receipt is null", taskReceipt.Id)
					return
				}
				m, err := c.HandlerReceipt(receipt)
				if err != nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("HandlerReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
					return
				}

				//update status
				_ = c.taskStore.UpdateNodeTaskStatus(key, 4) //成功

				//reset key
				newKey := fmt.Sprintf(KeyReceiptById, taskReceipt.BlockChain, receipt.TransactionHash)
				_ = c.taskStore.ResetNodeTask(key, newKey)

				kf <- []*kafka.Message{m}

			} else {

				//根据区块
				key := fmt.Sprintf(KeyReceipt, taskReceipt.BlockChain, taskReceipt.BlockHash+taskReceipt.BlockNumber)
				_ = c.taskStore.UpdateNodeTaskStatus(key, 3)

				//读取区块
				receiptList := c.blockChain.GetReceiptByBlock(taskReceipt.BlockHash, taskReceipt.BlockNumber, c.chain.ReceiptTask, log)

				if receiptList == nil || len(receiptList) < 1 {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "receipt task is empty", taskReceipt.Id)
					return
				}

				receiptTaskList := make([]*service.NodeTask, 0, 10)
				kafkaMsgList := make([]*kafka.Message, 0, 10)

				for _, receipt := range receiptList {

					t := &service.NodeTask{NodeId: taskReceipt.NodeId, BlockNumber: taskReceipt.BlockNumber, BlockHash: taskReceipt.BlockHash, TxHash: receipt.TransactionHash, TaskType: 3, BlockChain: receiptTask.BlockChain, TaskStatus: 3}
					receiptTaskList = append(receiptTaskList, t)
					m, err := c.HandlerReceipt(receipt)
					if err != nil {
						log.Errorf("HandlerReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
						t.TaskStatus = 2 //失败
						//continue
					} else {
						t.TaskStatus = 4
					}
					kafkaMsgList = append(kafkaMsgList, m)
					c.taskStore.StoreExecTask(fmt.Sprintf(KeyReceiptById, c.chain.BlockChainCode, receipt.TransactionHash), t)
				}

				if len(kafkaMsgList) > 0 {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 1) //成功、
					if len(kafkaMsgList) > 0 {
						kf <- kafkaMsgList
					}
				}

			}

		}(receiptTask, buffCh, kf, log)
	}

}

func (c *Cmd) ExecTxTask(nodeId string, txCh chan *service.NodeTask, kf chan []*kafka.Message) {
	buffCh := make(chan struct{}, 5)
	for true {
		//控制协程数量
		buffCh <- struct{}{}
		//执行任务
		txTask := <-txCh
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "ExecTxTask",
		})

		//log.Printf("txTask=%+v", txTask)

		if txTask.TaskType == 1 && len(txTask.TxHash) > 10 {
			//单个交易
			go func(taskTx *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
				defer func() {
					<-buffCh
				}()

				key := fmt.Sprintf(KeyTx, txTask.BlockChain, txTask.TxHash)

				err := c.taskStore.UpdateNodeTaskStatus(key, 3)
				if err != nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
					return
				}

				var tx *service.Tx
				if len(taskTx.TxHash) > 1 {
					tx = c.blockChain.GetTx(taskTx.TxHash, c.chain.TxTask, log)
				}
				if tx == nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("GetTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "tx is null", taskTx.Id)
					return
				}
				m, err := c.HandlerTx(tx)
				if err != nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("HandlerTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
					return
				}

				//写入receipt task
				if c.chain.PullReceipt {
					_ = c.taskStore.SendNodeTask([]*service.NodeTask{{BlockChain: c.chain.BlockChainCode, TxHash: tx.TxHash, BlockHash: tx.BlockHash, BlockNumber: tx.BlockNumber, TaskType: 3, TaskStatus: 0, NodeId: taskTx.NodeId}})
				}

				//update status
				_ = c.taskStore.UpdateNodeTaskStatus(key, 4) //成功

				//reset key
				newKey := fmt.Sprintf(KeyTxById, txTask.BlockChain, tx.TxHash)
				_ = c.taskStore.ResetNodeTask(key, newKey)

				kf <- []*kafka.Message{m}

			}(txTask, buffCh, kf, log)

		} else if txTask.TaskType == 1 && (len(txTask.BlockHash) > 10 || len(txTask.BlockNumber) > 2) {
			//批量交易
			go func(taskTx *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
				defer func() {
					<-buffCh
				}()
				key := fmt.Sprintf(KeyTx, txTask.BlockChain, txTask.BlockHash+txTask.BlockNumber)
				err := c.taskStore.UpdateNodeTaskStatus(key, 3)
				if err != nil {
					log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					return
				}

				var block *service.Block
				var txList []*service.Tx
				if len(taskTx.BlockHash) > 10 {
					block, txList = c.blockChain.GetBlockByHash(taskTx.BlockHash, c.chain.BlockTask, log)
				} else if len(taskTx.BlockNumber) > 5 {
					block, txList = c.blockChain.GetBlockByNumber(taskTx.BlockNumber, c.chain.BlockTask, log)
				}
				if block == nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("GetBlockByHashOrGetBlockByNumber｜BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskTx.Id)
					return
				}

				//发出tx
				if txList != nil && len(txList) > 0 {
					txMessageList := make([]*kafka.Message, 0, len(txList))
					//txTaskList := make([]*service.NodeTask, 0, len(txList))
					receiptsTaskList := make([]*service.NodeTask, 0, len(txList))

					for _, v := range txList {
						//交易任务
						t := &service.NodeTask{NodeId: taskTx.NodeId, BlockNumber: taskTx.BlockNumber, BlockHash: taskTx.BlockHash, TxHash: v.TxHash, TaskType: 1, BlockChain: taskTx.BlockChain, TaskStatus: 3}
						//txTaskList = append(txTaskList, t)

						//收据任务
						s := &service.NodeTask{NodeId: taskTx.NodeId, BlockChain: c.chain.BlockChainCode, TxHash: v.TxHash, BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, TaskType: 3, TaskStatus: 0}
						receiptsTaskList = append(receiptsTaskList, s)

						//交易消息
						m, err := c.HandlerTx(v)
						if err != nil {
							//更改交易状态：失败
							t.TaskStatus = 2
							log.Warnf("HandlerTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
						} else {
							//kf <- m
							//写入Kafka
							t.TaskStatus = 4
							txMessageList = append(txMessageList, m)
						}

						c.taskStore.StoreExecTask(fmt.Sprintf(KeyTxById, c.chain.BlockChainCode, t.TxHash), t)
					}

					//未分配的收据任务
					//写入receipt task
					if len(receiptsTaskList) > 0 && c.chain.PullReceipt {

						//tron 不支持 批量获取收据
						if c.chain.BlockChainCode == 205 {
							_ = c.taskStore.SendNodeTask(receiptsTaskList)
						} else {
							//其他公链，则 批量获取收据
							//根据区块获取收据，则 tx_hash 必需未空
							task := &service.NodeTask{BlockChain: c.chain.BlockChainCode, TxHash: "", BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, TaskType: 3, TaskStatus: 0, NodeId: taskTx.NodeId}
							err := c.taskStore.SendNodeTask([]*service.NodeTask{task})
							if err != nil {
								log.Errorf("AddNodeSource|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
							}
						}

					}

					//批量发送到Kafka
					if len(txMessageList) > 0 {
						kf <- txMessageList
					}
				}

				//批量交易任务成功结束
				_ = c.taskStore.UpdateNodeTaskStatus(key, 1)

			}(txTask, buffCh, kf, log)
		}

	}

}

func (c *Cmd) ExecBlockTask(blockCh chan *service.NodeTask, kf chan []*kafka.Message) {
	buffCh := make(chan struct{}, 5)
	for true {
		//控制协程数量
		buffCh <- struct{}{}

		//执行任务
		blockTask := <-blockCh
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "ExecBlockTask",
		})
		//log.Printf("taskBlock=%+v", blockTask)

		go func(taskBlock *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
			defer func() {
				<-buffCh
			}()

			var key string
			if len(taskBlock.BlockNumber) > 0 {
				key = fmt.Sprintf(KeyBlock, taskBlock.BlockChain, taskBlock.BlockNumber)
			} else if len(taskBlock.BlockHash) > 0 {
				key = fmt.Sprintf(KeyBlock, taskBlock.BlockChain, taskBlock.BlockHash)
			}

			err := c.taskStore.UpdateNodeTaskStatus(key, 3)
			if err != nil {
				log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskBlock.Id)
				_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
				return
			}

			var block *service.Block
			//var txList []*service.Tx
			if len(taskBlock.BlockHash) > 10 {
				block, _ = c.blockChain.GetBlockByHash(taskBlock.BlockHash, c.chain.BlockTask, log)
			} else if len(taskBlock.BlockNumber) > 5 {
				block, _ = c.blockChain.GetBlockByNumber(taskBlock.BlockNumber, c.chain.BlockTask, log)
			}
			if block == nil {
				_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
				log.Errorf("GetBlockByHashOrGetBlockByNumber|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskBlock.Id)
				return
			}

			//发出批量交易任务
			if c.chain.PullTx {
				s := &service.NodeTask{NodeId: taskBlock.NodeId, BlockChain: taskBlock.BlockChain, BlockNumber: taskBlock.BlockNumber, BlockHash: taskBlock.BlockHash, TaskType: 1, TaskStatus: 0}
				err = c.taskStore.SendNodeTask([]*service.NodeTask{s})
				if err != nil {
					_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
					log.Errorf("AddNodeSource|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "push task is error with batch tx by block", taskBlock.Id)
					return
				}
			}

			//发出block
			m, err := c.HandlerBlock(block)
			if err != nil {
				_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
				log.Errorf("HandlerBlock|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskBlock.Id)
				return
			}

			err = c.taskStore.UpdateNodeTaskStatus(key, 4) //成功
			if err != nil {
				log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskBlock.Id)
				return
			}

			//重置缓存key
			newKey := fmt.Sprintf(KeyBlockById, taskBlock.BlockChain, block.BlockHash)
			_ = c.taskStore.ResetNodeTask(key, newKey)

			//区块发送到Kafka
			kf <- []*kafka.Message{m}
		}(blockTask, buffCh, kf, log)
	}

}

func (c *Cmd) HandlerBlock(block *service.Block) (*kafka.Message, error) {
	k := c.chain.BlockTask.Kafka
	b, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(block.BlockHash), Value: b}, nil
}

func (c *Cmd) HandlerTx(tx *service.Tx) (*kafka.Message, error) {
	k := c.chain.TxTask.Kafka
	b, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(tx.TxHash), Value: b}, nil
}

func (c *Cmd) HandlerReceipt(receipt *service.Receipt) (*kafka.Message, error) {
	k := c.chain.ReceiptTask.Kafka
	b, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(receipt.TransactionHash), Value: b}, nil
}

func (c *Cmd) Stop() {
}
