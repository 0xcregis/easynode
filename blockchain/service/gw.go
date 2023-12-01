package service

import (
	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/0xcregis/easynode/blockchain/service/bnb"
	"github.com/0xcregis/easynode/blockchain/service/btc"
	"github.com/0xcregis/easynode/blockchain/service/ether"
	"github.com/0xcregis/easynode/blockchain/service/filecoin"
	"github.com/0xcregis/easynode/blockchain/service/polygon"
	"github.com/0xcregis/easynode/blockchain/service/tron"
	"github.com/0xcregis/easynode/blockchain/service/xrp"
	"github.com/0xcregis/easynode/common/chain"
	"github.com/sunjiangjun/xlog"
)

func NewApi(blockchain int64, cluster []*config.NodeCluster, xlog *xlog.XLog) blockchain.API {
	if blockchain == chain.GetChainCode("ETH", xlog) {
		return ether.NewEth(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("TRON", xlog) {
		return tron.NewTron(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("POLYGON", xlog) {
		return polygon.NewPolygonPos(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("BSC", xlog) {
		return bnb.NewBnb(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("BTC", xlog) {
		return btc.NewBtc(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("FIL", xlog) {
		return filecoin.NewFileCoin(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("XRP", xlog) {
		return xrp.NewXRP(cluster, blockchain, xlog)
	}
	return nil
}

func NewNftApi(blockchain int64, cluster []*config.NodeCluster, xlog *xlog.XLog) blockchain.NftApi {
	if blockchain == chain.GetChainCode("ETH", xlog) {
		return ether.NewNftEth(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("POLYGON", xlog) {
		return polygon.NewNftPolygonPos(cluster, blockchain, xlog)
	} else if blockchain == chain.GetChainCode("BSC", xlog) {
		return bnb.NewNftBnb(cluster, blockchain, xlog)
	}
	return nil
}

func NewApis(clusters map[int64][]*config.NodeCluster, xlog *xlog.XLog) map[int64]blockchain.API {
	blockChainClients := make(map[int64]blockchain.API, 0)
	for chainCode, cluster := range clusters {
		blockChainClients[chainCode] = NewApi(chainCode, cluster, xlog)
	}
	return blockChainClients
}

func NewNftApis(clusters map[int64][]*config.NodeCluster, xlog *xlog.XLog) map[int64]blockchain.NftApi {
	blockChainClients := make(map[int64]blockchain.NftApi, 0)
	for chainCode, cluster := range clusters {
		blockChainClients[chainCode] = NewNftApi(chainCode, cluster, xlog)
	}
	return blockChainClients
}
