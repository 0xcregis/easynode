package config

type TaskDb struct {
	Addr     string `json:"Addr"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DbName   string `json:"DbName"`
	Table    string `json:"Table"`
}

type FromCluster struct {
	Host       string `json:"Host"`
	Port       int    `json:"Port"`
	Key        string `json:"Key"`
	Weight     int64  `json:"Weight"`
	ErrorCount int64  `json:"ErrorCount"`
}
type Kafka struct {
	Host        string  `json:"Host"`
	Port        int     `json:"Port"`
	Topic       string  `json:"Topic"`
	Partition   int     `json:"Partition"`
	Partitions  []int64 `json:"WritePartitions"`
	Group       string  `json:"Group"`
	StartOffset int64   `json:"StartOffset"`
}

type BlockTask struct {
	FromCluster []*FromCluster `json:"FromCluster"`
	Kafka       *Kafka         `json:"Kafka"`
}

type TxTask struct {
	FromCluster []*FromCluster `json:"FromCluster"`
	Kafka       *Kafka         `json:"Kafka"`
}

type ReceiptTask struct {
	FromCluster []*FromCluster `json:"FromCluster"`
	Kafka       *Kafka         `json:"Kafka"`
}

type Chain struct {
	//Etcd           *Etcd        `json:"Etcd"`
	BlockChainName string `json:"BlockChainName"`
	BlockChainCode int    `json:"BlockChainCode"`
	//NodeWeight     int          `json:"NodeWeight"`
	PullReceipt uint8        `json:"PullReceipt"` //0:没有收据任务，1:批量收据任务，2:单笔收据任务
	PullTx      uint8        `json:"PullTx"`      //0:没有交易任务，1:批量交易任务，2:单笔交易任务
	Kafka       *Kafka       `json:"Kafka"`       //结果数据Kafka
	TaskKafka   *Kafka       `json:"TaskKafka"`   //任务kafka
	BlockTask   *BlockTask   `json:"BlockTask"`
	TxTask      *TxTask      `json:"TxTask"`
	ReceiptTask *ReceiptTask `json:"ReceiptTask"`
	Redis       *Redis       `json:"Redis"`
}

type Redis struct {
	Addr string `json:"Addr"`
	Port int64  `json:"Port"`
	DB   int    `json:"DB"`
}

type LogConfig struct {
	Path  string `json:"Path"`
	Delay int64  `json:"Delay"`
}

type Config struct {
	Chains    []*Chain   `json:"Chains"`
	LogConfig *LogConfig `json:"Log"`
	KeyPath   string     `json:"KeyPath"`
}
