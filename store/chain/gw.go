package chain

import (
	"fmt"

	chainCode "github.com/0xcregis/easynode/common/chain"
	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/chain/bnb"
	"github.com/0xcregis/easynode/store/chain/btc"
	"github.com/0xcregis/easynode/store/chain/ether"
	"github.com/0xcregis/easynode/store/chain/filecoin"
	"github.com/0xcregis/easynode/store/chain/polygonpos"
	"github.com/0xcregis/easynode/store/chain/tron"
	"github.com/0xcregis/easynode/store/chain/xrp"
	"github.com/segmentio/kafka-go"
)

func GetReceiptFromKafka(value []byte, blockChain int64) (*store.Receipt, error) {
	if chainCode.GetChainCode(blockChain, "ETH", nil) {
		return ether.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "TRON", nil) {
		return tron.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "POLYGON", nil) {
		return polygonpos.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BSC", nil) {
		return bnb.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "FIL", nil) {
		return filecoin.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "XRP", nil) {
		return xrp.GetReceiptFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BTC", nil) {
		return btc.GetReceiptFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func GetBlockFromKafka(value []byte, blockChain int64) (*store.Block, error) {
	if chainCode.GetChainCode(blockChain, "ETH", nil) {
		return ether.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "TRON", nil) {
		return tron.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "POLYGON", nil) {
		return polygonpos.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BSC", nil) {
		return bnb.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "FIL", nil) {
		return filecoin.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "XRP", nil) {
		return xrp.GetBlockFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BTC", nil) {
		return btc.GetBlockFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func GetTxFromKafka(value []byte, blockChain int64) (*store.Tx, error) {
	if chainCode.GetChainCode(blockChain, "ETH", nil) {
		return ether.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "TRON", nil) {
		return tron.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "POLYGON", nil) {
		return polygonpos.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BSC", nil) {
		return bnb.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "FIL", nil) {
		return filecoin.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "XRP", nil) {
		return xrp.GetTxFromKafka(value)
	} else if chainCode.GetChainCode(blockChain, "BTC", nil) {
		return btc.GetTxFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func ParseTx(blockchain int64, msg *kafka.Message) (*store.SubTx, error) {
	if chainCode.GetChainCode(blockchain, "ETH", nil) {
		return ether.ParseTx(msg.Value, store.EthTopic, store.EthTransferSingleTopic, blockchain)
	}
	if chainCode.GetChainCode(blockchain, "TRON", nil) {
		return tron.ParseTx(msg.Value, store.TronTopic, blockchain)
	}
	if chainCode.GetChainCode(blockchain, "POLYGON", nil) {
		return polygonpos.ParseTx(msg.Value, store.PolygonTopic, store.EthTransferSingleTopic, blockchain)
	}
	if chainCode.GetChainCode(blockchain, "BSC", nil) {
		return bnb.ParseTx(msg.Value, store.EthTopic, store.EthTransferSingleTopic, blockchain)
	}
	if chainCode.GetChainCode(blockchain, "FIL", nil) {
		return filecoin.ParseTx(msg.Value, store.PolygonTopic, blockchain)
	}
	if chainCode.GetChainCode(blockchain, "XRP", nil) {
		return xrp.ParseTx(msg.Value, "", blockchain)
	}
	if chainCode.GetChainCode(blockchain, "BTC", nil) {
		return btc.ParseTx(msg.Value, blockchain)
	}
	return nil, nil
}

func GetTxType(blockchain int64, msg *kafka.Message) (uint64, error) {
	if chainCode.GetChainCode(blockchain, "ETH", nil) {
		return ether.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "BSC", nil) {
		return bnb.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "TRON", nil) {
		return tron.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "POLYGON", nil) {
		return polygonpos.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "FIL", nil) {
		return filecoin.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "XRP", nil) {
		return xrp.GetTxType(msg.Value)
	}
	if chainCode.GetChainCode(blockchain, "BTC", nil) {
		return btc.GetTxType(msg.Value)
	}
	return 0, nil
}

func GetCoreAddress(blockChain int64, address string) string {
	if chainCode.GetChainCode(blockChain, "ETH", nil) {
		return ether.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "TRON", nil) {
		//转换地址格式
		address = tron.FormatAddr(address)
		return tron.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "POLYGON", nil) {
		return polygonpos.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "BSC", nil) {
		return bnb.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "FIL", nil) {
		return filecoin.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "XRP", nil) {
		return xrp.GetCoreAddr(address)
	} else if chainCode.GetChainCode(blockChain, "BTC", nil) {
		return btc.GetCoreAddr(address)
	} else {
		return address
	}
}

func CheckAddress(blockChain int64, msg *kafka.Message, list map[string]*store.MonitorAddress) bool {
	if len(list) < 1 {
		return false
	}
	if chainCode.GetChainCode(blockChain, "ETH", nil) {
		return ether.CheckAddress(msg.Value, list, store.EthTopic, store.EthTransferSingleTopic)
	} else if chainCode.GetChainCode(blockChain, "TRON", nil) {
		return tron.CheckAddress(msg.Value, list, store.TronTopic)
	} else if chainCode.GetChainCode(blockChain, "POLYGON", nil) {
		return polygonpos.CheckAddress(msg.Value, list, store.PolygonTopic, store.EthTransferSingleTopic)
	} else if chainCode.GetChainCode(blockChain, "BSC", nil) {
		return bnb.CheckAddress(msg.Value, list, store.EthTopic, store.EthTransferSingleTopic)
	} else if chainCode.GetChainCode(blockChain, "FIL", nil) {
		return filecoin.CheckAddress(msg.Value, list, store.PolygonTopic)
	} else if chainCode.GetChainCode(blockChain, "XRP", nil) {
		return xrp.CheckAddress(msg.Value, list, "")
	} else if chainCode.GetChainCode(blockChain, "BTC", nil) {
		return btc.CheckAddress(msg.Value, list)
	} else {
		return false
	}
}
