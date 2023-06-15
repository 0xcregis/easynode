package taskcreate

import (
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service/taskcreate/ether"
	"github.com/uduncloud/easynode/task/service/taskcreate/tron"
)

func GetLastBlockNumber(blockchain int64, log *xlog.XLog, v *config.BlockConfig) (int64, error) {
	var lastNumber int64
	var err error
	if blockchain == 200 {
		lastNumber, err = ether.NewEther(log).GetLatestBlockNumber(v)
	} else if blockchain == 205 {
		lastNumber, err = tron.NewTron(log).GetLatestBlockNumber(v)
	}

	if err != nil {
		return 0, err
	}

	return lastNumber, nil
}
