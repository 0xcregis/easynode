package chain

import (
	"fmt"
	"time"

	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/ether"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/polygonpos"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/tron2"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

func GetBlockchain(blockchain int, c *config.Chain, store service.StoreTaskInterface, logConfig *config.LogConfig, nodeId string) service.BlockChainInterface {
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/chain_info", logConfig.Path), 24*time.Hour)
	var srv service.BlockChainInterface
	if blockchain == 200 {
		srv = ether.NewService(c, x, store, nodeId, service.EthTopic)
	} else if blockchain == 205 {
		srv = tron2.NewService(c, x, store, nodeId, service.TronTopic)
	} else if blockchain == 201 {
		srv = polygonpos.NewService(c, x, store, nodeId, service.PolygonTopic)
	}

	return srv
}

func GetTxHashFromKafka(blockchain int, msg []byte) string {
	r := gjson.ParseBytes(msg)
	var txHash string
	if blockchain == 200 {
		txHash = r.Get("hash").String()
	} else if blockchain == 205 {
		tx := r.Get("tx").String()
		txHash = gjson.Parse(tx).Get("txID").String()
	} else if blockchain == 201 {
		txHash = r.Get("hash").String()
	}

	return txHash
}

func GetBlockHashFromKafka(blockchain int, msg []byte) string {
	r := gjson.ParseBytes(msg)
	var blockHash string
	if blockchain == 200 {
		blockHash = r.Get("hash").String()
	} else if blockchain == 205 {
		blockHash = r.Get("blockID").String()
	} else if blockchain == 201 {
		blockHash = r.Get("hash").String()
	}
	return blockHash
}

func GetReceiptHashFromKafka(blockchain int, msg []byte) string {
	r := gjson.ParseBytes(msg)
	var txHash string
	if blockchain == 200 {
		txHash = r.Get("transactionHash").String()
	} else if blockchain == 205 {
		txHash = r.Get("id").String()
	} else if blockchain == 201 {
		txHash = r.Get("transactionHash").String()
	}
	return txHash
}
