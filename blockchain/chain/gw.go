package chain

import (
	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/ether"
	"github.com/0xcregis/easynode/blockchain/chain/filecoin"
	"github.com/0xcregis/easynode/blockchain/chain/polygonpos"
	"github.com/0xcregis/easynode/blockchain/chain/tron"
	"github.com/0xcregis/easynode/blockchain/chain/xrp"
)

func NewChain(blockchain int64) blockchain.ChainConn {
	if blockchain == 200 {
		//eth
		return ether.NewChainClient()
	} else if blockchain == 205 {
		//tron
		return tron.NewChainClient()
	} else if blockchain == 201 {
		//polygon-pos
		return polygonpos.NewChainClient()
	} else if blockchain == 301 {
		//file-coin
		return filecoin.NewChainClient()
	} else if blockchain == 300 {
		//btc
		return nil
	} else if blockchain == 310 {
		//xrp
		return xrp.NewChainClient()
	} else {
		return nil
	}
}

func NewNFT(blockchain int64) blockchain.NFT {
	if blockchain == 200 {
		//eth
		return ether.NewNFTClient()
	} else {
		return nil
	}
}
