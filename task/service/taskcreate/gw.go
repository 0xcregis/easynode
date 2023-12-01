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
	if blockchain == chain.GetChainCode("ETH", log) {
		return ether.NewEther(log, v)
	} else if blockchain == chain.GetChainCode("TRON", log) {
		return tron.NewTron(log, v)
	} else if blockchain == chain.GetChainCode("POLYGON", log) {
		return polygonpos.NewPolygonPos(log, v)
	} else if blockchain == chain.GetChainCode("BSC", log) {
		return bnb.NewBnb(log, v)
	} else if blockchain == chain.GetChainCode("FIL", log) {
		return filecoin.NewFileCoin(log, v)
	} else if blockchain == chain.GetChainCode("XRP", log) {
		return xrp.NewXRP(log, v)
	} else if blockchain == chain.GetChainCode("BTC", log) {
		return btc.NewBtc(log, v)
	}
	return nil
}
