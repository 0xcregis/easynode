package kafka

import (
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"log"
	"testing"
	"time"
)

func Init() *EasyKafka {
	x := xlog.NewXLogger().BuildOutType(xlog.STD)
	return NewEasyKafka(x)
}
func TestEasyKafka_Write(t *testing.T) {

	k := Init()
	ch := make(chan *kafka.Message, 10)
	go func() {
		for true {
			msg := &kafka.Message{Key: []byte("901"), Value: []byte(`{"name":"sunhongtao","age":11}`), Topic: "test2", Partition: 1}
			select {
			case ch <- msg:
				log.Printf("==  %+v\n", msg)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	k.Write(&Config{Brokers: []string{"192.168.2.20:9092"}, Topic: "test", Partition: 1, Group: "g1"}, ch)
}

func TestEasyKafka_Read(t *testing.T) {
	k := Init()
	ch := make(chan *kafka.Message, 10)

	go func() {
		for true {
			select {
			case msg := <-ch:
				log.Printf("==  %+v\n", msg)
			}
		}
	}()

	k.Read(&Config{Brokers: []string{"192.168.2.20:9092"}, Topic: "test2", Partition: 0, Group: "g1"}, ch)

}
