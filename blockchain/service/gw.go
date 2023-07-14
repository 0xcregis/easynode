package service

import (
	"github.com/0xcregis/easynode/blockchain/config"
	"github.com/sunjiangjun/xlog"
)

func NewApi(blockchain int64, cluster []*config.NodeCluster, xlog *xlog.XLog) API {
	if blockchain == 200 {
		return NewEth(cluster, xlog)
	} else if blockchain == 205 {
		return NewTron(cluster, xlog)
	}
	return nil
}

func NewApis(clusters map[int64][]*config.NodeCluster, xlog *xlog.XLog) map[int64]API {
	blockChainClients := make(map[int64]API, 0)
	for blockchain, cluster := range clusters {
		blockChainClients[blockchain] = NewApi(blockchain, cluster, xlog)
	}
	return blockChainClients
}
