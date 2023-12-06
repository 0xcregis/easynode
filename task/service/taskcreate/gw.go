package taskcreate

import (
	"github.com/0xcregis/easynode/common/chain"
	"github.com/0xcregis/easynode/task"
	"github.com/0xcregis/easynode/task/config"
	"github.com/0xcregis/easynode/task/service/taskcreate/bnb"
	"github.com/0xcregis/easynode/task/service/taskcreate/btc"
	"github.com/0xcregis/easynode/task/service/taskcreate/ether"
	"github.com/0xcregis/easynode/task/service/taskcreate/filecoin"
	"github.com/0xcregis/easynode/task/service/taskcreate/polygonpos"
	"github.com/0xcregis/easynode/task/service/taskcreate/tron"
	"github.com/0xcregis/easynode/task/service/taskcreate/xrp"
	"github.com/sunjiangjun/xlog"
)

func NewApi(blockchain int64, log *xlog.XLog, v *config.BlockConfig) task.BlockChainInterface {
	if chain.GetChainCode(blockchain, "ETH", log) {
		return ether.NewEther(log, v)
	} else if chain.GetChainCode(blockchain, "TRON", log) {
		return tron.NewTron(log, v)
	} else if chain.GetChainCode(blockchain, "POLYGON", log) {
		return polygonpos.NewPolygonPos(log, v)
	} else if chain.GetChainCode(blockchain, "BSC", log) {
		return bnb.NewBnb(log, v)
	} else if chain.GetChainCode(blockchain, "FIL", log) {
		return filecoin.NewFileCoin(log, v)
	} else if chain.GetChainCode(blockchain, "XRP", log) {
		return xrp.NewXRP(log, v)
	} else if chain.GetChainCode(blockchain, "BTC", log) {
		return btc.NewBtc(log, v)
	}
	return nil
}
