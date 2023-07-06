package taskcreate

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"github.com/uduncloud/easynode/task/service/taskcreate/ether"
	"github.com/uduncloud/easynode/task/service/taskcreate/tron"
)

func NewApi(blockchain int64, log *xlog.XLog, v *config.BlockConfig) service.BlockChainInterface {
	if blockchain == 200 {
		return ether.NewEther(log, v)
	} else if blockchain == 205 {
		return tron.NewTron(log, v)
	}
	return nil
}
