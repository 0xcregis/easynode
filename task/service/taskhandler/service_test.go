package taskhandler

import (
	"github.com/uduncloud/easynode/task/config"
	"log"
	"testing"
	"time"
)

func Init() *Service {
	cfg := config.LoadConfig("./../../config.json")
	return NewService(&cfg)
}

func TestService_GetTaskForExec(t *testing.T) {
	s := Init()
	log.Println(s.GetTaskForExec([]string{"200"}))
}

func TestService_GetNodeInfo(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeInfo())
}

func TestService_Start(t *testing.T) {
	s := Init()
	s.Start()
	time.Sleep(10 * time.Minute)
}

func TestBalanceForCluster(t *testing.T) {
	//s := Init()
	s := Service{}
	s.balanceForCluster([]string{"0ea23437-ccd9-4abb-a528-abd1e28182f0#1", "74f14233-9948-4a74-ae47-bc47ffe1c246#1", "ca11057d-add8-4501-b6ea-11d34d6df675#1"}, 3473727)
}

func TestService_RebuildWeight(t *testing.T) {
	s := Init()
	log.Println(s.RebuildWeight("0ea23437-ccd9-4abb-a528-abd1e28182f0", 1))
}

func TestService_GetNodeTaskCountWithIng(t *testing.T) {
	s := Init()
	nodeInfoMp, chainCodeList, err := s.GetNodeInfo()
	if err != nil {
		panic(err)
	}
	ss, mp := s.GetNodeTaskCountWithIng(chainCodeList, nodeInfoMp, nil)
	log.Println(ss, mp)
}

func TestService_GetNodeTaskCountWithMap(t *testing.T) {
	s := Init()
	log.Println(s.GetNodeTaskCountWithMap())
}
