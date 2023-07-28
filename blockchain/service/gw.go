package service

import (
	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/ether"
	"github.com/0xcregis/easynode/blockchain/chain/polygonpos"
	"github.com/0xcregis/easynode/blockchain/chain/tron"
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func NewApi(blockchain int64, cluster []*config.NodeCluster, xlog *xlog.XLog) blockchain.API {
	if blockchain == 200 {
		blockChainClient := ether.NewChainClient()
		return NewEth(cluster, blockChainClient, xlog)
	} else if blockchain == 205 {
		blockChainClient := tron.NewChainClient()
		return NewTron(cluster, blockChainClient, xlog)
	} else if blockchain == 201 {
		blockChainClient := polygonpos.NewChainClient()
		return NewPolygonPos(cluster, blockChainClient, xlog)
	}
	return nil
}

func NewApis(clusters map[int64][]*config.NodeCluster, xlog *xlog.XLog) map[int64]blockchain.API {
	blockChainClients := make(map[int64]blockchain.API, 0)
	for blockchain, cluster := range clusters {
		blockChainClients[blockchain] = NewApi(blockchain, cluster, xlog)
	}
	return blockChainClients
}
