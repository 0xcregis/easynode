package db

import (
	"errors"
	"fmt"
	"github.com/sunjiangjun/xlog"
	"github.com/uduncloud/easynode/collect/common/pg"
	"github.com/uduncloud/easynode/collect/config"
	"github.com/uduncloud/easynode/collect/service"
	"gorm.io/gorm"
	"sync"
	"time"
)

const (
	DayFormat = "20060102"
)

type Service struct {
	log     *xlog.XLog
	taskOrm *gorm.DB
	taskDb  *config.TaskDb
	lock    *sync.RWMutex
}

func (s *Service) UpdateTaskStatus(id int64, status int) error {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", s.taskDb.Table, time.Now().Format(DayFormat))
	return table
}

func (s *Service) GetNodeTaskWithTxs(txHash []string, taskType int, blockChain int, taskStatus int) ([]int64, error) {
	var ids []int64
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Select("id").Where("tx_hash in (?) and task_type=? and block_chain=? and task_status=?", txHash, taskType, blockChain, taskStatus).Order("create_time desc").Pluck("id", &ids).Error
	if err != nil || len(ids) < 1 {
		s.log.Printf("GetNodeTaskStatusWithTx|error=%v", "no record")
		return nil, errors.New("no record")
	}
	return ids, nil
}

func (s *Service) GetNodeTaskByBlockNumber(number string, taskType int, blockChain int) (*service.NodeTask, error) {
	if len(number) < 1 {
		return nil, errors.New("no record")
	}
	var task service.NodeTask
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("block_number=? and task_type=? and block_chain=? and task_status=?", number, taskType, blockChain, 4).Order("create_time desc").Limit(1).First(&task).Error
	if err != nil {
		return nil, err
	}
	if task.Id < 1 {
		return nil, errors.New("no record")
	}
	return &task, nil
}

func (s *Service) GetNodeTaskByBlockHash(hash string, taskType int, blockChain int) (*service.NodeTask, error) {
	if len(hash) < 1 {
		return nil, errors.New("no record")
	}
	var task service.NodeTask
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("block_hash=? and task_type=? and block_chain=? and task_status=?", hash, taskType, blockChain, 4).Order("create_time desc").Limit(1).First(&task).Error
	if err != nil {
		return nil, err
	}
	if task.Id < 1 {
		return nil, errors.New("no record")
	}
	return &task, nil
}

func (s *Service) AddNodeTask(list []*service.NodeTask) error {
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Omit("id,create_time,log_time").CreateInBatches(list, 10).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateNodeTaskStatus(id int64, status int) error {
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("id=?", id).UpdateColumn("task_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateNodeTaskStatusWithBatch(ids []int64, status int) error {
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("id in (?)", ids).UpdateColumn("task_status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetTaskWithReceipt(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	//s.lock.Lock()
	//defer func() {
	//	s.lock.Unlock()
	//}()
	var result []*service.NodeTask
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("node_id=? and block_chain=? and task_type=? and task_status=0", nodeId, blockChain, 3).Order("create_time ASC").Limit(100).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, errors.New("no record")
	}

	ids := make([]int64, 0, len(result))
	for _, v := range result {
		ids = append(ids, v.Id)
	}
	_ = s.UpdateNodeTaskStatusWithBatch(ids, 3)

	return result, nil
}

func (s *Service) GetTaskWithTx(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	//s.lock.Lock()
	//defer func() {
	//	s.lock.Unlock()
	//}()
	var result []*service.NodeTask
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("node_id=? and block_chain=? and task_type=? and task_status=0", nodeId, blockChain, 1).Order("create_time ASC").Limit(100).Scan(&result).Error
	if err != nil {
		return nil, err
	}
	if len(result) < 1 {
		return nil, errors.New("no record")
	}
	ids := make([]int64, 0, len(result))
	for _, v := range result {
		ids = append(ids, v.Id)
	}
	_ = s.UpdateNodeTaskStatusWithBatch(ids, 3)
	return result, nil
}

func (s *Service) GetTaskWithBlock(blockChain int, nodeId string) ([]*service.NodeTask, error) {
	//s.lock.Lock()
	//defer func() {
	//	s.lock.Unlock()
	//}()
	var result []*service.NodeTask
	err := s.taskOrm.Table(s.GetNodeTaskTable()).Where("node_id=? and block_chain=? and task_type=? and task_status=0", nodeId, blockChain, 2).Order("create_time ASC").Limit(10).Scan(&result).Error
	if err != nil {
		return nil, err
	}

	if len(result) < 1 {
		return nil, errors.New("no record")
	}
	ids := make([]int64, 0, len(result))
	for _, v := range result {
		ids = append(ids, v.Id)
	}
	_ = s.UpdateNodeTaskStatusWithBatch(ids, 3)

	return result, nil
}

func NewMySQLTaskService(taskDb *config.TaskDb, x *xlog.XLog) service.TaskDbInterface {
	g, err := pg.Open(taskDb.User, taskDb.Password, taskDb.Addr, taskDb.DbName, taskDb.Port, x)
	if err != nil {
		panic(err)
	}

	return &Service{
		log:     x,
		taskOrm: g,
		taskDb:  taskDb,
		lock:    &sync.RWMutex{},
	}
}
