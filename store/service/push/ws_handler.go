package push

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	kafkaClient "github.com/uduncloud/easynode/common/kafka"
	"github.com/uduncloud/easynode/store/config"
	"github.com/uduncloud/easynode/store/service"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

type WsHandler struct {
	log     *xlog.XLog
	cfg     map[int64]*config.Chain
	kafka   *kafkaClient.EasyKafka
	db      service.DbMonitorAddressInterface
	connMap map[string]*websocket.Conn
	cmdMap  map[string]service.WsReqMessage
	ctxMap  map[string]context.CancelFunc
}

func NewWsHandler(cfg *config.Config, xlog *xlog.XLog) *WsHandler {
	kfk := kafkaClient.NewEasyKafka(xlog)
	db := NewChService(cfg, xlog)
	mp := make(map[int64]*config.Chain, 2)
	for _, v := range cfg.Chains {
		mp[v.BlockChain] = v
	}
	return &WsHandler{
		log:     xlog,
		db:      db,
		cfg:     mp,
		kafka:   kfk,
		connMap: make(map[string]*websocket.Conn, 5),
		cmdMap:  make(map[string]service.WsReqMessage, 5),
		ctxMap:  make(map[string]context.CancelFunc, 5),
	}
}

func (ws *WsHandler) Start(ctx *gin.Context, w http.ResponseWriter, r *http.Request) {

	token := ctx.Param("token")
	if len(token) < 1 {
		resp := "{\"code\":1,\"data\":\"not found token\"}"
		_, _ = ctx.Writer.Write([]byte(resp))
		return
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// 解决跨域问题
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ws.log.Errorf("upgrade|err=%v", err)
		return
	}
	defer c.Close()

	if old, ok := ws.connMap[token]; ok {
		old.Close()
		delete(ws.connMap, token)
	}
	ws.connMap[token] = c

	c.SetPingHandler(func(message string) error {
		err := c.WriteControl(websocket.PongMessage, []byte(message), time.Now().Add(1*time.Second))
		if err == websocket.ErrCloseSent {
			return nil
		} else if e, ok := err.(net.Error); ok && e.Temporary() {
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

	//监听
	go func(PongCh chan string, ws *WsHandler, token string, cancel context.CancelFunc) {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for {
			select {
			case <-PongCh: //收到 客户端响应
				ticker.Reset(10 * time.Minute)
				continue
			case <-ticker.C: //已经超时，仍然未收到客户端的响应
				//清理
				delete(ws.cmdMap, token)
				if c, ok := ws.connMap[token]; ok {
					_ = c.Close()
				}
				delete(ws.connMap, token)
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

	//tx->push
	//go ws.sendMessage(token, ws.cfg.KafkaCfg["Tx"], ws.cfg.BlockChain, cx)

	//read cmd
	for {

		log := ws.log.WithFields(logrus.Fields{
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
		ws.handlerMessage(cx, token, c, mt, message, log)
	}

	//延迟时间，待其他服务清理
	<-time.After(3 * time.Second)

}

func (ws *WsHandler) sendMessageEx(token string, kafkaConfig map[int64]*config.KafkaConfig, ctx context.Context) {
	for b, k := range kafkaConfig {
		go ws.sendMessage(token, k, b, ctx)
	}
}

// kafka->ws.push
func (ws *WsHandler) sendMessage(token string, kafkaConfig *config.KafkaConfig, blockChain int64, ctx context.Context) {
	receiver := make(chan *kafka.Message)
	go func(ctx context.Context) {
		broker := fmt.Sprintf("%v:%v", kafkaConfig.Host, kafkaConfig.Port)
		group := fmt.Sprintf("group_push_%v", kafkaConfig.Group)
		ws.kafka.Read(&kafkaClient.Config{Brokers: []string{broker}, Topic: kafkaConfig.Topic, Group: group, Partition: kafkaConfig.Partition, StartOffset: kafkaConfig.StartOffset}, receiver, ctx)
	}(ctx)

	for {
		select {
		case <-ctx.Done():
			if _, ok := ws.ctxMap[token]; ok {
				delete(ws.ctxMap, token)
			}
			if _, ok := ws.cmdMap[token]; ok {
				delete(ws.cmdMap, token)
			}
			return
		case msg := <-receiver:
			//消息过滤
			//根据这个用户（token）最新的订阅命令，筛选符合条件的交易且推送出去
			//todo filter
			if _, ok := ws.cmdMap[token]; !ok {
				//用户不在线
				continue
			} else {
				//用户在线
				if _, ok := ws.cmdMap[token]; !ok { //没有订阅
					continue
				}

				// 其他各类交易
				//该用户订阅的地址 是否和该交易相匹配
				//不匹配 则返回
				list, err := ws.db.GetAddressByToken(blockChain, token)
				if err != nil || len(list) < 1 {
					continue
				}

				//检查地址 是否和交易 相关
				{
					has := false
					if blockChain == 200 {
						has = ws.CheckAddressForEther(msg, list)
					}

					if blockChain == 205 {
						has = ws.CheckAddressForTron(msg, list)
					}

					if !has {
						continue
					}
				}

			}

			data := service.ParseTx(blockChain, msg)
			wpm := service.WsPushMessage{Code: 1, BlockChain: blockChain, Data: data}
			bs, _ := json.Marshal(wpm)
			err := ws.connMap[token].WriteMessage(websocket.TextMessage, bs)
			if err != nil {
				ws.log.Errorf("sendMessage|error=%v", err.Error())
			}

		}

		time.Sleep(1 * time.Second)
	}
}

func (ws *WsHandler) CheckAddressForEther(msg *kafka.Message, list []*service.MonitorAddress) bool {
	root := gjson.ParseBytes(msg.Value)
	fromAddr := root.Get("from").String()
	toAddr := root.Get("to").String()
	has := false
	for _, v := range list {
		//已经判断出该交易 符合要求了，不需要在检查其他地址了
		if has {
			break
		}

		// 该交易是否是订阅的类型交易
		//if v.TxType== tx.Type {
		// has=true
		//break
		//}

		// 普通交易且 地址包含订阅地址
		if strings.HasPrefix(fromAddr, v.Address) || strings.HasPrefix(toAddr, v.Address) {
			has = true
			break
		}

		//合约交易
		monitorAddr := strings.TrimLeft(v.Address, "0x") //去丢0x
		if root.Get("receipt").Exists() {
			receipt := root.Get("receipt").String()
			receiptRoot := gjson.Parse(receipt)
			list := receiptRoot.Get("logs").Array()
			for _, v := range list {
				topics := v.Get("topics").Array()
				//Transfer()
				if len(topics) >= 3 && topics[0].String() == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
					if strings.HasSuffix(topics[1].String(), monitorAddr) || strings.HasSuffix(topics[2].String(), monitorAddr) {
						has = true
						break
					}

				}

			}

		}
	}
	return has
}

func (ws *WsHandler) CheckAddressForTron(msg *kafka.Message, list []*service.MonitorAddress) bool {
	root := gjson.ParseBytes(msg.Value)
	tx := root.Get("tx").String()
	txRoot := gjson.Parse(tx)
	contracts := txRoot.Get("raw_data.contract").Array()
	if len(contracts) < 1 {
		return false
	}
	r := contracts[0]
	txType := r.Get("type").String()

	var fromAddr, toAddr string
	var logs []gjson.Result
	var internalTransactions []gjson.Result
	if txType == "TransferContract" {
		fromAddr = r.Get("parameter.value.owner_address").String()
		toAddr = r.Get("parameter.value.to_address").String()
		//r.Get("parameter.value.amount").String()
	} else if txType == "TriggerSmartContract" {

		receipt := root.Get("receipt").String()
		receiptRoot := gjson.Parse(receipt)
		if receiptRoot.Get("receipt.result").String() != "SUCCESS" {
			return false
		}
		logs = receiptRoot.Get("log").Array()
		internalTransactions = receiptRoot.Get("internal_transactions").Array()
	}

	has := false
	for _, v := range list {
		//已经判断出该交易 符合要求了，不需要在检查其他地址了
		if has {
			break
		}

		// 该交易是否是订阅的类型交易
		//if v.TxType== tx.Type {
		// has=true
		//break
		//}

		// 普通交易且 地址包含订阅地址
		if txType == "TransferContract" {
			if strings.HasSuffix(fromAddr, v.Address) || strings.HasSuffix(toAddr, v.Address) {
				has = true
				break
			}
		} else if txType == "TriggerSmartContract" {
			//合约交易
			var monitorAddr string
			if strings.HasPrefix(v.Address, "0x") {
				monitorAddr = strings.TrimLeft(v.Address, "0x") //去丢0x
			}
			if strings.HasPrefix(v.Address, "41") {
				monitorAddr = strings.TrimLeft(v.Address, "41") //去丢0x
			}

			//合约调用下的TRC20
			if len(logs) > 0 {
				for _, v := range logs {
					topics := v.Get("topics").Array()
					//Transfer()
					if len(topics) >= 3 && topics[0].String() == "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef" {
						if strings.HasSuffix(topics[1].String(), monitorAddr) || strings.HasSuffix(topics[2].String(), monitorAddr) {
							has = true
							break
						}
					}
				}
			}

			//合约调用下的内部交易TRX转帐和TRC10转账：
			if len(internalTransactions) > 0 {
				for _, v := range internalTransactions {
					fromAddr = v.Get("caller_address").String()
					toAddr = v.Get("transferTo_address").String()
					if strings.HasSuffix(fromAddr, monitorAddr) || strings.HasSuffix(toAddr, monitorAddr) {
						has = true
						break
					}
				}
			}
		}
	}
	return has
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

	// 订阅链的配置
	txKafkaParams := make(map[int64]*config.KafkaConfig, 2)
	for _, b := range msg.BlockChain {
		if c, ok := ws.cfg[b]; ok {
			f := c.KafkaCfg["Tx"]
			txKafkaParams[b] = f
		}
	}

	//最新一次命令
	if msg.Code == 1 {
		//仅能保存一个订阅的命令
		if _, ok := ws.cmdMap[token]; ok {
			//已经订阅了，则返回订阅失败
			returnMsg.Err = fmt.Sprintf("sub cmd is already existed")
			returnMsg.Status = 1
		} else {
			//未订阅时，则订阅成功
			ws.cmdMap[token] = msg
			kafkaCtx, cancel := context.WithCancel(ctx)
			ws.ctxMap[token] = cancel
			ws.sendMessageEx(token, txKafkaParams, kafkaCtx)
		}
	} else if msg.Code == 2 {
		if f, ok := ws.ctxMap[token]; ok {
			f()
		}
		delete(ws.ctxMap, token)
		delete(ws.cmdMap, token)
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
