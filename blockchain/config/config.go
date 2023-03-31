package config

type Config struct {
	RootPath   string                   `json:"RootPath"`
	Port       int                      `json:"Port"`
	BlockChain []int64                  `json:"BlockChain"`
	Cluster    map[int64][]*NodeCluster `json:"Cluster"`
	TaskDb     *TaskDb                  `json:"TaskDb"`
}

type NodeCluster struct {
	NodeUrl    string `json:"NodeUrl"`
	NodeToken  string `json:"NodeToken"`
	Weight     int64  `json:"Weight"`
	ErrorCount int64  `json:"ErrorCount"`
}

type TaskDb struct {
	Addr     string `json:"Addr"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DbName   string `json:"DbName"`
}
