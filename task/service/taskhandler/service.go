package taskhandler

import (
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
	"github.com/uduncloud/easynode/task/common/sql"
	"github.com/uduncloud/easynode/task/config"
	"github.com/uduncloud/easynode/task/service"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	config       *config.Config
	nodeSourceDb *gorm.DB
	taskDb       *gorm.DB
	nodeInfoDb   *gorm.DB
	log          *xlog.XLog
	nodeMp       map[string]int64
}

func (s *Service) Start() {
	go func() {
		s.initNodeMap()
	}()

	go func() {
		s.loop()
	}()
}

func (s *Service) initNodeMap() {
	for true {
		<-time.After(5 * time.Second)
		mp := s.GetNodeTaskCountWithMap()
		if mp != nil && len(mp) > 0 {
			s.nodeMp = mp
		}
	}
}

func (s *Service) loop() {
	for true {
		log := s.log.WithFields(logrus.Fields{
			"id":    time.Now().UnixMilli(),
			"model": "Start",
		})

		//读取节点信息
		nodeInfoMp, chainCodeList, err := s.GetNodeInfo()
		if err != nil {
			log.Errorf("GetNodeInfo|err=%v", err)
			<-time.After(10 * time.Second)
			continue
		}

		if len(nodeInfoMp) < 1 {
			log.Warnf("nodeInfoMp|err=%v", errors.New("has not found node info"))
			<-time.After(10 * time.Second)
			continue
		}

		//根据待执行和正在执行的任务数量，判断
		chainCodeList, nodeInfoMp = s.GetNodeTaskCountWithIng(chainCodeList, nodeInfoMp, log)
		if len(chainCodeList) < 1 {
			log.Warnf("chainCodeList|err=%v", "chainCode list is empty")
			<-time.After(10 * time.Second)
			continue
		}

		log.Printf("active chainCode:%v", chainCodeList)
		log.Printf("active nodeInfo:%v", nodeInfoMp)

		//读取任务
		taskList, err := s.GetTaskForExec(chainCodeList)
		if err != nil {
			log.Errorf("GetTaskForExec|err=%v", err.Error())
			<-time.After(10 * time.Second)
			continue
		}

		//暂没有需要执行的任务
		if len(taskList) < 1 {
			log.Warnf("taskList|info=%v", "taskList is empty")
			time.Sleep(10 * time.Second)
			continue
		}

		log.Printf("active taskList:%v", len(taskList))

		//平衡算法
		//下发任务
		nodeTaskList := make([]*service.NodeTask, 0)
		nodeTaskIdList := make([]int64, 0)

		for _, v := range taskList {
			sp := v.SourceType
			code := v.BlockChain
			var key string
			var taskType int8
			if sp == 1 {
				key = fmt.Sprintf("tx_%v", code)
				taskType = 1
			}

			if sp == 2 {
				key = fmt.Sprintf("block_%v", code)
				taskType = 2
			}

			if sp == 3 {
				key = fmt.Sprintf("receipt_%v", code)
				taskType = 3
			}

			nodeIds := nodeInfoMp[key]
			if len(nodeIds) < 1 {
				s.log.Printf("nodeInfoMp|key=%v|err=%v", key, errors.New("has not found node info"))
				continue
			}
			//平衡算法
			var nodeId string
			nodeId = s.balanceForCluster(nodeIds, float64(v.Id))

			task := service.NodeTask{TaskStatus: 0, BlockChain: code, TaskType: taskType, TxHash: v.TxHash, BlockNumber: v.BlockNumber, BlockHash: v.BlockHash, NodeId: nodeId}
			nodeTaskList = append(nodeTaskList, &task)
			nodeTaskIdList = append(nodeTaskIdList, v.Id)
		}

		if len(nodeTaskList) > 0 {
			err = s.AddNodeTask(nodeTaskList)
			if err == nil && len(nodeTaskIdList) > 0 {
				_ = s.DeleteNodeSource(nodeTaskIdList)
			}
		}

	}
}

// RebuildWeight 根据积压的任务，动态的调整 最终权重
func (s *Service) RebuildWeight(nodeId string, weight int64) int64 {

	//检查缓存空间，减少数据库频繁的读取
	var count int64
	if c, ok := s.nodeMp[nodeId]; !ok {
		count = s.GetNodeTaskCountByNodeId(nodeId)
		s.nodeMp[nodeId] = count
	} else {
		count = c
	}

	//根据其他维度，调整权重
	if count < 10 {
		return weight
	} else if count >= 10 && count <= 100 {
		return weight / 2
	} else if count > 100 {
		return weight / 10
	}
	return weight
}

func (s *Service) balanceForCluster(nodes []string, shardId float64) string {
	if len(nodes) == 1 {
		str := nodes[0]
		//array[0]:node_id, array[1]:weight
		array := strings.Split(str, "#")
		return array[0]
	}

	mp := make(map[string][]int64, 0)
	for _, v := range nodes {
		array := strings.Split(v, "#")
		if len(array) == 2 {
			key := array[0]
			weight, _ := strconv.ParseInt(array[1], 0, 64)
			//根据其他维度，调整最终权重计算
			weight = s.RebuildWeight(key, weight)
			mp[key] = []int64{weight}
		}
	}

	var sum int64
	for k, v := range mp {
		sum += v[0]
		mp[k] = append(v, sum)
	}

	f := math.Mod(shardId, float64(sum))
	var nodeId string

	for k, v := range mp {
		if len(v) == 2 && f <= float64(v[1]) && f >= float64(v[1]-v[0]) {
			nodeId = k
			break
		}
	}

	if len(nodeId) < 1 {
		nodeId = nodes[0]
	}

	return nodeId
}
func (s *Service) getNodeTaskTable() string {
	table := fmt.Sprintf("%v_%v", s.config.NodeTaskDb.Table, time.Now().Format(service.DayFormat))
	return table
}
func (s *Service) AddNodeTask(list []*service.NodeTask) error {
	err := s.taskDb.Table(s.getNodeTaskTable()).Omit("id,create_time,log_time").CreateInBatches(list, 10).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) GetNodeTaskCountWithIng(codes []string, nodeMp map[string][]string, log *logrus.Entry) ([]string, map[string][]string) {

	mp := make(map[string]int64, 0)
	for _, v := range codes {
		mp[v] = 1
	}

	returnCode := make([]string, 0, 10)
	returnNode := make(map[string][]string, 0)

	for k, vs := range nodeMp {

		temp := make([]string, 0, 10)
		for _, v := range vs {
			ks := strings.Split(v, "#")
			nodeId := ks[0]

			var Num int64
			err := s.taskDb.Table(s.getNodeTaskTable()).Where("task_status in (0,4,3) and node_id=?", nodeId).Count(&Num).Error
			if err != nil {
				log.Errorf("GetNodeTaskCountWithIng|blockChain=%v | error=%v", k, err)
				continue
			}
			if Num < 200 {
				//r = append(r, k)
				temp = append(temp, v)
			} else {
				log.Printf("blockChain=%v | This task assignment is cancelled ,becase many tasks waiting to execute", k)
			}
		}

		if len(temp) > 0 {
			returnNode[k] = temp
		}

	}

	for k, _ := range mp {
		//var Num int64
		//err := s.taskDb.Table(s.getNodeTaskTable()).Where("task_status in (0,4,3) and block_chain=?", k).Count(&Num).Error
		//if err != nil {
		//	s.log.Printf("GetNodeTaskCountWithIng|blockChain=%v | error=%v", k, err)
		//	continue
		//}
		//if Num < 200 {
		returnCode = append(returnCode, k)
		//} else {
		//	s.log.Printf("blockChain=%v | This task assignment is cancelled ,becase many tasks waiting to execute", k)
		//}
	}

	return returnCode, returnNode
}

func (s *Service) GetNodeTaskCountWithMap() map[string]int64 {
	var list []*struct {
		NodeId string
		Number int64
	}
	err := s.taskDb.Table(s.getNodeTaskTable()).Select("node_id,count(1)as number").Where("task_status in (0,4,3)").Group("node_id").Scan(&list).Error
	if err != nil {
		return nil
	}

	mp := make(map[string]int64, 2)
	for _, v := range list {
		mp[v.NodeId] = v.Number
	}
	return mp
}

func (s *Service) GetNodeTaskCountByNodeId(nodeId string) int64 {
	var Num int64
	err := s.taskDb.Table(s.getNodeTaskTable()).Where("task_status in (0,4,3) and node_id=?", nodeId).Count(&Num).Error
	if err != nil {
		return 0
	}
	return Num
}

func (s *Service) GetNodeInfo() (map[string][]string, []string, error) {
	var list []*service.NodeInfo
	start := time.Now().Add(-5 * time.Minute).Format(service.TimeFormat)
	err := s.nodeInfoDb.Table(s.config.NodeInfoDb.Table).Where("create_time>?", start).Scan(&list).Error
	if err != nil {
		return nil, nil, err
	}

	mp := make(map[string][]string, 0)
	chainCode := make([]string, 0, 10)
	for _, v := range list {
		nodeId := v.NodeId
		array := gjson.Parse(v.Info).Array()
		for _, r := range array {
			code := r.Get("BlockChainCode").String()
			weight := r.Get("NodeWeight").Int()
			if weight < 10 {
				weight = 10
			}

			//采集节点支持哪些公链
			chainCode = append(chainCode, code)

			if r.Get("TxTask").Exists() {
				key := fmt.Sprintf("tx_%v", code)
				mp[key] = append(mp[key], fmt.Sprintf("%v#%v", nodeId, weight))
			}
			if r.Get("BlockTask").Exists() {
				key := fmt.Sprintf("block_%v", code)
				mp[key] = append(mp[key], fmt.Sprintf("%v#%v", nodeId, weight))
			}
			if r.Get("ReceiptTask").Exists() {
				key := fmt.Sprintf("receipt_%v", code)
				mp[key] = append(mp[key], fmt.Sprintf("%v#%v", nodeId, weight))
			}

		}

	}

	return mp, chainCode, nil
}

func (s *Service) GetTaskForExec(codes []string) ([]*service.NodeSource, error) {

	//codes 值可能出现重复，暂忽略不及，不影响SQL 性能

	r := make([]*service.NodeSource, 0, 100)

	for _, v := range codes {

		//区块任务
		var blockList []*service.NodeSource
		err := s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Where("source_type=? and block_chain =?", 2, v).Order("create_time asc").Limit(5).Scan(&blockList).Error
		if err == nil && len(blockList) > 0 {
			r = append(r, blockList...)
		} else if err != nil {
			s.log.Printf("GetTaskForExec|error=%v", err)
		}

		temps := make([]string, 0, len(blockList))
		for _, v := range blockList {
			temps = append(temps, v.BlockNumber)
		}

		//交易任务
		var list []*service.NodeSource
		if len(temps) > 0 {
			err = s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Where("source_type=? and block_number in (?) and block_chain =?", 1, temps, v).Order("create_time asc").Scan(&list).Error
		}
		if len(list) < 1 {
			err = s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Where("source_type=? and block_chain =?", 1, v).Order("create_time asc").Limit(200).Scan(&list).Error
		}
		if err != nil {
			s.log.Printf("GetTaskForExec|error=%v", err)
		}
		if len(list) > 0 {
			r = append(r, list...)
		}

		//收据任务
		var list3 []*service.NodeSource
		if len(temps) > 0 {
			err = s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Where("source_type=? and block_number in (?) and block_chain =?", 3, temps, v).Order("create_time asc").Scan(&list3).Error
		}
		if len(list3) < 1 {
			err = s.nodeSourceDb.Table(s.config.NodeSourceDb.Table).Where("source_type=? and block_chain =?", 3, v).Order("create_time asc").Limit(200).Scan(&list3).Error
		}

		if err != nil {
			s.log.Printf("GetTaskForExec|error=%v", err)
		}

		if len(list3) > 0 {
			r = append(r, list3...)
		}
	}
	return r, nil
}

func (s *Service) DeleteNodeSource(list []int64) error {
	str := "delete from %v where id in (?)"
	str = fmt.Sprintf(str, s.config.NodeSourceDb.Table)
	return s.nodeSourceDb.Exec(str, list).Error
}

func NewService(config *config.Config) *Service {
	xg := xlog.NewXLogger().BuildOutType(xlog.FILE).BuildFile("./log/task/handler_task", 24*time.Hour)
	s, err := sql.Open(config.NodeSourceDb.User, config.NodeSourceDb.Password, config.NodeSourceDb.Addr, config.NodeSourceDb.DbName, config.NodeSourceDb.Port, xg)
	if err != nil {
		panic(err)
	}

	info, err := sql.Open(config.NodeInfoDb.User, config.NodeInfoDb.Password, config.NodeInfoDb.Addr, config.NodeInfoDb.DbName, config.NodeInfoDb.Port, xg)
	if err != nil {
		panic(err)
	}

	task, err := sql.Open(config.NodeTaskDb.User, config.NodeTaskDb.Password, config.NodeTaskDb.Addr, config.NodeTaskDb.DbName, config.NodeTaskDb.Port, xg)
	if err != nil {
		panic(err)
	}

	return &Service{
		config:       config,
		nodeSourceDb: s,
		nodeInfoDb:   info,
		taskDb:       task,
		log:          xg,
		nodeMp:       make(map[string]int64, 2),
	}
}
