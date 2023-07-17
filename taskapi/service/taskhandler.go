package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	kafkaClient "github.com/0xcregis/easynode/common/kafka"
	"github.com/0xcregis/easynode/taskapi/config"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
)

type TaskHandler struct {
	cfg         *config.Config
	kafkaClient *kafkaClient.EasyKafka
	sendCh      chan []*kafka.Message
}

func NewTaskHandler(cfg *config.Config, log *xlog.XLog) TaskApiInterface {
	kf := kafkaClient.NewEasyKafka(log)
	sendCh := make(chan []*kafka.Message, 10)

	m := &TaskHandler{
		cfg:         cfg,
		kafkaClient: kf,
		sendCh:      sendCh,
	}

	go func() {
		m.startKafka()
	}()
	return m
}

func (m *TaskHandler) startKafka() {
	broker := fmt.Sprintf("%v:%v", m.cfg.TaskKafka.Host, m.cfg.TaskKafka.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	m.kafkaClient.Write(kafkaClient.Config{Brokers: []string{broker}}, m.sendCh, nil, ctx)
}

func (m *TaskHandler) SendNodeTask(task *NodeTask) error {
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
