package chain

import (
	"encoding/json"
	"io"
	"os"

	"github.com/sunjiangjun/xlog"
)

var defaultChain = `
  {
  "ETH": 200,
  "POLYGON": 201,
  "BSC": 202,
  "TRON": 205,
  "BTC": 300,
  "FIL": 301,
  "XRP": 310
}
`

func LoadConfig(path string) (string, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = f.Close()
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func LoadChainCodeFile(file string) map[string]int64 {
	//default
	body := defaultChain

	//set customer config
	if len(file) > 1 {
		f, _ := LoadConfig(file)
		if len(f) > 1 {
			body = f
		}
	}

	if len(body) < 1 {
		return nil
	}

	mp := make(map[string]int64)
	err := json.Unmarshal([]byte(body), &mp)
	if err != nil {
		return nil
	}
	return mp
}

func GetChainCode(chainName string, log *xlog.XLog) int64 {
	if log == nil {
		log = xlog.NewXLogger()
	}
	mp := LoadChainCodeFile("./chain.json")
	if mp == nil {
		log.Errorf("unknown all chainCode，this is a fatal error")
		return 0
	}

	if code, ok := mp[chainName]; ok {
		return code
	} else {
		log.Errorf("unknown chainCode:%v，please check whether the system supports this chain", code)
		return 0
	}
}
