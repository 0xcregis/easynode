package config

type Config struct {
	RootPath   string                   `json:"RootPath"`
	Port       int                      `json:"Port"`
	BlockChain []int64                  `json:"BlockChain"`
	Cluster    map[int64][]*NodeCluster `json:"Cluster"`
	Kafka      *Kafka                   `json:"Kafka"`
}

type Kafka struct {
	Host      string `json:"Host"`
	Port      int    `json:"Port"`
	Topic     string `json:"Topic"`
	Partition int    `json:"Partition"`
}

type NodeCluster struct {
	NodeUrl    string `json:"NodeUrl"`
	NodeToken  string `json:"NodeToken"`
	Weight     int64  `json:"Weight"`
	ErrorCount int64  `json:"ErrorCount"`
}
