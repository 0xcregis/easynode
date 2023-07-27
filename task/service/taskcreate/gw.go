package taskcreate

import (
	"github.com/0xcregis/easynode/task/config"
	"github.com/0xcregis/easynode/task/service"
	"github.com/0xcregis/easynode/task/service/taskcreate/ether"
	"github.com/0xcregis/easynode/task/service/taskcreate/polygonpos"
	"github.com/0xcregis/easynode/task/service/taskcreate/tron"
	"github.com/sunjiangjun/xlog"
)

func NewApi(blockchain int64, log *xlog.XLog, v *config.BlockConfig) service.BlockChainInterface {
	if blockchain == 200 {
		return ether.NewEther(log, v)
	} else if blockchain == 205 {
		return tron.NewTron(log, v)
	} else if blockchain == 201 {
		return polygonpos.NewPolygonPos(log, v)
	}
	return nil
}
