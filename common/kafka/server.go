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
	Brokers     []string
	Topic       string
	Group       string
	Partition   int
	StartOffset int64
}

type EasyKafka struct {
	log  *xlog.XLog
	Logf func(string, ...interface{})
}

func NewEasyKafka(xLog *xlog.XLog) *EasyKafka {
	f := func(msg string, a ...interface{}) {
		//xLog.Printf("kafka|msg=%v , other=%v", msg, a)
		fmt.Println(msg, a)
	}
	return &EasyKafka{
		log:  xLog,
		Logf: f,
	}
}

func (easy *EasyKafka) Read(c *Config, ch chan *kafka.Message, ctx context.Context) {
	//if c.StartOffset == 0 {
	//	c.StartOffset = kafka.LastOffset
	//}
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     c.Brokers,
		Topic:       c.Topic,
		Partition:   c.Partition,
		GroupID:     c.Group,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: c.StartOffset,
		Logger:      kafka.LoggerFunc(easy.Logf),
		ErrorLogger: kafka.LoggerFunc(easy.Logf),
	})
	//r.SetOffset(0)

	defer func() {
		if err := r.Close(); err != nil {
			easy.log.Fatal("kafka|read| failed to close reader:", err)
		}
	}()

	errCh := make(chan int)

	interrupt := true

	go func(r *kafka.Reader, ch chan *kafka.Message, errCh chan int) {
		for interrupt {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				easy.log.Errorf("kafka|read|ReadMessage err:%v", err)
				errCh <- 1
				break
			}
			//easy.log.Printf("kafka|read message at offset %d: %s ", m.Offset, string(m.Key))
			ch <- &m
		}

	}(r, ch, errCh)

	select {
	case <-errCh:
		break
	case <-ctx.Done():
		break
	}
	interrupt = false
	time.Sleep(2 * time.Second)
}

func (easy *EasyKafka) WriteBatch(c *Config, ch chan []*kafka.Message, resp chan []*kafka.Message) {
	query := make(chan int, 5)
	for true {
		query <- 1
		go func(c *Config, ch chan []*kafka.Message, resp chan []*kafka.Message) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			easy.Write(*c, ch, resp, ctx)
			<-query
		}(c, ch, resp)
	}
}

func (easy *EasyKafka) Write(c Config, ch chan []*kafka.Message, resp chan []*kafka.Message, ctx context.Context) {
	// make a writer that produces to topic-A, using the least-bytes distribution
	var f kafka.BalancerFunc = func(message kafka.Message, i ...int) int {
		return message.Partition
	}

	w := &kafka.Writer{
		Addr: kafka.TCP(c.Brokers...),
		//Topic:       c.Topic,
		Balancer:               f,
		Logger:                 kafka.LoggerFunc(easy.Logf),
		ErrorLogger:            kafka.LoggerFunc(easy.Logf),
		BatchBytes:             70e6,
		BatchSize:              20,
		AllowAutoTopicCreation: true,
	}

	defer func() {
		if err := w.Close(); err != nil {
			easy.log.Fatal("kafka|write| failed to close writer:", err)
		}
	}()

	running := true

	go func(ctx context.Context) {
		for true {
			select {
			case <-ctx.Done():
				running = false
			}
		}
	}(ctx)

	for running {
		//easy.log.Printf("kafka|Write|length=%v", len(ch))
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

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
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
