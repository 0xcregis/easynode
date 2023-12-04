package chain

import (
	"log"
	"testing"

	"github.com/sunjiangjun/xlog"
)

func TestGetChainCode(t *testing.T) {
	code := GetChainCode(200, "ETH", xlog.NewXLogger())
	log.Println(code)
}
