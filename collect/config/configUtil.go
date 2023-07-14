package config

import (
	"encoding/json"
	"io"
	"os"
)

func LoadConfig(path string) Config {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	_ = json.Unmarshal(b, &cfg)
	return cfg
}

func (k *Kafka) CopyKafka() Kafka {
	//var ka Kafka
	ka := *k
	return ka
}

func (f *FromCluster) CopyFromCluster() FromCluster {
	//var fc FromCluster
	fc := *f
	return fc
}

func (b *BlockTask) CopyBlockTask() BlockTask {
	var task BlockTask
	task.Kafka = b.Kafka
	for _, v := range b.FromCluster {
		c := v.CopyFromCluster()
		task.FromCluster = append(task.FromCluster, &c)
	}
	return task
}

func (b *TxTask) CopyTxTask() TxTask {
	var task TxTask
	task.Kafka = b.Kafka
	for _, v := range b.FromCluster {
		c := v.CopyFromCluster()
		task.FromCluster = append(task.FromCluster, &c)
	}
	return task
}

func (b *ReceiptTask) CopyReceiptTask() ReceiptTask {
	var task ReceiptTask
	task.Kafka = b.Kafka
	for _, v := range b.FromCluster {
		c := v.CopyFromCluster()
		task.FromCluster = append(task.FromCluster, &c)
	}
	return task
}

func (c *Chain) CopyChain() Chain {
	chain := *c
	k := c.Kafka.CopyKafka()
	chain.Kafka = &k
	if c.BlockTask != nil {
		b := c.BlockTask.CopyBlockTask()
		chain.BlockTask = &b
	}

	if c.TxTask != nil {
		t := c.TxTask.CopyTxTask()
		chain.TxTask = &t
	}

	if c.ReceiptTask != nil {
		r := c.ReceiptTask.CopyReceiptTask()
		chain.ReceiptTask = &r
	}
	return chain
}
