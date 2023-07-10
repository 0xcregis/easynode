package service

import (
	"context"
	"github.com/segmentio/kafka-go"
)

type Process interface {
	Start(ctx context.Context)
}

type BlockChainInterface interface {
	GetLatestBlockNumber() (int64, error)
}

type StoreTaskInterface interface {
	AddNodeTask(list []*NodeTask) ([]*kafka.Message, error)
	UpdateLastNumber(blockChainCode int64, latestNumber int64) error
	UpdateRecentNumber(blockChainCode int64, recentNumber int64) error
	GetRecentNumber(blockCode int64) (int64, int64, error)
	GetNodeId(blockChainCode int64) ([]string, error)
	GetAndDelNodeId(blockChainCode int64) ([]string, error)
}
