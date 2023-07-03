package config

type Config struct {
	RootPath string            `json:"RootPath"`
	Port     int               `json:"Port"`
	Chains   []*Chain          `json:"Chains"`
	BaseDb   *ClickhouseBaseDb `json:"BaseDb"`
}

type Chain struct {
	BlockChain   int64                   `json:"BlockChain"`
	TxStore      bool                    `json:"TxStore"`
	BlockStore   bool                    `json:"BlockStore"`
	ReceiptStore bool                    `json:"ReceiptStore"`
	SubStore     bool                    `json:"SubStore"`
	KafkaCfg     map[string]*KafkaConfig `json:"Kafka"`
	ClickhouseDb *ClickhouseChainDb      `json:"ChainDb"`
	Redis        *Redis                  `json:"Redis"`
}

type Redis struct {
	Addr string `json:"Addr"`
	Port int64  `json:"Port"`
	DB   int    `json:"DB"`
}

type ClickhouseChainDb struct {
	Addr         string `json:"Addr"`
	Port         int    `json:"Port"`
	User         string `json:"User"`
	Password     string `json:"Password"`
	DbName       string `json:"DbName"`
	TxTable      string `json:"TxTable"`
	BlockTable   string `json:"BlockTable"`
	ReceiptTable string `json:"ReceiptTable"`
	SubTxTable   string `json:"SubTxTable"`
}

type ClickhouseBaseDb struct {
	Addr         string `json:"Addr"`
	Port         int    `json:"Port"`
	User         string `json:"User"`
	Password     string `json:"Password"`
	DbName       string `json:"DbName"`
	AddressTable string `json:"AddressTable"`
	TokenTable   string `json:"TokenTable"`
}

type KafkaConfig struct {
	Port        int    `json:"Port" `
	Host        string `json:"Host"`
	Topic       string `json:"Topic"`
	Partition   int    `json:"Partition"`
	Group       string `json:"Group"`
	StartOffset int64  `json:"StartOffset"` //-1:latest,-2:first,0:commited offset
}

type Topic struct {
	Topic     string `json:"Topic"`
	Partition int    `json:"Partition"`
}
