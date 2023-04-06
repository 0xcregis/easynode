package kafka

import (
	"context"
	"errors"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/sunjiangjun/xlog"
	"time"
)

type Config struct {
	Brokers   []string
	Topic     string
	Group     string
	Partition int
}

type EasyKafka struct {
	log   *xlog.XLog
	query chan int
	Logf  func(string, ...interface{})
}

func NewEasyKafka(xLog *xlog.XLog) *EasyKafka {
	f := func(msg string, a ...interface{}) {
		xLog.Printf("kafka|msg=%v , other=%v", msg, a)
		fmt.Println(msg, a)
		fmt.Println()
	}
	query := make(chan int, 5)
	return &EasyKafka{
		log:   xLog,
		Logf:  f,
		query: query,
	}
}

func (easy *EasyKafka) Read(c *Config, ch chan *kafka.Message) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Brokers,
		Topic:       c.Topic,
		Partition:   c.Partition,
		GroupID:     c.Group,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		Logger:      kafka.LoggerFunc(easy.Logf),
		ErrorLogger: kafka.LoggerFunc(easy.Logf),
	})
	//r.SetOffset(0)

	errCh := make(chan int)

	go func(r *kafka.Reader, ch chan *kafka.Message, errCh chan int) {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				errCh <- 1
				break
			}
			easy.log.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			ch <- &m
		}
	}(r, ch, errCh)

	<-errCh
	if err := r.Close(); err != nil {
		easy.log.Printf("failed to close reader:%v", err)
	}
}

func (easy *EasyKafka) WriteBatch(c *Config, ch chan []*kafka.Message, resp chan []*kafka.Message) {
	for true {
		easy.query <- 1
		go func(c *Config, ch chan []*kafka.Message, resp chan []*kafka.Message) {
			easy.Write(*c, ch, resp)
		}(c, ch, resp)
	}
}

func (easy *EasyKafka) Write(c Config, ch chan []*kafka.Message, resp chan []*kafka.Message) {
	// make a writer that produces to topic-A, using the least-bytes distribution
	w := &kafka.Writer{
		Addr: kafka.TCP(c.Brokers...),
		//Topic:       c.Topic,
		Balancer:               &kafka.LeastBytes{},
		Logger:                 kafka.LoggerFunc(easy.Logf),
		ErrorLogger:            kafka.LoggerFunc(easy.Logf),
		BatchBytes:             70e6,
		BatchSize:              20,
		AllowAutoTopicCreation: true,
	}

	defer func() {
		if err := w.Close(); err != nil {
			easy.log.Fatal("failed to close writer:", err)
		}
		<-easy.query
	}()

	running := true

	for running {
		easy.log.Printf("kafka|Write|length=%v", len(ch))
		ms := <-ch
		step := 10
		for i := 0; i < len(ms); i += step {
			n := i + step
			if len(ms) <= n {
				n = len(ms)
			}
			err := easy.sendToKafka(w, ms[i:n], resp)
			if err != nil {
				running = false
				break
			}
		}

	}

}

func (easy *EasyKafka) sendToKafka(w *kafka.Writer, ms []*kafka.Message, resp chan []*kafka.Message) error {
	temp := make([]kafka.Message, 0, len(ms))
	for _, v := range ms {
		//easy.log.Printf("kafka|topic=%v p=%v,key=%v,value=%v", v.Topic, v.Partition, string(v.Key), string(v.Value))
		temp = append(temp, *v)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	const retries = 5
	for i := 0; i < retries; i++ {
		err := w.WriteMessages(ctx, temp...)
		if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		if err != nil {
			easy.log.Errorf("kafka|sendToKafka|failed to write messages:%v", err)
			return err
		}
		break
	}

	if resp != nil {
		resp <- ms
	}
	return nil
}
