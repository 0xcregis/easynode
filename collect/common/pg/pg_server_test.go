package pg

import (
	"github.com/sunjiangjun/xlog"
	"log"
	"testing"
)

func TestOpen(t *testing.T) {
	x := xlog.NewXLogger().BuildOutType(xlog.STD)
	g, err := Open("root", "123456789", "192.168.2.11", "information_schema", 3306, x)
	if err != nil {
		panic(err)
	}
	db, err := g.DB()
	if err != nil {
		panic(err)
	}
	log.Println(db.Ping())
}
