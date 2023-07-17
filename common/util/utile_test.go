package util

import (
	"log"
	"testing"
)

func TestGetLocalNodeId(t *testing.T) {
	log.Println(GetLocalNodeId("./temp/log/name"))
}

func TestReadLatestBlock(t *testing.T) {
	log.Println(ReadLatestBlock())
}
