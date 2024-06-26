package chain

import (
	"fmt"
	"time"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/bnb"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/btc"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/ether"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/filecoin"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/polygonpos"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/tron2"
	"github.com/0xcregis/easynode/collect/service/cmd/chain/xrp"
	chainCode "github.com/0xcregis/easynode/common/chain"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

func GetBlockchain(blockchain int, c *config.Chain, store collect.StoreTaskInterface, logConfig *config.LogConfig, nodeId string) collect.BlockChainInterface {
	if logConfig.LogLevel == 0 {
		logConfig.LogLevel = 4
	}
	code := int64(blockchain)
	x := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFormatter(xlog.FORMAT_JSON).BuildLevel(xlog.Level(logConfig.LogLevel)).BuildFile(fmt.Sprintf("%v/chain_info", logConfig.Path), 24*time.Hour)
	var srv collect.BlockChainInterface
	if chainCode.GetChainCode(code, "ETH", x) {
		srv = ether.NewService(c, x, store, nodeId, collect.EthTopic, collect.EthNftTransferSingleTopic)
	} else if chainCode.GetChainCode(code, "TRON", x) {
		srv = tron2.NewService(c, x, store, nodeId, collect.TronTopic)
	} else if chainCode.GetChainCode(code, "POLYGON", x) {
		srv = polygonpos.NewService(c, x, store, nodeId, collect.PolygonTopic, collect.EthNftTransferSingleTopic)
	} else if chainCode.GetChainCode(code, "BSC", x) {
		srv = bnb.NewService(c, x, store, nodeId, collect.EthTopic, collect.EthNftTransferSingleTopic)
	} else if chainCode.GetChainCode(code, "FIL", x) {
		srv = filecoin.NewService(c, x, store, nodeId, "")
	} else if chainCode.GetChainCode(code, "XRP", x) {
		srv = xrp.NewService(c, x, store, nodeId, "")
	} else if chainCode.GetChainCode(code, "BTC", x) {
		srv = btc.NewService(c, x, store, nodeId)
	}

	return srv
}

func GetTxHashFromKafka(blockchain int, txMsg []byte) string {
	code := int64(blockchain)
	r := gjson.ParseBytes(txMsg)
	var txHash string
	if chainCode.GetChainCode(code, "ETH", nil) {
		txHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "TRON", nil) {
		tx := r.Get("tx").String()
		txHash = gjson.Parse(tx).Get("txID").String()
	} else if chainCode.GetChainCode(code, "POLYGON", nil) {
		txHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "BSC", nil) {
		txHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "FIL", nil) {
		txHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "XRP", nil) {
		txHash = r.Get("hash").String()
	}

	return txHash
}

func GetBlockHashFromKafka(blockchain int, blockMsg []byte) string {
	code := int64(blockchain)
	r := gjson.ParseBytes(blockMsg)
	var blockHash string
	if chainCode.GetChainCode(code, "ETH", nil) {
		blockHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "TRON", nil) {
		blockHash = r.Get("blockID").String()
	} else if chainCode.GetChainCode(code, "POLYGON", nil) {
		blockHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "BSC", nil) {
		blockHash = r.Get("hash").String()
	} else if chainCode.GetChainCode(code, "FIL", nil) {
		blockHash = r.Get("blockHash").String()
	} else if chainCode.GetChainCode(code, "XRP", nil) {
		blockHash = r.Get("ledger_hash").String()
	}
	return blockHash
}

func GetReceiptHashFromKafka(blockchain int, receiptMsg []byte) string {
	code := int64(blockchain)
	r := gjson.ParseBytes(receiptMsg)
	var txHash string
	if chainCode.GetChainCode(code, "ETH", nil) {
		txHash = r.Get("transactionHash").String()
	} else if chainCode.GetChainCode(code, "TRON", nil) {
		txHash = r.Get("id").String()
	} else if chainCode.GetChainCode(code, "POLYGON", nil) {
		txHash = r.Get("transactionHash").String()
	} else if chainCode.GetChainCode(code, "BSC", nil) {
		txHash = r.Get("transactionHash").String()
	} else if chainCode.GetChainCode(code, "FIL", nil) {
		txHash = r.Get("transactionHash").String()
	} else if chainCode.GetChainCode(code, "XRP", nil) {
		txHash = r.Get("hash").String()
	}
	return txHash
}
