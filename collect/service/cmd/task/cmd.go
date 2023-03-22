package task

import (
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	kafkaClient "github.com/uduncloud/easynode/collect/common/kafka"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/db"
	"github.com/uduncloud/easynode/collect/service/cmd/task/ether"
	"github.com/uduncloud/easynode/collect/service/cmd/task/tron"
	"github.com/uduncloud/easynode/collect/util"
	"time"
)

type Cmd struct {
	blockChain service.BlockChainInterface
	log        *xlog.XLog
	task       *db.Service
	chain      *config.Chain
	kafka      *kafkaClient.EasyKafka
}

func GetBlockChainService(c *config.Chain, db *config.TaskDb, sourceDb *config.SourceDb, logConfig *config.LogConfig) service.BlockChainInterface {
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/chain_info", logConfig.Path), 24*time.Hour)
	var srv service.BlockChainInterface
	if c.BlockChainCode == 200 {
		srv = ether.NewService(c, db, sourceDb, x)
	} else if c.BlockChainCode == 205 {
		srv = tron.NewService(c, db, sourceDb, x)
	}
	//第三方节点监控
	srv.Monitor()
	return srv
}

func NewService(c *config.Chain, taskDb *config.TaskDb, sourceDb *config.SourceDb, logConfig *config.LogConfig) *Cmd {
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/cmd", logConfig.Path), 24*time.Hour)
	t := db.NewTaskService(taskDb, sourceDb, x)
	client := kafkaClient.NewEasyKafka(x)

	block := GetBlockChainService(c, taskDb, sourceDb, logConfig)

	return &Cmd{
		log:        x,
		task:       t,
		chain:      c,
		kafka:      client,
		blockChain: block,
	}
}

func (c *Cmd) Start() {
	c.log.Printf("start %v service...", c.chain.BlockChainName)

	nodeId := util.GetLocalNodeId()
	txChan := make(chan *service.NodeTask, 10)
	receiptChan := make(chan *service.NodeTask, 10)
	blockChan := make(chan *service.NodeTask, 10)
	kafkaCh := make(chan []*kafka.Message, 100)
	kafkaRespCh := make(chan []*kafka.Message, 30)

	//配置了区块执行计划
	if c.chain.BlockTask != nil {
		//获取区块任务
		go c.getBlockTask(nodeId, c.chain.BlockChainCode, blockChan)

		//执行 exec_block 任务
		go c.ExecBlockTask(blockChan, kafkaCh)
	}

	//配置了 交易执行计划
	if c.chain.TxTask != nil {
		//获取交易任务
		go c.getTxTask(nodeId, c.chain.BlockChainCode, txChan)

		//执行交易
		go c.ExecTxTask(txChan, kafkaCh)
	}

	//配置了 收据执行计划
	if c.chain.ReceiptTask != nil {
		//获取收据任务
		go c.getReceiptTask(nodeId, c.chain.BlockChainCode, receiptChan)

		//执行收据
		go c.ExecReceiptTask(receiptChan, kafkaCh)
	}

	//kafka 消息回调
	go c.HandlerKafkaRespMessage(kafkaRespCh)

	//发送Kafka
	broker := fmt.Sprintf("%v:%v", c.chain.Kafka.Host, c.chain.Kafka.Port)
	c.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, kafkaCh, kafkaRespCh)
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

		ids := make([]int64, 0, 20)
		txHashList := make([]string, 0, 20)
		receiptsHashList := make([]string, 0, 20)
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "HandlerKafkaRespMessage",
		})

		for _, msg := range msList {

			topic := msg.Topic
			r := gjson.ParseBytes(msg.Value)

			log.Printf("topic=%v,msg.value=%v", topic, string(msg.Value))

			//交易
			if c.chain.TxTask != nil && topic == c.chain.TxTask.Kafka.Topic {

				txHash := r.Get("hash").String()

				txHashList = append(txHashList, txHash)
			}

			//收据
			if c.chain.ReceiptTask != nil && topic == c.chain.ReceiptTask.Kafka.Topic {

				txHash := r.Get("transactionHash").String()

				receiptsHashList = append(receiptsHashList, txHash)
			}

			//区块
			if c.chain.BlockTask != nil && topic == c.chain.BlockTask.Kafka.Topic {
				number := r.Get("number").String()
				blockHash := r.Get("hash").String()

				t, err := c.task.GetNodeTaskByBlockNumber(number, 2, c.chain.BlockChainCode)
				if err == nil && t.Id > 1 {
					ids = append(ids, t.Id)
					//_ = c.task.UpdateNodeTaskStatus(t.Id, 1)
				} else {
					t, err = c.task.GetNodeTaskByBlockHash(blockHash, 2, c.chain.BlockChainCode)
					if err == nil && t.Id > 1 {
						ids = append(ids, t.Id)
						//_ = c.task.UpdateNodeTaskStatus(t.Id, 1)
					}
				}
			}

		}

		if len(txHashList) > 0 {
			list, err := c.task.GetNodeTaskWithTxs(txHashList, 1, c.chain.BlockChainCode, 4)
			if err != nil {
				log.Errorf("GetNodeTaskWithTxs|tx|err=%v", err)
			} else {
				ids = append(ids, list...)
			}
		}

		if len(receiptsHashList) > 0 {
			list, err := c.task.GetNodeTaskWithTxs(receiptsHashList, 3, c.chain.BlockChainCode, 4)
			if err != nil {
				log.Errorf("GetNodeTaskWithTxs|receipt|err=%v", err)
			} else {
				ids = append(ids, list...)
			}
		}

		if len(ids) > 0 {
			err := c.task.UpdateNodeTaskStatusWithBatch(ids, 1)
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
		log.Printf("receiptTask=%+v", receiptTask)

		go func(taskReceipt *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
			defer func() {
				<-buffCh
			}()

			err := c.task.UpdateNodeTaskStatus(taskReceipt.Id, 3)
			if err != nil {
				err := c.task.UpdateNodeTaskStatus(taskReceipt.Id, 2)
				log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskReceipt.Id)
				return
			}

			//receiptList := make([]*service.Receipt, 0, 10)
			if len(taskReceipt.TxHash) > 1 {
				//根据交易
				receipt := c.blockChain.GetReceipt(taskReceipt.TxHash, c.chain.ReceiptTask, log)
				if receipt == nil {
					_ = c.task.UpdateNodeTaskStatus(taskReceipt.Id, 2) //失败
					log.Errorf("GetReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "receipt is null", taskReceipt.Id)
					return
				}
				m, err := c.HandlerReceipt(receipt)
				if err != nil {
					_ = c.task.UpdateNodeTaskStatus(taskReceipt.Id, 2) //失败
					log.Errorf("HandlerReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
					return
				}
				kf <- []*kafka.Message{m}
				_ = c.task.UpdateNodeTaskStatus(taskReceipt.Id, 4) //成功

			} else {
				//根据区块
				receiptList := c.blockChain.GetReceiptByBlock(taskReceipt.BlockHash, taskReceipt.BlockNumber, c.chain.ReceiptTask, log)

				if receiptList == nil || len(receiptList) < 1 {
					_ = c.task.UpdateNodeTaskStatus(taskReceipt.Id, 2) //失败
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
						continue
					} else {
						t.TaskStatus = 4
					}
					kafkaMsgList = append(kafkaMsgList, m)
				}

				if len(receiptTaskList) > 0 {
					err = c.task.AddNodeTask(receiptTaskList)
					if err != nil {
						log.Errorf("AddNodeTask|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
						return
					}

					_ = c.task.UpdateNodeTaskStatus(taskReceipt.Id, 1) //成功
					if len(kafkaMsgList) > 0 {
						kf <- kafkaMsgList
					}
				}

			}

		}(receiptTask, buffCh, kf, log)
	}

}

func (c *Cmd) ExecTxTask(txCh chan *service.NodeTask, kf chan []*kafka.Message) {
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

		log.Printf("txTask=%+v", txTask)

		if txTask.TaskType == 1 && len(txTask.TxHash) > 10 {
			//单个交易
			go func(taskTx *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
				defer func() {
					<-buffCh
				}()

				err := c.task.UpdateNodeTaskStatus(taskTx.Id, 3)
				if err != nil {
					_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 2) //失败
					log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
					return
				}

				var tx *service.Tx
				if len(taskTx.TxHash) > 1 {
					tx = c.blockChain.GetTx(taskTx.TxHash, c.chain.TxTask, log)
				}
				if tx == nil {
					_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 2) //失败
					log.Errorf("GetTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "tx is null", taskTx.Id)
					return
				}
				m, err := c.HandlerTx(tx)
				if err != nil {
					_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 2) //失败
					log.Errorf("HandlerTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
					return
				}

				//写入receipt task
				if c.chain.PullReceipt {
					_ = c.task.AddNodeSource(&service.NodeSource{BlockChain: c.chain.BlockChainCode, TxHash: tx.TxHash, BlockHash: tx.BlockHash, BlockNumber: tx.BlockNumber, SourceType: 3})
				}

				kf <- []*kafka.Message{m}
				_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 4) //成功

			}(txTask, buffCh, kf, log)

		} else if txTask.TaskType == 1 && (len(txTask.BlockHash) > 10 || len(txTask.BlockNumber) > 2) {
			//批量交易
			go func(taskTx *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
				defer func() {
					<-buffCh
				}()

				err := c.task.UpdateNodeTaskStatus(taskTx.Id, 3)
				if err != nil {
					log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
					_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 2) //失败
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
					_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 2) //失败
					log.Errorf("GetBlockByHashOrGetBlockByNumber｜BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskTx.Id)
					return
				}

				//发出tx
				if txList != nil && len(txList) > 0 {
					txMessageList := make([]*kafka.Message, 0, len(txList))
					txTaskList := make([]*service.NodeTask, 0, len(txList))
					nodeSourceList := make([]*service.NodeSource, 0, len(txList))

					for _, v := range txList {
						//交易任务
						t := &service.NodeTask{NodeId: taskTx.NodeId, BlockNumber: taskTx.BlockNumber, BlockHash: taskTx.BlockHash, TxHash: v.TxHash, TaskType: 1, BlockChain: taskTx.BlockChain, TaskStatus: 3}
						txTaskList = append(txTaskList, t)

						//收据任务
						s := &service.NodeSource{BlockChain: c.chain.BlockChainCode, TxHash: v.TxHash, BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, SourceType: 3}
						nodeSourceList = append(nodeSourceList, s)

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
					}

					//add source task: 未分配的收据任务
					//写入receipt task
					if len(txTaskList) > 0 && c.chain.PullReceipt {

						//tron 不支持 批量获取收据
						if c.chain.BlockChainCode == 205 {
							_ = c.task.AddNodeSourceList(nodeSourceList)
						} else {
							//其他公链，则 批量获取收据
							//根据区块获取收据，则 tx_hash 必需未空
							ns := &service.NodeSource{BlockChain: c.chain.BlockChainCode, TxHash: "", BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, SourceType: 3}
							err := c.task.AddNodeSource(ns)
							if err != nil {
								log.Errorf("AddNodeSource|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
							}
						}

					}

					//更新交易状态
					if len(txTaskList) > 0 {
						_ = c.task.AddNodeTask(txTaskList)
					}

					//批量发送到Kafka
					if len(txMessageList) > 0 {
						kf <- txMessageList
					}
				}

				//批量交易任务成功结束
				_ = c.task.UpdateNodeTaskStatus(taskTx.Id, 1)

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
		c.log.Printf("taskBlock=%+v", blockTask)

		go func(taskBlock *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
			defer func() {
				<-buffCh
			}()

			err := c.task.UpdateNodeTaskStatus(taskBlock.Id, 3)
			if err != nil {
				log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskBlock.Id)
				_ = c.task.UpdateNodeTaskStatus(taskBlock.Id, 2) //失败
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
				_ = c.task.UpdateNodeTaskStatus(taskBlock.Id, 2) //失败
				log.Errorf("GetBlockByHashOrGetBlockByNumber|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskBlock.Id)
				return
			}

			//发出批量交易任务
			if c.chain.PullTx {
				s := &service.NodeSource{BlockChain: taskBlock.BlockChain, BlockNumber: taskBlock.BlockNumber, BlockHash: taskBlock.BlockHash, SourceType: 1}
				err = c.task.AddNodeSource(s)
				if err != nil {
					_ = c.task.UpdateNodeTaskStatus(taskBlock.Id, 2) //失败
					log.Errorf("AddNodeSource|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "push task is error with batch tx by block", taskBlock.Id)
					return
				}
			}

			//发出block
			m, err := c.HandlerBlock(block)
			if err != nil {
				_ = c.task.UpdateNodeTaskStatus(taskBlock.Id, 2) //失败
				log.Errorf("HandlerBlock|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskBlock.Id)
				return
			}

			err = c.task.UpdateNodeTaskStatus(taskBlock.Id, 4) //成功
			if err != nil {
				log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskBlock.Id)
				return
			}

			//区块发送到Kafka
			kf <- []*kafka.Message{m}

		}(blockTask, buffCh, kf, log)
	}

}

func (c *Cmd) getReceiptTask(nodeId string, blockChain int, receiptChan chan *service.NodeTask) {
	for true {
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "GetReceiptTask",
		})
		log.Printf("nodeId=%v ,blockChain=%v", nodeId, blockChain)
		var list []*service.NodeTask
		var err error
		//if c.lock() {
		list, err = c.task.GetTaskWithReceipt(blockChain, nodeId)
		if err != nil {
			log.Errorf("GetTaskWithReceipt|error=%v", err)
		}
		//c.unlock()
		//}
		log.Printf("ReceiptsTask.list=%v", len(list))
		if len(list) > 0 {
			for _, v := range list {
				receiptChan <- v
			}
		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Cmd) getTxTask(nodeId string, blockChain int, txCh chan *service.NodeTask) {
	for true {
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "GetTxTask",
		})
		log.Printf("nodeId=%v ,blockChain=%v", nodeId, blockChain)
		var list []*service.NodeTask
		var err error
		//if c.lock() {
		list, err = c.task.GetTaskWithTx(blockChain, nodeId)
		if err != nil {
			log.Errorf("GetTaskWithTx|err=%v", err)
		}
		//c.unlock()
		//}
		log.Printf("TxTask.list=%v", len(list))
		if len(list) > 0 {
			for _, v := range list {
				txCh <- v
			}

		} else {
			time.Sleep(5 * time.Second)
		}
	}
}

func (c *Cmd) getBlockTask(nodeId string, blockChain int, blockCh chan *service.NodeTask) {
	for true {
		log := c.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "GetBlockTask",
		})
		log.Printf("nodeId=%v ,blockChain=%v", nodeId, blockChain)
		var list []*service.NodeTask
		var err error
		//if c.lock() {
		list, err = c.task.GetTaskWithBlock(blockChain, nodeId)
		if err != nil {
			log.Errorf("GetTaskWithBlock|err=%v", err)
		}
		//c.unlock()
		//}
		log.Printf("BlockTask.list=%v", len(list))

		if len(list) > 0 {
			for _, v := range list {
				blockCh <- v
			}
		} else {
			time.Sleep(3 * time.Second)
		}
	}
}

func (c *Cmd) HandlerBlock(block *service.Block) (*kafka.Message, error) {
	k := c.chain.BlockTask.Kafka
	nodeId := util.GetLocalNodeId()
	b, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(nodeId), Value: b}, nil
}

func (c *Cmd) HandlerTx(tx *service.Tx) (*kafka.Message, error) {
	k := c.chain.TxTask.Kafka
	nodeId := util.GetLocalNodeId()
	b, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(nodeId), Value: b}, nil
}

func (c *Cmd) HandlerReceipt(receipt *service.Receipt) (*kafka.Message, error) {
	k := c.chain.ReceiptTask.Kafka
	nodeId := util.GetLocalNodeId()
	b, err := json.Marshal(receipt)
	if err != nil {
		return nil, err
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(nodeId), Value: b}, nil
}

func (c *Cmd) Stop() {
}
