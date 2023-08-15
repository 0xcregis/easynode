package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	kafkaClient "github.com/0xcregis/easynode/common/kafka"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/chain"
	"github.com/0xcregis/easynode/store/config"
	"github.com/0xcregis/easynode/store/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

type WsHandler struct {
	log     *logrus.Entry
	cfg     map[int64]*config.Chain
	kafka   *kafkaClient.EasyKafka
	store   store.DbStoreInterface
	cache   *db.CacheService
	connMap map[string]*websocket.Conn //token:conn
	//cmdMap  map[string]map[int64]store.CmdMessage //token:code-wsReq
	//ctxMap         map[string]map[int64]context.CancelFunc
	//monitorAddress map[int64]map[string]*TokenAddress //blockchain:token-tokenAddress
	writer chan *store.WsPushMessage
}

type TokenAddress struct {
	Token string
	List  map[string]*store.MonitorAddress
}

func NewWsHandler(cfg *config.Config, log *xlog.XLog) *WsHandler {
	kfk := kafkaClient.NewEasyKafka(log)
	ch := db.NewChService(cfg, log)
	mp := make(map[int64]*config.Chain, 2)

	//mp2 := make(map[int64]map[string]*TokenAddress, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
		//mp2[v.BlockChain] = make(map[string]*TokenAddress, 2)
	}
	cache := db.NewCacheService(cfg.Chains, log)
	return &WsHandler{
		log:     log.WithField("model", "wsSrv"),
		store:   ch,
		cfg:     mp,
		kafka:   kfk,
		cache:   cache,
		connMap: make(map[string]*websocket.Conn, 5),
		//cmdMap:  make(map[string]map[int64]store.CmdMessage, 5),
		//ctxMap:         make(map[string]map[int64]context.CancelFunc, 5),
		//monitorAddress: mp2,
		writer: make(chan *store.WsPushMessage),
	}
}

func (ws *WsHandler) Start(kafkaCtx context.Context) {

	//更新监控地址池
	go func() {
		for {
			for _, w := range ws.cfg {
				l, err := ws.store.GetAddressByToken3(w.BlockChain)
				if err != nil {
					continue
				}
				_ = ws.cache.SetMonitorAddress(w.BlockChain, l)
			}

			<-time.After(100 * time.Second)
		}
	}()

	//push message
	go func(ctx context.Context, writer chan *store.WsPushMessage) {
		interrupt := true
		for interrupt {
			select {
			case <-ctx.Done():
				interrupt = false
				break
			case ms := <-writer:
				if _, ok := ws.connMap[ms.Token]; ok {
					bs, _ := json.Marshal(ms)
					err := ws.connMap[ms.Token].WriteMessage(websocket.TextMessage, bs)
					if err != nil {
						ws.log.Errorf("sendMessage|error=%v", err.Error())
					}
				}
			}
			time.Sleep(300 * time.Millisecond)
		}
	}(kafkaCtx, ws.writer)

	// 消费kafka.tx
	txKafkaParams := make(map[int64]*config.KafkaConfig, 2)
	subKafkaParams := make(map[int64]*config.KafkaConfig, 2)
	for b, c := range ws.cfg {
		if v, ok := c.KafkaCfg["Tx"]; ok {
			txKafkaParams[b] = v
		}
		if v, ok := c.KafkaCfg["SubTx"]; ok {
			subKafkaParams[b] = v
		}
	}
	ws.sendMessageEx(txKafkaParams, subKafkaParams, kafkaCtx)
}

func (ws *WsHandler) Sub2(ctx *gin.Context, w http.ResponseWriter, r *http.Request) {

	serialId := ctx.Query("serialId")
	token := ctx.Param("token")
	if len(token) < 1 {
		resp := "{\"code\":1,\"data\":\"not found token\"}"
		_, _ = ctx.Writer.Write([]byte(resp))
		return
	}

	if len(serialId) > 0 {
		token = fmt.Sprintf("%v_%v", serialId, token)
	}

	ws.Sub(w, r, token)
}

func (ws *WsHandler) Sub(w http.ResponseWriter, r *http.Request, token string) {
	upGrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
	c, err := upGrader.Upgrade(w, r, nil)
	if err != nil {
		ws.log.Errorf("upgrade|err=%v", err)
		return
	}
	defer func() {
		_ = c.Close()
	}()

	if old, ok := ws.connMap[token]; ok {
		_ = old.Close()
		delete(ws.connMap, token)
	}
	ws.connMap[token] = c

	c.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(1*time.Second))
		if errors.Is(err, websocket.ErrCloseSent) {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Timeout() {
			return nil
		}
		return err
	})

	cx, cancel := context.WithCancel(context.Background())
	defer cancel()

	PongCh := make(chan string)
	c.SetPongHandler(func(appData string) error {
		PongCh <- appData
		return nil
	})

	//收到pong消息并处理
	go func(PongCh chan string, ws *WsHandler, token string, cancel context.CancelFunc) {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-PongCh: //收到 客户端响应
				ticker.Reset(10 * time.Minute)
				continue
			case <-ticker.C: //已经超时，仍然未收到客户端的响应
				cancel()
				return
			}
		}
	}(PongCh, ws, token, cancel)

	//定时发送 ping 命令
	go func(ws *websocket.Conn, token string, ctx2 context.Context) {
		ticker := time.NewTicker(2 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				if err := ws.WriteControl(websocket.PingMessage, []byte(token), time.Now().Add(10*time.Second)); err != nil {
					log.Println("ping:", err)
				}
			case <-ctx2.Done():
				return
			}
		}
	}(c, token, cx)

	interrupt := true
	//read cmd
	for interrupt {
		select {
		case <-cx.Done():
			//清理
			//delete(ws.cmdMap, token)
			if c, ok := ws.connMap[token]; ok {
				_ = c.Close()
			}
			delete(ws.connMap, token)
			interrupt = false
		}
	}

	//延迟时间，待其他服务清理
	<-time.After(4 * time.Second)
}

func (ws *WsHandler) sendMessageEx(kafkaConfig map[int64]*config.KafkaConfig, subKafkaConfig map[int64]*config.KafkaConfig, ctx context.Context) {
	for b, k := range kafkaConfig {
		s := subKafkaConfig[b]
		go ws.sendMessage(s, k, b, ctx)
	}
}

// kafka->ws.push
func (ws *WsHandler) sendMessage(SubKafkaConfig *config.KafkaConfig, kafkaConfig *config.KafkaConfig, blockChain int64, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	sender := make(chan []*kafka.Message, 10)
	bufferMessage := make([]*kafka.Message, 0, 10)
	monitorAddress := make(map[string]*TokenAddress, 10)
	filters := make(map[string][]*store.SubFilter, 10)
	lock := sync.RWMutex{}

	//read tx from kafka
	go func(ctx context.Context) {
		broker := fmt.Sprintf("%v:%v", kafkaConfig.Host, kafkaConfig.Port)
		group := fmt.Sprintf("gr_push_%v", kafkaConfig.Group)
		c, cancel := context.WithCancel(ctx)
		defer cancel()
		ws.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: kafkaConfig.Topic, Group: group, Partition: kafkaConfig.Partition, StartOffset: kafkaConfig.StartOffset}, receiver, c)
	}(ctx)

	//write sub-tx to kafka
	if SubKafkaConfig != nil {
		go func(ctx context.Context) {
			broker := fmt.Sprintf("%v:%v", SubKafkaConfig.Host, SubKafkaConfig.Port)
			c, cancel := context.WithCancel(ctx)
			defer cancel()
			ws.kafka.Write(kafkaClient.Config{Brokers: []string{broker}}, sender, nil, c)
		}(ctx)
	}

	// 定时更新监控地址池
	go func(blockChain int64, ctx2 context.Context) {
		interrupt := true
		for interrupt {
			select {
			case <-ctx2.Done():
				interrupt = false
				break
			default:
				mp := make(map[string]*TokenAddress, 10)
				for token := range ws.connMap {
					//path中包含serialId请求方式
					newToken := token
					if strings.Contains(token, "_") {
						ls := strings.Split(token, "_")
						if len(ls) == 2 {
							newToken = ls[1]
						}
					}

					ls, _ := ws.store.GetAddressByToken(blockChain, newToken)
					addressMp := make(map[string]*store.MonitorAddress, len(ls))
					for _, a := range ls {
						addressMp[chain.GetCoreAddress(blockChain, a.Address)] = a
					}

					mp[token] = &TokenAddress{
						Token: token,
						List:  addressMp,
					}
				}
				lock.Lock()
				monitorAddress = mp
				lock.Unlock()
			}

			<-time.After(3 * time.Second)
		}
	}(blockChain, ctx)

	//推送的数据备份
	go func(blockChain int64, ctx context.Context) {
		tk := time.NewTicker(5 * time.Second)
		interrupt := true
		for interrupt {
			select {
			case <-ctx.Done():
				interrupt = false
				tk.Stop()
				break
			case <-tk.C:
				if len(bufferMessage) > 0 {
					bf := bufferMessage[:]
					bufferMessage = bufferMessage[len(bufferMessage):]
					sender <- bf
				}
			}

		}
	}(blockChain, ctx)

	//get filters from db for verify tx
	go func(blockChain int64, ctx context.Context) {
		interrupt := true
		for interrupt {
			select {
			case <-ctx.Done():
				interrupt = false
				break
			default:
				list, err := ws.store.GetSubFilter("", blockChain, "")
				if err != nil || len(list) < 1 {
					ws.log.Warnf("SubFilter|blockchain:%v,err:%v", blockChain, "filter is null")
				}

				temp := make(map[string][]*store.SubFilter, 10)
				for _, v := range list {
					if l, ok := temp[v.Token]; ok {
						l = append(l, v)
					} else {
						temp[v.Token] = []*store.SubFilter{v}
					}
				}
				lock.Lock()
				filters = temp
				lock.Unlock()
			}

			<-time.After(60 * time.Second)

		}
	}(blockChain, ctx)

	for {
		msg := <-receiver
		//消息过滤
		//根据这个用户（token）最新的订阅命令，筛选符合条件的交易且推送出去
		if len(ws.connMap) < 1 {
			//用户不在线
			continue
		}

		//用户在线 但没有订阅
		lock.RLock()
		if len(filters) < 1 {
			lock.RUnlock()
			continue
		}
		lock.RUnlock()

		//无用户监控
		lock.RLock()
		if len(monitorAddress) < 1 {
			lock.RUnlock()
			continue
		}
		lock.RUnlock()

		tp, err := chain.GetTxType(blockChain, msg)
		if err != nil {
			continue
		}

		lock.Lock()
		pushMp := ws.TxEngine(filters, monitorAddress, msg, blockChain, tp)
		lock.Unlock()

		//parse tx
		var tx *store.SubTx
		if len(pushMp) > 0 {
			tx, err = chain.ParseTx(blockChain, msg)
			if err != nil {
				ws.log.Warnf("ParseTx|blockchain:%v,kafka.msg:%v,err:%v", blockChain, string(msg.Value), err)
				continue
			}
		}

		//push
		for token, code := range pushMp {
			wpm := &store.WsPushMessage{Code: code, BlockChain: blockChain, Data: tx, Token: token}
			ws.writer <- wpm
		}

		//save to kafka
		if len(pushMp) > 0 && SubKafkaConfig != nil {
			r, _ := json.Marshal(tx)
			m := &kafka.Message{Topic: SubKafkaConfig.Topic, Key: []byte(uuid.New().String()), Value: r}
			bufferMessage = append(bufferMessage, m)
		}

	}

}

// TxEngine the method for verify tx and return
func (ws *WsHandler) TxEngine(filters map[string][]*store.SubFilter, monitorAddress map[string]*TokenAddress, msg *kafka.Message, blockChain int64, txType uint64) map[string]int64 {
	pushMp := make(map[string]int64, 1)
	for token, filter := range filters {
		//push := false
		//var code int64
		code, push := ws.checkCode(filter, txType)

		//不符合订阅条件
		if !push {
			continue
		}

		//检查地址 是否和交易 相关
		//if tokenMp, ok := ws.monitorAddress[blockChain]; ok {
		if tokenAddress, ok := monitorAddress[token]; ok {
			if !chain.CheckAddress(blockChain, msg, tokenAddress.List) {
				continue
			}
		} else {
			continue
		}

		pushMp[token] = code
	}

	return pushMp
}

func (ws *WsHandler) checkCode(mp []*store.SubFilter, tp uint64) (int64, bool) {
	push := false
	var code int64
	for _, q := range mp {
		c, _ := strconv.ParseInt(q.TxCode, 0, 64)
		switch c {
		case 1: //资产交易
			if tp == 1 || tp == 2 {
				push = true
				code = c
			}
		case 3: //质押
			if tp == 6 {
				push = true
				code = c
			}
		case 5: //解质押
			if tp == 7 {
				push = true
				code = c
			}
		case 7: //解提取
			if tp == 8 {
				push = true
				code = c
			}
		case 9: //代理资源
			if tp == 3 {
				push = true
				code = c
			}
		case 11: //代理资源（资源回收）
			if tp == 4 {
				push = true
				code = c
			}
		case 13: //激活账户
			if tp == 5 {
				push = true
				code = c
			}
		}
	}

	return code, push
}
