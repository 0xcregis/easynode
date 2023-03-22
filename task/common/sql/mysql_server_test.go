package sql

import (
	"github.com/sunjiangjun/xlog"
	"log"
	"testing"
)

func TestOpen(t *testing.T) {
	x := xlog.NewXLogger().BuildOutType(xlog.STD)
	g, err := Open("root", "123456789", "192.168.2.11", "easy_node", 5432, x)
	if err != nil {
		panic(err)
	}
	db, err := g.DB()
	if err != nil {
		panic(err)
	}
	log.Println(db.Ping())
}
