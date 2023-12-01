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
	if blockchain == chain.GetChainCode("ETH", log) {
		//eth
		return ether.NewChainClient()
	} else if blockchain == chain.GetChainCode("TRON", log) {
		//tron
		return tron.NewChainClient()
	} else if blockchain == chain.GetChainCode("POLYGON", log) {
		//polygon-pos
		return polygonpos.NewChainClient()
	} else if blockchain == chain.GetChainCode("BSC", log) {
		//bnb
		return bnb.NewChainClient()
	} else if blockchain == chain.GetChainCode("FIL", log) {
		//file-coin
		return filecoin.NewChainClient()
	} else if blockchain == chain.GetChainCode("BTC", log) {
		//btc
		return btc.NewChainClient()
	} else if blockchain == chain.GetChainCode("XRP", log) {
		//xrp
		return xrp.NewChainClient()
	} else {
		return nil
	}
}

func NewNFT(blockchain int64, log *xlog.XLog) blockchain.NFT {
	if blockchain == chain.GetChainCode("ETH", log) {
		//eth
		return ether.NewNFTClient()
	} else if blockchain == chain.GetChainCode("POLYGON", log) {
		return polygonpos.NewNFTClient()
	} else if blockchain == chain.GetChainCode("BSC", log) {
		return bnb.NewNFTClient()
	} else {
		return nil
	}
}
