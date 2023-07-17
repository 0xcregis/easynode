package kafka

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
)

func Init() *EasyKafka {
	x := xlog.NewXLogger().BuildOutType(xlog.STD)
	return NewEasyKafka(x)
}
func TestEasyKafka_Write(t *testing.T) {

	k := Init()
	ch := make(chan []*kafka.Message, 10)
	go func() {
		for {
			msg := &kafka.Message{Key: []byte("901"), Value: []byte(`{"name":"sunhongtao","age":11}`), Topic: "test2", Partition: 0}

			ch <- []*kafka.Message{msg}
			log.Printf("==  %+v\n", msg)

			time.Sleep(3 * time.Second)
		}
	}()

	k.Write(Config{Brokers: []string{"kafka:9092"}, Group: "g1"}, ch, nil, context.Background())
}

func TestEasyKafka_Read(t *testing.T) {
	k := Init()
	ch := make(chan *kafka.Message, 10)

	go func() {
		for {
			msg := <-ch
			log.Printf("==  %+v\n", msg)
		}
	}()

	k.Read(&Config{Brokers: []string{"192.168.2.20:9092"}, Topic: "test2", Partition: 0, Group: "g1"}, ch, context.Background())
}
