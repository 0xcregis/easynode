package config

type Config struct {
	NodeTaskDb          *NodeTaskDb    `json:"NodeTaskDb"`
	AutoCreateBlockTask bool           `json:"AutoCreateBlockTask"`
	BlockConfigs        []*BlockConfig `json:"Chains"`
}

type NodeTaskDb struct {
	User     string `json:"User" gorm:"column:User"`
	Table    string `json:"Table" gorm:"column:Table"`
	Port     int    `json:"Port" gorm:"column:Port"`
	DbName   string `json:"DbName" gorm:"column:DbName"`
	Addr     string `json:"Addr" gorm:"column:Addr"`
	Password string `json:"Password" gorm:"column:Password"`
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
	BlockChainName string `json:"BlockChainName"`
	BlockMin       int64  `json:"BlockMin"`
	BlockChainCode int64  `json:"BlockChainCode"`
	BlockMax       int64  `json:"BlockMax"`
	NodeKey        string `json:"NodeKey"`
	NodeHost       string `json:"NodeHost"`
}
