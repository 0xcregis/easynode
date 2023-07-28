package chain

import (
	"github.com/0xcregis/easynode/blockchain"
	"github.com/0xcregis/easynode/blockchain/chain/ether"
	"github.com/0xcregis/easynode/blockchain/chain/polygonpos"
	"github.com/0xcregis/easynode/blockchain/chain/tron"
)

func NewChain(blockchain int64) blockchain.BlockChain {
	if blockchain == 200 {
		return ether.NewChainClient()
	} else if blockchain == 205 {
		return tron.NewChainClient()
	} else if blockchain == 201 {
		return polygonpos.NewChainClient()
	} else {
		return nil
	}
}
