package xrp

import (
	"encoding/json"

	"github.com/tidwall/gjson"
)

func GetBlockHead(block string) string {
	mp := gjson.Parse(block).Map()
	delete(mp, "transactions")
	r := make(map[string]any, len(mp))
	for k, v := range mp {
		r[k] = v.String()
	}
	bs, _ := json.Marshal(r)
	return string(bs)
}
