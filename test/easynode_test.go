package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/coreos/etcd/pkg/fileutil"
)

func TestCheckConfig(t *testing.T) {
	p := "./../build"
	dirList, err := os.ReadDir(p)
	if err != nil {
		panic(err)
	}
	for _, v := range dirList {
		if v.IsDir() && v.Name() == "config" {
			configPath := filepath.Join(p, v.Name())
			t.Log(fileutil.ReadDir(configPath))
		}

		if v.IsDir() && v.Name() == "scripts" {
			configPath := filepath.Join(p, v.Name())
			t.Log(fileutil.ReadDir(configPath))
		}
	}
}
