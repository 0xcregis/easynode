package config

type Config struct {
	RootPath   string                   `json:"RootPath"`
	Port       int                      `json:"Port"`
	BlockChain []int64                  `json:"BlockChain"`
	Cluster    map[int64][]*NodeCluster `json:"Cluster"`
}

type NodeCluster struct {
	NodeUrl    string `json:"NodeUrl"`
	NodeToken  string `json:"NodeToken"`
	Weight     int64  `json:"Weight"`
	ErrorCount int64  `json:"ErrorCount"`
}
