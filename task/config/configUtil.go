package config

import (
	"encoding/json"
	"io"
	"os"
)

func LoadConfig(path string) Config {
	f, err := os.OpenFile(path, os.O_RDONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	json.Unmarshal(b, &cfg)
	return cfg
}
