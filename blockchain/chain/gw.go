package chain

import (
	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/bnb"
	"github.com/0xcregis/easynode/blockchain/chain/btc"
	"github.com/0xcregis/easynode/blockchain/chain/ether"
	"github.com/0xcregis/easynode/blockchain/chain/filecoin"
	"github.com/0xcregis/easynode/blockchain/chain/polygonpos"
	"github.com/0xcregis/easynode/blockchain/chain/tron"
	"github.com/0xcregis/easynode/blockchain/chain/xrp"
	"github.com/0xcregis/easynode/common/chain"
	"github.com/sunjiangjun/xlog"
)

func NewChain(blockchain int64, log *xlog.XLog) blockchain.ChainConn {
	if chain.GetChainCode(blockchain, "ETH", log) {
		//eth
		return ether.NewChainClient()
	} else if chain.GetChainCode(blockchain, "TRON", log) {
		//tron
		return tron.NewChainClient()
	} else if chain.GetChainCode(blockchain, "POLYGON", log) {
		//polygon-pos
		return polygonpos.NewChainClient()
	} else if chain.GetChainCode(blockchain, "BSC", log) {
		//bnb
		return bnb.NewChainClient()
	} else if chain.GetChainCode(blockchain, "FIL", log) {
		//file-coin
		return filecoin.NewChainClient()
	} else if chain.GetChainCode(blockchain, "BTC", log) {
		//btc
		return btc.NewChainClient()
	} else if chain.GetChainCode(blockchain, "XRP", log) {
		//xrp
		return xrp.NewChainClient()
	} else {
		return nil
	}
}

func NewNFT(blockchain int64, log *xlog.XLog) blockchain.NFT {
	if chain.GetChainCode(blockchain, "ETH", log) {
		//eth
		return ether.NewNFTClient()
	} else if chain.GetChainCode(blockchain, "POLYGON", log) {
		return polygonpos.NewNFTClient()
	} else if chain.GetChainCode(blockchain, "BSC", log) {
		return bnb.NewNFTClient()
	} else {
		return nil
	}
}
