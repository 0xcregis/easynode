package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/blockchain/config"
	"net/http"
	"strconv"
	"time"
)

type WsHandler struct {
	log               *logrus.Entry
	nodeCluster       map[int64][]*config.NodeCluster
	blockChainClients map[int64]API
}

func NewWsHandler(cluster map[int64][]*config.NodeCluster, xlog *xlog.XLog) *WsHandler {
	return &WsHandler{
		log:               xlog.WithField("model", "wsSrv"),
		nodeCluster:       cluster,
		blockChainClients: NewApis(cluster, xlog),
	}
}

func (s *WsHandler) Start(w http.ResponseWriter, r *http.Request) {
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
		s.log.Errorf("upgrade|err=%v", err)
		return
	}
	defer c.Close()

	for {
		log := s.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "ws.start",
		})

		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Errorf("ReadMessage|err=%v", err)
			continue
		}
		//log.Printf("ReadMessage: %s", message)

		go s.handlerMessage(c, mt, message, log)
	}
}

func (s *WsHandler) handlerMessage(c *websocket.Conn, mt int, message []byte, log *logrus.Entry) {

	//根据命令不同执行不同函数
	var msg WsReqMessage
	var returnMsg WsRespMessage

	err := json.Unmarshal(message, &msg)
	if err != nil {
		errMsg := &WsRespMessage{Id: msg.Id, Code: msg.Code, Err: err.Error(), Params: msg.Params, Status: 1, Resp: nil}
		s.returnMsg(errMsg, log, mt, c)
		return
	}
	//初始化返回
	returnMsg.Id = msg.Id
	returnMsg.Code = msg.Code
	returnMsg.Params = msg.Params
	returnMsg.Status = 0

	//最终返回
	defer func(r *WsRespMessage) {
		s.returnMsg(r, log, mt, c)
	}(&returnMsg)

	returnCh := make(chan interface{})
	errCh := make(chan error)

	log.Printf("request.WsReqMessage=%+v", msg)
	blockChain, err := strconv.ParseInt(msg.Params["blockChain"], 0, 64)
	if err != nil {
		returnMsg.Status = 1
		returnMsg.Err = err.Error()
		returnMsg.Resp = nil
		return
	}
	go s.requestMsg(&msg, log, errCh, returnCh, blockChain)

	select {
	case err = <-errCh:
		{
			returnMsg.Status = 1
			returnMsg.Err = err.Error()
			returnMsg.Resp = nil
		}
	case r := <-returnCh:
		{
			returnMsg.Status = 0
			returnMsg.Resp = r
		}
	}
}

func (s *WsHandler) requestMsg(msg *WsReqMessage, log *logrus.Entry, errCh chan error, returnCh chan interface{}, blockChain int64) {
	switch msg.Code {
	case 1: //单笔交易
		{
			//set txHash value  from message
			var txHash string
			if v, ok := msg.Params["txHash"]; ok {
				txHash = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found txHash")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			//是否同步参数
			sync := "0"
			if v, ok := msg.Params["sync"]; ok {
				sync = v
			} else {
				sync = "0"
			}

			//异步请求 sync="0"
			if sync == "0" {
				t, err := s.blockChainClients[blockChain].GetTxByHash(blockChain, txHash)
				if err != nil {
					errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetTxByHash", msg.Code, err.Error())
					log.Errorf(errMsg)
					errCh <- errors.New(errMsg)
					return
				}
				returnCh <- t
				return
			}

			//同步请求 sync=1
			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()
			interrupt := true
			tempCh := make(chan interface{})
			go func(blockChain int64, txHash string, tempCh chan interface{}, cancel context.CancelFunc) {
				for sync == "1" && interrupt {
					t, err := s.blockChainClients[blockChain].GetTxByHash(blockChain, txHash)
					if err != nil {
						cancel()
						break
					}

					r := gjson.Parse(t)
					log.Println("result=", r.Get("result").String())
					if !r.Get("result").Exists() || (r.Get("result").String() == "") {
						time.Sleep(10 * time.Second)
						continue
					}

					tempCh <- t
				}
			}(blockChain, txHash, tempCh, cancel)

			select {
			case <-ctx.Done():
				interrupt = false
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetTxByHash", msg.Code, "timeout")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
			case r := <-tempCh:
				returnCh <- r
			}

		}

	case 2: //单笔区块
		{
			//set blockHash value  from message
			var blockHash string
			if v, ok := msg.Params["blockHash"]; ok {
				blockHash = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found blockHash")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			block, err := s.blockChainClients[blockChain].GetBlockByHash(blockChain, blockHash, true)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetBlockByHash", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- block
		}
	case 3: //单笔区块
		{
			//set blockHash value  from message
			var number string
			if v, ok := msg.Params["number"]; ok {
				number = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found number")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			block, err := s.blockChainClients[blockChain].GetBlockByNumber(blockChain, number, true)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetBlockByNumber", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- block
		}
	case 4: //最新区块
		{
			block, err := s.blockChainClients[blockChain].LatestBlock(blockChain)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "LatestBlock", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- block
		}
	case 8: //sendRawTx
		{
			var signed string
			if v, ok := msg.Params["signed"]; ok {
				signed = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found signed")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			r, err := s.blockChainClients[blockChain].SendRawTransaction(blockChain, signed)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "SendRawTransaction", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- r
		}
	case 9: //call
		{
			var jsonRpc string
			if v, ok := msg.Params["jsonRpc"]; ok {
				jsonRpc = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found jsonRpc")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			r, err := s.blockChainClients[blockChain].SendJsonRpc(blockChain, jsonRpc)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "SendJsonRpc", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- r
		}
	case 10: //主币余额查询
		{
			//set txHash value  from message
			var addr string
			if v, ok := msg.Params["address"]; ok {
				addr = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found addr")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			balance, err := s.blockChainClients[blockChain].Balance(blockChain, addr, "latest")
			//err := s.chain.SubBalance(addr, balanceCh)
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetMainBalance", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- balance
		}

	case 11: //代币余额查询
		{
			//set txHash value  from message
			var addr string
			if v, ok := msg.Params["address"]; ok {
				addr = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found addr")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			var contract string
			if v, ok := msg.Params["contract"]; ok {
				contract = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found contract")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			tokenBalance, err := s.blockChainClients[blockChain].TokenBalance(blockChain, addr, contract, "")
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "GetTokenBalance", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- tokenBalance
		}
	case 12: //nonce
		{
			//set txHash value  from message
			var addr string
			if v, ok := msg.Params["address"]; ok {
				addr = v
			} else {
				errMsg := fmt.Sprintf("msg.code=%v,err=%v", msg.Code, "not found addr")
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}

			nonce, err := s.blockChainClients[blockChain].Nonce(blockChain, addr, "latest")
			if err != nil {
				errMsg := fmt.Sprintf("function=%v,msg.code=%v,err=%v", "Nonce", msg.Code, err.Error())
				log.Errorf(errMsg)
				errCh <- errors.New(errMsg)
				return
			}
			returnCh <- nonce
		}

	}
}

func (s *WsHandler) returnMsg(r *WsRespMessage, log *logrus.Entry, mt int, c *websocket.Conn) {
	bs, err := json.Marshal(r)
	if err != nil {
		log.Errorf("response|err=%v", err)
	}
	err = c.WriteMessage(mt, bs)
	if err != nil {
		log.Errorf("response|err=%v", err)
	}
}
