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
	defer func() {
		_ = f.Close()
	}()
	b, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	cfg := Config{}
	_ = json.Unmarshal(b, &cfg)
	return cfg
}
