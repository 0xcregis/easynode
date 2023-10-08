package config

type Config struct {
	AutoCreateBlockTask bool           `json:"AutoCreateBlockTask"`
	LogLevel            int            `json:"LogLevel"`
	BlockConfigs        []*BlockConfig `json:"Chains"`
	TaskKafka           *Kafka         `json:"TaskKafka"`
	Redis               *Redis         `json:"Redis"`
}

type NodeCluster struct {
	NodeHost string `json:"NodeHost"`
	NodeKey  string `json:"NodeKey"`
	Weight   int64  `json:"Weight"`
}

type Redis struct {
	Addr string `json:"Addr"`
	Port int64  `json:"Port"`
	DB   int    `json:"DB"`
}

type Kafka struct {
	Host      string `json:"Host"`
	Port      int    `json:"Port"`
	Topic     string `json:"Topic"`
	Partition int    `json:"Partition"`
}

/**
[
    {
      "BlockChain_Name": "eth",
      "BlockChain_Code": 200,
      "BlockMin": 1000,
      "BlockMax": 1200
    }
  ]
*/

type BlockConfig struct {
	BlockChainName string         `json:"BlockChainName"`
	BlockMin       int64          `json:"BlockMin"`
	BlockChainCode int64          `json:"BlockChainCode"`
	BlockMax       int64          `json:"BlockMax"`
	Cluster        []*NodeCluster `json:"Cluster"`
	//NodeKey        string `json:"NodeKey"`
	//NodeHost       string `json:"NodeHost"`
}
