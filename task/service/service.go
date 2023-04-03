package service

import "github.com/uduncloud/easynode/task/config"

type Process interface {
	Start()
}

type CreateTask interface {
	GetLastBlockNumber(v *config.BlockConfig) error
}

type DbTaskCreateInterface interface {
	AddNodeTask(list []*NodeTask) error
	UpdateLastNumber(blockChainCode int64, latestNumber int64) error
	UpdateRecentNumber(blockChainCode int64, recentNumber int64) error
	GetRecentNumber(blockCode int64) (int64, int64, error)
}

type DbTaskMonitorInterface interface {
	CheckTable()
	CreateNodeTaskTable()
	RetryTaskForFail()
	HandlerDeadTask()
	HandlerManyFailTask()
}
