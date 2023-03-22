package config

type Config struct {
	RootPath     string                  `json:"RootPath"`
	Port         int                     `json:"Port"`
	BlockChain   []int64                 `json:"BlockChain"`
	TaskDb       *TaskDb                 `json:"TaskDb"`
	ClickhouseDb map[int64]*ClickhouseDb `json:"ClickhouseDb"`
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
