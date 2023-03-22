package config

/**
{
    "TaskDb":{
        "Addr":"192.168.2.20",
        "Port":5432,
        "User":"root",
        "Password":"1234567789",
        "DbName":"easy_node",
        "Table":"node.node_task"
    },
    "NodeInfoDb":{
        "Addr":"192.168.2.20",
        "Port":5432,
        "User":"root",
        "Password":"1234567789",
        "DbName":"easy_node",
        "Table":"node.node_info"
    },
    "Chains":[
        {
            "BlockChain_Name":"ether",
            "BlockChain_Code":100,
            "Block_Task":{
                "From_Cluster":[
                    {
                        "Host":"192.168.2.20",
                        "Port":9092
                    }
                ],
                "Kafka":{
                    "Host":"192.168.2.20",
                    "Port":9092,
                    "Topic":"test",
                    "Partition":0
                }
            },
            "Tx_Task":{
                "From_Cluster":[
                    {
                        "Host":"192.168.2.20",
                        "Port":9092
                    }
                ],
                "Kafka":{
                    "Host":"192.168.2.20",
                    "Port":9092,
                    "Topic":"test2",
                    "Partition":0
                }
            }
        }
    ]
}
*/

type TaskDb struct {
	Addr     string `json:"Addr"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DbName   string `json:"DbName"`
	Table    string `json:"Table"`
}

type SourceDb struct {
	Addr     string `json:"Addr"`
	Port     int    `json:"Port"`
	User     string `json:"User"`
	Password string `json:"Password"`
	DbName   string `json:"DbName"`
	Table    string `json:"Table"`
}

type NodeInfoDb struct {
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
	Host      string `json:"Host"`
	Port      int    `json:"Port"`
	Topic     string `json:"Topic"`
	Partition int    `json:"Partition"`
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

type Etcd struct {
	Host string `json:"Host"`
	Key  string `json:"Key"`
	Ttl  int64  `json:"Ttl"`
}

type Chain struct {
	Etcd           *Etcd        `json:"Etcd"`
	BlockChainName string       `json:"BlockChainName"`
	BlockChainCode int          `json:"BlockChainCode"`
	NodeWeight     int          `json:"NodeWeight"`
	PullReceipt    bool         `json:"PullReceipt"`
	PullTx         bool         `json:"PullTx"`
	Kafka          *Kafka       `json:"Kafka"`
	BlockTask      *BlockTask   `json:"BlockTask"`
	TxTask         *TxTask      `json:"TxTask"`
	ReceiptTask    *ReceiptTask `json:"ReceiptTask"`
}

type LogConfig struct {
	Path  string `json:"Path"`
	Delay int64  `json:"Delay"`
}

type Config struct {
	TaskDb     *TaskDb     `json:"NodeTaskDb"`
	SourceDb   *SourceDb   `json:"NodeSourceDb"`
	NodeInfoDb *NodeInfoDb `json:"NodeInfoDb"`
	Chains     []*Chain    `json:"Chains"`
	LogConfig  *LogConfig  `json:"Log"`
}
