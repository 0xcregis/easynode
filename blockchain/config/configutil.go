package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/tidwall/gjson"
)

func LoadConfig(path string) Config {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = f.Close()
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(err)
	}

	list := gjson.ParseBytes(b).Get("Nodes").Array()
	for _, v := range list {
		chainCode := v.Get("BlockChain").Int()
		var node NodeCluster
		err := json.Unmarshal([]byte(v.String()), &node)
		if err != nil {
			panic(err)
		}

		if m, ok := cfg.Cluster[chainCode]; ok {
			m = append(m, &node)
			cfg.Cluster[chainCode] = m
		} else {
			m := make(map[int64][]*NodeCluster, 1)
			nodeList := []*NodeCluster{&node}
			m[chainCode] = nodeList
			cfg.Cluster = m
		}
	}

	return cfg
}
