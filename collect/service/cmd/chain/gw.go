package chain

import (
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"github.com/uduncloud/easynode/collect/service/cmd/chain/ether"
	"github.com/uduncloud/easynode/collect/service/cmd/chain/tron2"
	"time"
)

func GetBlockchain(blockchain int, c *config.Chain, store service.StoreTaskInterface, logConfig *config.LogConfig) service.BlockChainInterface {
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildFile(fmt.Sprintf("%v/chain_info", logConfig.Path), 24*time.Hour)
	var srv service.BlockChainInterface
	if blockchain == 200 {
		srv = ether.NewService(c, x, store)
	} else if blockchain == 205 {
		srv = tron2.NewService(c, x, store)
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
	}
	return txHash
}
