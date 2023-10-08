package chain

import (
	"fmt"

	"github.com/0xcregis/easynode/store"
	"github.com/0xcregis/easynode/store/chain/btc"
	"github.com/0xcregis/easynode/store/chain/ether"
	"github.com/0xcregis/easynode/store/chain/filecoin"
	"github.com/0xcregis/easynode/store/chain/polygonpos"
	"github.com/0xcregis/easynode/store/chain/tron"
	"github.com/0xcregis/easynode/store/chain/xrp"
	"github.com/segmentio/kafka-go"
)

func GetReceiptFromKafka(value []byte, blockChain int64) (*store.Receipt, error) {
	if blockChain == 200 {
		return ether.GetReceiptFromKafka(value)
	} else if blockChain == 205 {
		return tron.GetReceiptFromKafka(value)
	} else if blockChain == 201 {
		return polygonpos.GetReceiptFromKafka(value)
	} else if blockChain == 301 {
		return filecoin.GetReceiptFromKafka(value)
	} else if blockChain == 310 {
		return xrp.GetReceiptFromKafka(value)
	} else if blockChain == 300 {
		return btc.GetReceiptFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func GetBlockFromKafka(value []byte, blockChain int64) (*store.Block, error) {
	if blockChain == 200 {
		return ether.GetBlockFromKafka(value)
	} else if blockChain == 205 {
		return tron.GetBlockFromKafka(value)
	} else if blockChain == 201 {
		return polygonpos.GetBlockFromKafka(value)
	} else if blockChain == 301 {
		return filecoin.GetBlockFromKafka(value)
	} else if blockChain == 310 {
		return xrp.GetBlockFromKafka(value)
	} else if blockChain == 300 {
		return btc.GetBlockFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func GetTxFromKafka(value []byte, blockChain int64) (*store.Tx, error) {
	if blockChain == 200 {
		return ether.GetTxFromKafka(value)
	} else if blockChain == 205 {
		return tron.GetTxFromKafka(value)
	} else if blockChain == 201 {
		return polygonpos.GetTxFromKafka(value)
	} else if blockChain == 301 {
		return filecoin.GetTxFromKafka(value)
	} else if blockChain == 310 {
		return xrp.GetTxFromKafka(value)
	} else if blockChain == 300 {
		return btc.GetTxFromKafka(value)
	} else {
		return nil, fmt.Errorf("blockchain:%v does not support", blockChain)
	}
}

func ParseTx(blockchain int64, msg *kafka.Message) (*store.SubTx, error) {
	if blockchain == 200 {
		return ether.ParseTx(msg.Value, store.EthTopic, store.EthTransferSingleTopic, blockchain)
	}
	if blockchain == 205 {
		return tron.ParseTx(msg.Value, store.TronTopic, blockchain)
	}
	if blockchain == 201 {
		return polygonpos.ParseTx(msg.Value, store.PolygonTopic, blockchain)
	}
	if blockchain == 301 {
		return filecoin.ParseTx(msg.Value, store.PolygonTopic, blockchain)
	}
	if blockchain == 310 {
		return xrp.ParseTx(msg.Value, "", blockchain)
	}
	if blockchain == 300 {
		return btc.ParseTx(msg.Value, blockchain)
	}
	return nil, nil
}

func GetTxType(blockchain int64, msg *kafka.Message) (uint64, error) {
	if blockchain == 200 {
		return ether.GetTxType(msg.Value)
	}
	if blockchain == 205 {
		return tron.GetTxType(msg.Value)
	}
	if blockchain == 201 {
		return polygonpos.GetTxType(msg.Value)
	}
	if blockchain == 301 {
		return filecoin.GetTxType(msg.Value)
	}
	if blockchain == 310 {
		return xrp.GetTxType(msg.Value)
	}
	if blockchain == 300 {
		return btc.GetTxType(msg.Value)
	}
	return 0, nil
}

func GetCoreAddress(blockChain int64, address string) string {
	if blockChain == 200 {
		return ether.GetCoreAddr(address)
	} else if blockChain == 205 {
		return tron.GetCoreAddr(address)
	} else if blockChain == 201 {
		return polygonpos.GetCoreAddr(address)
	} else if blockChain == 301 {
		return filecoin.GetCoreAddr(address)
	} else if blockChain == 310 {
		return xrp.GetCoreAddr(address)
	} else if blockChain == 300 {
		return btc.GetCoreAddr(address)
	} else {
		return address
	}
}

func CheckAddress(blockChain int64, msg *kafka.Message, list map[string]*store.MonitorAddress) bool {
	if len(list) < 1 {
		return false
	}
	if blockChain == 200 {
		return ether.CheckAddress(msg.Value, list, store.EthTopic)
	} else if blockChain == 205 {
		return tron.CheckAddress(msg.Value, list, store.TronTopic)
	} else if blockChain == 201 {
		return polygonpos.CheckAddress(msg.Value, list, store.PolygonTopic)
	} else if blockChain == 301 {
		return filecoin.CheckAddress(msg.Value, list, store.PolygonTopic)
	} else if blockChain == 310 {
		return xrp.CheckAddress(msg.Value, list, "")
	} else if blockChain == 300 {
		return btc.CheckAddress(msg.Value, list)
	} else {
		return false
	}
}
