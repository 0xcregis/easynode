package chain

import (
	"io"
	"os"

	"github.com/sunjiangjun/xlog"
	"github.com/tidwall/gjson"
)

// dev
//var defaultChainCode = map[string]map[int64]int8{
//	"ETH":     {200: 1, 2001: 1},
//	"POLYGON": {201: 1, 2011: 1},
//	"BSC":     {202: 1},
//	"TRON":    {205: 1},
//	"BTC":     {300: 1},
//	"FIL":     {301: 1},
//	"XRP":     {310: 1},
//}

// main
var defaultChainCode = map[string]map[int64]int8{
	"ETH":     {60: 1, 6001: 1},
	"POLYGON": {62: 1, 6201: 1},
	"BSC":     {2510: 1, 2610: 1},
	"TRON":    {195: 1, 198: 1},
	"BTC":     {0: 1, 1: 1},
	"FIL":     {2307: 1},
	"XRP":     {144: 1},
}

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

func LoadChainCodeFile(file string) map[string]map[int64]int8 {
	//set customer config
	mp := make(map[string]map[int64]int8)
	if len(file) > 1 {
		body, _ := LoadConfig(file)
		if len(body) > 1 {
			gjson.Parse(body).ForEach(func(key, v gjson.Result) bool {
				k := key.String()
				list := v.Array()
				m := make(map[int64]int8)
				for _, v := range list {
					code := v.Int()
					m[code] = 1
				}
				mp[k] = m
				return true
			})
		}

	}

	return mp
}

func GetChainCode(chainCode int64, chainName string, log *xlog.XLog) bool {
	if log == nil {
		log = xlog.NewXLogger()
	}
	mp := defaultChainCode

	//todo load chainCode if it is necessary, but it is not efficient because It loads configuration files very frequently
	//mp = LoadChainCodeFile("./chain.json")

	if mp == nil {
		log.Errorf("unknown all chainCode，this is a fatal error")
		return false
	}

	if m, ok := mp[chainName]; ok {
		if _, ok := m[chainCode]; ok {
			return true
		} else {
			//log.Errorf("unknown chainCode:%v，please check whether the system supports this chain", chainCode)
			return false
		}
	} else {
		//log.Errorf("unknown chainCode:%v，please check whether the system supports this chain", chainCode)
		return false
	}
}
