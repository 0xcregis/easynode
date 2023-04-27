package config

type Config struct {
	RootPath     string                  `json:"RootPath"`
	Port         int                     `json:"Port"`
	BlockChain   []int64                 `json:"BlockChain"`
	ClickhouseDb map[int64]*ClickhouseDb `json:"ClickhouseDb"`
	TaskKafka    *Kafka                  `json:"TaskKafka"`
}

type Kafka struct {
	Host      string `json:"Host"`
	Port      int    `json:"Port"`
	Topic     string `json:"Topic"`
	Partition int    `json:"Partition"`
}

type TaskDb struct {
	Addr     string `json:"Addr"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DbName   string `json:"DbName"`
}

type ClickhouseDb struct {
	Addr         string `json:"Addr"`
	Port         int    `json:"Port"`
	User         string `json:"User"`
	Password     string `json:"Password"`
	DbName       string `json:"DbName"`
	TxTable      string `json:"TxTable"`
	BlockTable   string `json:"BlockTable"`
	ReceiptTable string `json:"ReceiptTable"`
}
