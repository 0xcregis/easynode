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
	if chain.GetChainCode(blockchain, "ETH", xlog) {
		return ether.NewEth(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "TRON", xlog) {
		return tron.NewTron(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "POLYGON", xlog) {
		return polygon.NewPolygonPos(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "BSC", xlog) {
		return bnb.NewBnb(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "BTC", xlog) {
		return btc.NewBtc(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "FIL", xlog) {
		return filecoin.NewFileCoin(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "XRP", xlog) {
		return xrp.NewXRP(cluster, blockchain, xlog)
	}
	return nil
}

func NewNftApi(blockchain int64, cluster []*config.NodeCluster, xlog *xlog.XLog) blockchain.NftApi {
	if chain.GetChainCode(blockchain, "ETH", xlog) {
		return ether.NewNftEth(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "POLYGON", xlog) {
		return polygon.NewNftPolygonPos(cluster, blockchain, xlog)
	} else if chain.GetChainCode(blockchain, "BSC", xlog) {
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

func NewTronApi(clusters map[int64][]*config.NodeCluster, xlog *xlog.XLog) *tron.Tron {
	for chainCode, cluster := range clusters {
		if chain.GetChainCode(chainCode, "TRON", xlog) {
			return tron.NewTron2(cluster, chainCode, xlog)
		}
	}
	return nil
}
