package db

import (
	"fmt"
	"github.com/sunjiangjun/xlog"
	"gorm.io/driver/clickhouse"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

func OpenCK(user, password, addr, dbName string, port int, xlog *xlog.XLog) (*gorm.DB, error) {

	//mysql
	//dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8mb4,utf8&parseTime=True&loc=%s",
	//	user,
	//	password,
	//	addr,
	//	port,
	//	dbName,
	//	"Asia%2FShanghai")
	//dsn := "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	//PostgresSQL
	//dsn := "host=%v user=%v password=%v dbname=%v port=%v sslmode=disable TimeZone=Asia/Shanghai"
	//dsn = fmt.Sprintf(dsn, addr, user, password, dbName, port)

	dsn := "clickhouse://%v:%v@%v:%v/%v"

	dsn = fmt.Sprintf(dsn, user, password, addr, port, dbName)

	l := logger.New(log.New(xlog.Out, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.LogLevel(xlog.Level),
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})

	db, err := gorm.Open(clickhouse.Open(dsn), &gorm.Config{
		Logger: l,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return nil, err
	}

	d, _ := db.DB()
	d.SetMaxOpenConns(30)
	d.SetConnMaxLifetime(7 * time.Hour)
	d.SetMaxIdleConns(20)
	d.SetConnMaxIdleTime(70 * time.Minute)
	return db, nil
}
