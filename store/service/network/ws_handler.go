package network

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	db2 "github.com/uduncloud/easynode/store/service/db"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type WsHandler struct {
	log     *logrus.Entry
	cfg     map[int64]*config.Chain
	kafka   *kafkaClient.EasyKafka
	db      service.DbMonitorAddressInterface
	cache   *db2.CacheService
	connMap map[string]*websocket.Conn                //token:conn
	cmdMap  map[string]map[int64]service.WsReqMessage //token:code-wsReq
	//ctxMap         map[string]map[int64]context.CancelFunc
	monitorAddress map[int64]map[string]*TokenAddress //blockchain:token-tokenAddress
}

type TokenAddress struct {
	Token string
	List  []*service.MonitorAddress
}

func NewWsHandler(cfg *config.Config, log *xlog.XLog) *WsHandler {
	kfk := kafkaClient.NewEasyKafka(log)
	db := db2.NewChService(cfg, log)
	mp := make(map[int64]*config.Chain, 2)
	mp2 := make(map[int64]map[string]*TokenAddress, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
		mp2[v.BlockChain] = make(map[string]*TokenAddress, 2)
	}
	cache := db2.NewCacheService(cfg.Chains, log)
	return &WsHandler{
		log:     log.WithField("model", "wsSrv"),
		db:      db,
		cfg:     mp,
		kafka:   kfk,
		cache:   cache,
		connMap: make(map[string]*websocket.Conn, 5),
		cmdMap:  make(map[string]map[int64]service.WsReqMessage, 5),
		//ctxMap:         make(map[string]map[int64]context.CancelFunc, 5),
		monitorAddress: mp2,
	}
}

func (ws *WsHandler) Start(kafkaCtx context.Context) {

	//更新监控地址池
	go func() {
		for {
			for _, w := range ws.cfg {
				l, err := ws.db.GetAddressByToken3(w.BlockChain)
				if err != nil {
					continue
				}
				_ = ws.cache.SetMonitorAddress(w.BlockChain, l)
			}

			<-time.After(30 * time.Second)
		}
	}()

	//go func() {
	// 订阅链的配置
	txKafkaParams := make(map[int64]*config.KafkaConfig, 2)
	subKafkaParams := make(map[int64]*config.KafkaConfig, 2)
	//c1, _ := context.WithCancel(context.Background())
	//defer cancel()
	for b, c := range ws.cfg {
		f := c.KafkaCfg["Tx"]
		txKafkaParams[b] = f

		s := c.KafkaCfg["SubTx"]
		subKafkaParams[b] = s
	}
	ws.sendMessageEx(txKafkaParams, subKafkaParams, kafkaCtx)
	//}()
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

	ws.Sub(ctx, w, r, token)
}

func (ws *WsHandler) Sub(ctx *gin.Context, w http.ResponseWriter, r *http.Request, token string) {
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
		if err == websocket.ErrCloseSent {
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
			delete(ws.cmdMap, token)
			if c, ok := ws.connMap[token]; ok {
				_ = c.Close()
			}
			delete(ws.connMap, token)
			//if fs, ok := ws.ctxMap[token]; ok {
			//	for _, f := range fs {
			//		f()
			//	}
			//}
			//delete(ws.ctxMap, token)
			interrupt = false
		default:
			logger := ws.log.WithFields(logrus.Fields{
				"id":    time.Now().UnixMilli(),
				"model": "ws.start",
			})

			mt, message, err := c.ReadMessage()
			if err != nil {
				fmt.Printf("ReadMessage|err=%v", err)
				cancel()
				break
			}
			//log.Printf("ReadMessage: %s", message)
			ws.handlerMessage(cx, token, c, mt, message, logger)
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
	//addressList := make([]*service.MonitorAddress, 0, 10)
	//控制消息处理速度
	//lock := make(chan int64)

	go func(ctx context.Context) {
		broker := fmt.Sprintf("%v:%v", kafkaConfig.Host, kafkaConfig.Port)
		group := fmt.Sprintf("gr_push_%v", kafkaConfig.Group)
		c, cancel := context.WithCancel(ctx)
		defer cancel()
		ws.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: kafkaConfig.Topic, Group: group, Partition: kafkaConfig.Partition, StartOffset: kafkaConfig.StartOffset}, receiver, c)
	}(ctx)

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

					l, _ := ws.db.GetAddressByToken(blockChain, newToken)

					mp[token] = &TokenAddress{
						Token: token,
						List:  l,
					}
				}
				ws.monitorAddress[blockChain] = mp
			}

			<-time.After(3 * time.Second)
		}
	}(blockChain, ctx)

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

	for {
		select {
		case msg := <-receiver:
			//消息过滤
			//根据这个用户（token）最新的订阅命令，筛选符合条件的交易且推送出去
			if len(ws.connMap) < 1 {
				//用户不在线
				continue
			}

			//用户在线 但没有订阅
			if len(ws.cmdMap) < 1 {
				continue
			}

			//无用户监控
			if l, ok := ws.monitorAddress[blockChain]; !ok || len(l) < 1 {
				continue
			}

			tp, err := service.GetTxType(blockChain, msg)
			if err != nil {
				continue
			}

			pushMp := make(map[string]int64, 1)
			for token, mp := range ws.cmdMap {
				//push := false
				//var code int64
				code, push := ws.CheckCode(mp, tp, blockChain)

				//不符合订阅条件
				if !push {
					continue
				}

				//检查地址 是否和交易 相关
				if v, ok := ws.monitorAddress[blockChain]; ok {
					if tokenAddress, ok := v[token]; ok {
						if !ws.checkTx(blockChain, msg, tokenAddress.List) {
							continue
						}
					} else {
						continue
					}
				} else {
					continue
				}

				pushMp[token] = code
			}

			//parse tx
			var tx *service.SubTx
			if len(pushMp) > 0 {
				tx, err = service.ParseTx(blockChain, msg)
				if err != nil {
					ws.log.Warnf("ParseTx|blockchain:%v,kafka.msg:%v,err:%v", blockChain, string(msg.Value), err)
					continue
				}
			}

			//push
			for token, code := range pushMp {
				wpm := service.WsPushMessage{Code: code, BlockChain: blockChain, Data: tx, Token: token}
				bs, _ := json.Marshal(wpm)
				err = ws.connMap[token].WriteMessage(websocket.TextMessage, bs)
				if err != nil {
					ws.log.Errorf("sendMessage|error=%v", err.Error())
				}
			}

			//save to kafka
			if len(pushMp) > 0 && SubKafkaConfig != nil {
				r, _ := json.Marshal(tx)
				m := &kafka.Message{Topic: SubKafkaConfig.Topic, Key: []byte(uuid.New().String()), Value: r}
				bufferMessage = append(bufferMessage, m)

				//if len(bufferMessage) > 10 {
				//	bf := bufferMessage[:]
				//	bufferMessage = bufferMessage[len(bufferMessage):]
				//	sender <- bf
				//}
			}

			//<-lock
		}

	}
}

func (ws *WsHandler) CheckCode(mp map[int64]service.WsReqMessage, tp uint64, blockChain int64) (int64, bool) {
	push := false
	var code int64
	for c, q := range mp {

		has := false
		for _, v := range q.BlockChain {
			if v == blockChain {
				has = true
				break
			}
		}

		if !has {
			continue
		}

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

func (ws *WsHandler) checkTx(blockChain int64, msg *kafka.Message, list []*service.MonitorAddress) bool {
	return service.CheckAddress(blockChain, msg, list)
}

func (ws *WsHandler) handlerMessage(ctx context.Context, token string, c *websocket.Conn, mt int, message []byte, log *logrus.Entry) {

	//根据命令不同执行不同函数
	var msg service.WsReqMessage
	var returnMsg service.WsRespMessage

	err := json.Unmarshal(message, &msg)
	if err != nil {
		errMsg := &service.WsRespMessage{Id: msg.Id, Code: msg.Code, Err: err.Error(), Params: msg.Params, Status: 1, BlockChain: msg.BlockChain, Resp: nil}
		ws.returnMsg(errMsg, log, mt, c)
		return
	}

	//初始化返回
	returnMsg.Id = msg.Id
	returnMsg.Code = msg.Code
	returnMsg.Params = msg.Params
	returnMsg.BlockChain = msg.BlockChain
	returnMsg.Status = 0

	//最终返回
	defer func(r *service.WsRespMessage) {
		ws.returnMsg(r, log, mt, c)
	}(&returnMsg)

	//不支持
	for _, blockchain := range msg.BlockChain {
		if _, ok := ws.cfg[blockchain]; !ok {
			returnMsg.Status = 1
			returnMsg.Err = fmt.Sprintf("the blockchain(%v) has not support", msg.BlockChain)
			return
		}
	}

	if msg.Code == 1 || msg.Code == 3 || msg.Code == 5 || msg.Code == 7 || msg.Code == 9 || msg.Code == 11 || msg.Code == 13 {
		//订阅请求
		v, ok := ws.cmdMap[token]
		if ok {
			_, ok = v[msg.Code]
		}
		if ok {
			//已经订阅了，则返回订阅失败
			returnMsg.Err = fmt.Sprintf("sub cmd is already existed")
			returnMsg.Status = 1
		} else {
			//未订阅时，则订阅成功
			if _, ok := ws.cmdMap[token]; !ok {
				ws.cmdMap[token] = map[int64]service.WsReqMessage{msg.Code: msg}
			} else {
				ws.cmdMap[token][msg.Code] = msg
			}

			//kafkaCtx, cancel := context.WithCancel(ctx)
			//if _, ok := ws.ctxMap[token]; !ok {
			//	ws.ctxMap[token] = map[int64]context.CancelFunc{msg.Code: cancel}
			//} else {
			//	ws.ctxMap[token][msg.Code] = cancel
			//}
			//ws.sendMessageEx(msg.Code, token, txKafkaParams, subKafkaParams, kafkaCtx)
		}
	} else if msg.Code == 2 || msg.Code == 4 || msg.Code == 6 || msg.Code == 8 || msg.Code == 10 || msg.Code == 12 || msg.Code == 14 {
		//取消订阅，回收资源
		//if fs, ok := ws.ctxMap[token]; ok {
		//	if f, ok2 := fs[msg.Code-1]; ok2 {
		//		f()
		//	}
		//}
		//
		//if v, ok := ws.ctxMap[token]; ok {
		//	delete(v, msg.Code-1)
		//}

		if v, ok := ws.cmdMap[token]; ok {
			delete(v, msg.Code-1)
		}
	}

}

func (ws *WsHandler) returnMsg(r *service.WsRespMessage, log *logrus.Entry, mt int, c *websocket.Conn) {
	bs, err := json.Marshal(r)
	if err != nil {
		log.Errorf("response|err=%v", err)
	}
	err = c.WriteMessage(mt, bs)
	if err != nil {
		log.Errorf("response|err=%v", err)
	}
}
