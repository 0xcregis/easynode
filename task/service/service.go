package service

type Process interface {
	Start()
}

type StoreTaskInterface interface {
	AddNodeTask(list []*NodeTask) error
	UpdateLastNumber(blockChainCode int64, latestNumber int64) error
	UpdateRecentNumber(blockChainCode int64, recentNumber int64) error
	GetRecentNumber(blockCode int64) (int64, int64, error)
}
