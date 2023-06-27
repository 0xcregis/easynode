package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/chain"
	"github.com/uduncloud/easynode/collect/service/db"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/common/util"
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

	kafkaCh chan []*kafka.Message
}

func NewService(cfg *config.Chain, logConfig *config.LogConfig) *Cmd {
	log := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/cmd", logConfig.Path), 24*time.Hour)
	txChan := make(chan *service.NodeTask, 10)
	receiptChan := make(chan *service.NodeTask, 10)
	blockChan := make(chan *service.NodeTask, 10)
	kafkaCh := make(chan []*kafka.Message, 10)
	store := db.NewTaskCacheService(cfg, log)
	kf := kafkaClient.NewEasyKafka(log)
	chainApi := getBlockChainApi(cfg, logConfig, store)

	return &Cmd{
		log:         log,
		taskStore:   store,
		chain:       cfg,
		kafka:       kf,
		blockChain:  chainApi,
		txChan:      txChan,
		receiptChan: receiptChan,
		blockChan:   blockChan,
		kafkaCh:     kafkaCh,
	}
}

func getBlockChainApi(c *config.Chain, logConfig *config.LogConfig, store service.StoreTaskInterface) service.BlockChainInterface {
	srv := chain.GetBlockchain(c.BlockChainCode, c, store, logConfig)
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
	//go func() {
	//	broker := fmt.Sprintf("%v:%v", c.chain.TaskKafka.Host, c.chain.TaskKafka.Port)
	//	c.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker}}, c.taskKafkaCh, nil)
	//}()

	go c.ReadNodeTaskFromKafka(nodeId, c.chain.BlockChainCode, c.blockChan, c.txChan, c.receiptChan)

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
	//go c.HandlerKafkaRespMessage(c.kafkaRespCh)

	//发送Kafka
	//tx,block,receipt 节点
	broker := fmt.Sprintf("%v:%v", c.chain.Kafka.Host, c.chain.Kafka.Port)
	//任务节点
	broker2 := fmt.Sprintf("%v:%v", c.chain.TaskKafka.Host, c.chain.TaskKafka.Port)

	c.kafka.WriteBatch(&kafkaClient.Config{Brokers: []string{broker, broker2}}, c.kafkaCh, c.HandlerKafkaRespMessage)
}

func (c *Cmd) ReadNodeTaskFromKafka(nodeId string, blockChain int, blockCh chan *service.NodeTask, txCh chan *service.NodeTask, receiptChan chan *service.NodeTask) {

	receiver := make(chan *kafka.Message)
	broker := fmt.Sprintf("%v:%v", c.chain.TaskKafka.Host, c.chain.TaskKafka.Port)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//taskKafka read
	go func() {
		group := fmt.Sprintf("gr_nodetask_%v_%v", blockChain, c.chain.TaskKafka.Group)
		topic := fmt.Sprintf("task_%v", blockChain)
		c.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: topic, Group: group, Partition: c.chain.TaskKafka.Partition, StartOffset: c.chain.TaskKafka.StartOffset}, receiver, ctx)
	}()

	log := c.log.WithFields(logrus.Fields{"model": "Execution", "id": time.Now().UnixMilli()})

	for true {
		select {
		case msg := <-receiver:
			//log.Printf("Read NodeTask offset=%v,topic=%v,value=%v", msg.Offset, msg.Topic, string(msg.Value))
			task := service.NodeTask{}
			err := json.Unmarshal(msg.Value, &task)
			if err != nil {
				log.Errorf("Read NodeTask offset=%v,error=%v", msg.Offset, err.Error())
				continue
			}

			if task.Id <= 0 || task.NodeId != nodeId || task.TaskStatus != 0 {
				log.Warnf("Read NodeTask offset=%v,task=%+v ,but it is wrong since id<=0 or nodeId !=local NodeId or status !=0 ", msg.Offset, task)
				continue
			}

			//交易任务
			{
				if task.TaskType == 1 { //单个交易
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyTx, task.BlockChain, task.TxHash), &task)
					log.Printf("Tx|blockchain:%v,txHash:%v", task.BlockChain, task.TxHash)
					txCh <- &task
				}

				if task.TaskType == 4 { //区块所有交易
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyTx, task.BlockChain, task.BlockHash+task.BlockNumber), &task)
					log.Printf("Tx:blockchain:%v,blockNumber:%v,blockHash:%v", task.BlockChain, task.BlockNumber, task.BlockHash)
					txCh <- &task
				}

			}

			//区块任务
			if task.TaskType == 2 {
				if len(task.BlockNumber) > 0 {
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyBlock, task.BlockChain, task.BlockNumber), &task)
				} else if len(task.BlockHash) > 0 {
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyBlock, task.BlockChain, task.BlockHash), &task)
				}
				log.Printf("Block:blockchain:%v,blockNumber:%v,blockHash:%v", task.BlockChain, task.BlockNumber, task.BlockHash)
				blockCh <- &task
			}

			//收据任务
			{
				if task.TaskType == 3 { //单收据
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyReceipt, task.BlockChain, task.TxHash), &task)
					receiptChan <- &task
				}

				if task.TaskType == 5 { //区块所有收据
					c.taskStore.StoreNodeTask(fmt.Sprintf(KeyReceipt, task.BlockChain, task.BlockHash+task.BlockNumber), &task)
					receiptChan <- &task
				}

			}
		}
	}

}

//func (c *Cmd) HandlerKafkaRespMessageEx(kafkaRespCh chan []*kafka.Message) {
//	buff := make(chan int64, 10)
//	for true {
//		buff <- 1
//		go func(kafkaRespCh chan []*kafka.Message) {
//			c.HandlerKafkaRespMessage(kafkaRespCh)
//			<-buff
//		}(kafkaRespCh)
//	}
//}

func (c *Cmd) HandlerKafkaRespMessage(msList []*kafka.Message) {

	//msList := <-kafkaRespCh
	ids := make([]string, 0, 20)
	log := c.log.WithFields(logrus.Fields{
		"id":    time.Now().UnixMilli(),
		"model": "HandlerKafkaRespMessage",
	})

	for _, msg := range msList {
		//交易消息回调，如下处理
		topic := msg.Topic
		//log.Printf("topic=%v,msg.offset=%v", topic, msg.Offset)

		//交易
		if c.chain.TxTask != nil && topic == c.chain.TxTask.Kafka.Topic {
			txHash := chain.GetTxHashFromKafka(c.chain.BlockChainCode, msg.Value)
			if len(txHash) > 0 {
				ids = append(ids, fmt.Sprintf(KeyTxById, c.chain.BlockChainCode, txHash))
			}
		}

		//收据
		if c.chain.ReceiptTask != nil && topic == c.chain.ReceiptTask.Kafka.Topic {
			txHash := chain.GetReceiptHashFromKafka(c.chain.BlockChainCode, msg.Value)
			if len(txHash) > 0 {
				ids = append(ids, fmt.Sprintf(KeyReceiptById, c.chain.BlockChainCode, txHash))
			}
		}

		//区块
		if c.chain.BlockTask != nil && topic == c.chain.BlockTask.Kafka.Topic {
			blockHash := chain.GetBlockHashFromKafka(c.chain.BlockChainCode, msg.Value)
			if len(blockHash) > 0 {
				ids = append(ids, fmt.Sprintf(KeyBlockById, c.chain.BlockChainCode, blockHash))
			}
		}

	}

	if len(ids) > 0 {
		err := c.taskStore.UpdateNodeTaskStatusWithBatch(ids, 1)
		if err != nil {
			log.Errorf("UpdateNodeTaskStatusWithBatch|err=%v", err)
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
			if taskReceipt.TaskType == 3 {
				list := c.execSingleReceipt(taskReceipt, log)
				if len(list) > 0 {
					kf <- list
				}
			}

			if taskReceipt.TaskType == 5 {
				list := c.execMultiReceipt(taskReceipt, log)
				if len(list) > 0 {
					kf <- list
				}

			}

		}(receiptTask, buffCh, kf, log)
	}

}

func (c *Cmd) execMultiReceipt(taskReceipt *service.NodeTask, log *logrus.Entry) []*kafka.Message {
	//根据区块
	key := fmt.Sprintf(KeyReceipt, taskReceipt.BlockChain, taskReceipt.BlockHash+taskReceipt.BlockNumber)
	_ = c.taskStore.UpdateNodeTaskStatus(key, 3)

	//读取区块
	receiptList := c.blockChain.GetReceiptByBlock(taskReceipt.BlockHash, taskReceipt.BlockNumber, c.chain.ReceiptTask, log)

	if receiptList == nil || len(receiptList) < 1 {
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		log.Errorf("GetReceiptByBlock|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "receipt task is empty", taskReceipt.Id)
		return nil
	}

	//receiptTaskList := make([]*service.NodeTask, 0, 10)
	kafkaMsgList := make([]*kafka.Message, 0, 10)

	for _, receipt := range receiptList {

		t := &service.NodeTask{NodeId: taskReceipt.NodeId, BlockNumber: taskReceipt.BlockNumber, BlockHash: taskReceipt.BlockHash, TxHash: receipt.TransactionHash, TaskType: 3, BlockChain: taskReceipt.BlockChain, TaskStatus: 3, Id: time.Now().UnixNano(), CreateTime: time.Now(), LogTime: time.Now()}
		//receiptTaskList = append(receiptTaskList, t)
		m, err := c.HandlerReceipt(receipt)
		if err != nil {
			log.Errorf("HandlerReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
			t.TaskStatus = 2 //失败
		} else {
			t.TaskStatus = 4
			kafkaMsgList = append(kafkaMsgList, m)
		}
		c.taskStore.StoreNodeTask(fmt.Sprintf(KeyReceiptById, c.chain.BlockChainCode, receipt.TransactionHash), t)
	}

	_ = c.taskStore.UpdateNodeTaskStatus(key, 1) //成功、

	if len(kafkaMsgList) > 0 {
		return kafkaMsgList
	} else {
		return nil
	}
}

func (c *Cmd) execSingleReceipt(taskReceipt *service.NodeTask, log *logrus.Entry) []*kafka.Message {
	if len(taskReceipt.TxHash) > 1 {
		key := fmt.Sprintf(KeyReceipt, taskReceipt.BlockChain, taskReceipt.TxHash)
		err := c.taskStore.UpdateNodeTaskStatus(key, 3)

		//根据交易
		receipt := c.blockChain.GetReceipt(taskReceipt.TxHash, c.chain.ReceiptTask, log)
		if receipt == nil {
			_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
			log.Errorf("GetReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "receipt is null", taskReceipt.Id)
			return nil
		}
		m, err := c.HandlerReceipt(receipt)
		if err != nil {
			_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
			log.Errorf("HandlerReceipt|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskReceipt.Id)
			return nil
		}

		//update status
		_ = c.taskStore.UpdateNodeTaskStatus(key, 4) //成功

		//reset key
		newKey := fmt.Sprintf(KeyReceiptById, taskReceipt.BlockChain, receipt.TransactionHash)
		_ = c.taskStore.ResetNodeTask(int64(taskReceipt.BlockChain), key, newKey)

		return []*kafka.Message{m}
	}

	return nil
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

				list := c.execSingleTx(taskTx, log)
				if len(list) > 0 {
					kf <- list
				}

			}(txTask, buffCh, kf, log)
		}

		if txTask.TaskType == 4 && (len(txTask.BlockHash) > 10 || len(txTask.BlockNumber) > 2) {
			//批量交易
			go func(taskTx *service.NodeTask, buffCh chan struct{}, kf chan []*kafka.Message, log *logrus.Entry) {
				defer func() {
					<-buffCh
				}()

				list := c.execMultiTx(taskTx, log)

				//批量发送到Kafka
				if len(list) > 0 {
					kf <- list
				}

			}(txTask, buffCh, kf, log)
		}

	}

}

func (c *Cmd) execMultiTx(taskTx *service.NodeTask, log *logrus.Entry) []*kafka.Message {

	start := time.Now()
	defer func() {
		log.Printf("execMultiTx,Duration=%v", time.Now().Sub(start))
	}()

	key := fmt.Sprintf(KeyTx, taskTx.BlockChain, taskTx.BlockHash+taskTx.BlockNumber)
	err := c.taskStore.UpdateNodeTaskStatus(key, 3)
	if err != nil {
		log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		return nil
	}

	var block *service.BlockInterface
	var txList []*service.TxInterface
	if len(taskTx.BlockNumber) > 5 {
		block, txList = c.blockChain.GetBlockByNumber(taskTx.BlockNumber, c.chain.BlockTask, log, true)
	} else if len(taskTx.BlockHash) > 10 {
		block, txList = c.blockChain.GetBlockByHash(taskTx.BlockHash, c.chain.BlockTask, log, true)
	}
	if block == nil {
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		log.Errorf("GetBlockByHashOrGetBlockByNumber｜BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskTx.Id)
		return nil
	}

	//发出tx
	txMessageList := make([]*kafka.Message, 0, len(txList))
	//txTaskList := make([]*service.NodeTask, 0, len(txList))
	receiptsTaskList := make([]*service.NodeTask, 0, len(txList))

	if txList != nil && len(txList) > 0 {

		for _, v := range txList {
			//交易任务
			t := &service.NodeTask{NodeId: taskTx.NodeId, BlockNumber: taskTx.BlockNumber, BlockHash: taskTx.BlockHash, TxHash: v.TxHash, TaskType: 1, BlockChain: taskTx.BlockChain, TaskStatus: 3, Id: time.Now().UnixNano(), CreateTime: time.Now(), LogTime: time.Now()}
			//txTaskList = append(txTaskList, t)

			//收据任务
			s := &service.NodeTask{NodeId: taskTx.NodeId, BlockChain: c.chain.BlockChainCode, TxHash: v.TxHash, BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, TaskType: 3, TaskStatus: 0, Id: time.Now().UnixNano(), CreateTime: time.Now(), LogTime: time.Now()}
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

			c.taskStore.StoreNodeTask(fmt.Sprintf(KeyTxById, c.chain.BlockChainCode, t.TxHash), t)
		}

		//未分配的收据任务
		//写入receipt task
		if len(receiptsTaskList) > 0 && c.chain.PullReceipt {

			//tron 不支持 批量获取收据
			//if c.chain.BlockChainCode == 205 {
			//c.taskKafkaCh <- c.taskStore.SendNodeTask(receiptsTaskList, c.chain.TaskKafka.Partitions)
			//} else {
			//	//其他公链，则 批量获取收据
			//	//根据区块获取收据，则 tx_hash 必需未空
			task := &service.NodeTask{BlockChain: c.chain.BlockChainCode, TxHash: "", BlockHash: block.BlockHash, BlockNumber: block.BlockNumber, TaskType: 5, TaskStatus: 0, NodeId: taskTx.NodeId}
			c.kafkaCh <- c.taskStore.SendNodeTask([]*service.NodeTask{task}, c.chain.TaskKafka.Partitions)
			//}

		}
	}

	//批量交易任务成功结束
	_ = c.taskStore.UpdateNodeTaskStatus(key, 1)

	//批量发送到Kafka
	if len(txMessageList) > 0 {
		return txMessageList
	} else {
		return nil
	}
}

func (c *Cmd) execSingleTx(taskTx *service.NodeTask, log *logrus.Entry) []*kafka.Message {

	start := time.Now()
	defer func() {
		log.Printf("execSingleTx,Duration=%v", time.Now().Sub(start))
	}()

	key := fmt.Sprintf(KeyTx, taskTx.BlockChain, taskTx.TxHash)

	err := c.taskStore.UpdateNodeTaskStatus(key, 3)
	if err != nil {
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		log.Errorf("UpdateNodeTaskStatus|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err, taskTx.Id)
		return nil
	}

	var tx *service.TxInterface
	if len(taskTx.TxHash) > 1 {
		tx = c.blockChain.GetTx(taskTx.TxHash, c.chain.TxTask, log)
	}
	if tx == nil {
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		log.Errorf("GetTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "tx is null", taskTx.Id)
		return nil
	}
	m, err := c.HandlerTx(tx)
	if err != nil {
		_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
		log.Errorf("HandlerTx|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, err.Error(), taskTx.Id)
		return nil
	}

	//写入receipt task
	if c.chain.PullReceipt {
		c.kafkaCh <- c.taskStore.SendNodeTask([]*service.NodeTask{{BlockChain: c.chain.BlockChainCode, TxHash: tx.TxHash, BlockHash: "", BlockNumber: "", TaskType: 3, TaskStatus: 0, NodeId: taskTx.NodeId}}, c.chain.TaskKafka.Partitions)
	}

	//update status
	_ = c.taskStore.UpdateNodeTaskStatus(key, 4) //成功

	//reset key
	newKey := fmt.Sprintf(KeyTxById, taskTx.BlockChain, tx.TxHash)
	_ = c.taskStore.ResetNodeTask(int64(taskTx.BlockChain), key, newKey)

	//kf <- []*kafka.Message{m}
	return []*kafka.Message{m}
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
			start := time.Now()
			defer func() {
				<-buffCh
				log.Printf("ExecBlockTask,Duration=%v", time.Now().Sub(start))
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

			var block *service.BlockInterface
			//var txList []*service.Tx
			if len(taskBlock.BlockNumber) > 5 {
				block, _ = c.blockChain.GetBlockByNumber(taskBlock.BlockNumber, c.chain.BlockTask, log, false)
			} else if len(taskBlock.BlockHash) > 10 {
				block, _ = c.blockChain.GetBlockByHash(taskBlock.BlockHash, c.chain.BlockTask, log, false)
			}

			if block == nil {
				_ = c.taskStore.UpdateNodeTaskStatus(key, 2) //失败
				log.Errorf("GetBlockByHashOrGetBlockByNumber|BlockChainName=%v,err=%v,taskId=%v", c.chain.BlockChainName, "block is null", taskBlock.Id)
				return
			}

			//发出批量交易任务
			if c.chain.PullTx {
				s := &service.NodeTask{NodeId: taskBlock.NodeId, BlockChain: taskBlock.BlockChain, BlockNumber: taskBlock.BlockNumber, BlockHash: taskBlock.BlockHash, TaskType: 4, TaskStatus: 0}
				c.kafkaCh <- c.taskStore.SendNodeTask([]*service.NodeTask{s}, c.chain.TaskKafka.Partitions)
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
			_ = c.taskStore.ResetNodeTask(int64(taskBlock.BlockChain), key, newKey)

			//区块发送到Kafka
			kf <- []*kafka.Message{m}
		}(blockTask, buffCh, kf, log)
	}

}

func (c *Cmd) HandlerBlock(block *service.BlockInterface) (*kafka.Message, error) {
	k := c.chain.BlockTask.Kafka
	var r []byte
	if v, ok := block.Block.(string); ok {
		r = []byte(v)
	} else if v, ok := block.Block.(*service.Block); ok {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r = b
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(block.BlockHash), Value: r}, nil
}

func (c *Cmd) HandlerTx(tx *service.TxInterface) (*kafka.Message, error) {
	k := c.chain.TxTask.Kafka

	var r []byte
	if v, ok := tx.Tx.(string); ok {
		r = []byte(v)
	} else if v, ok := tx.Tx.(*service.Tx); ok {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r = b
	} else if v, ok := tx.Tx.(map[string]interface{}); ok {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r = b
	}

	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(tx.TxHash), Value: r}, nil
}

func (c *Cmd) HandlerReceipt(receipt *service.ReceiptInterface) (*kafka.Message, error) {
	k := c.chain.ReceiptTask.Kafka
	var r []byte
	if v, ok := receipt.Receipt.(string); ok {
		r = []byte(v)
	} else if v, ok := receipt.Receipt.(*service.Receipt); ok {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r = b
	} else if v, ok := receipt.Receipt.(*service.TronReceipt); ok {
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		r = b
	}
	return &kafka.Message{Topic: k.Topic, Partition: k.Partition, Time: time.Now(), Key: []byte(receipt.TransactionHash), Value: r}, nil
}

func (c *Cmd) Stop() {
}
