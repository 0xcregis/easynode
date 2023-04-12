package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadConfig(path string) Config {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	json.Unmarshal(b, &cfg)
	return cfg
}
