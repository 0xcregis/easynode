package util

import (
	"log"
	"testing"
)

func TestGetLocalMachine(t *testing.T) {
	log.Println(GetLocalNodeId())
}

func TestDeleteFile(t *testing.T) {

	log.Println(DeleteFile("temp/log"))
}
