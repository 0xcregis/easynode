package service

import "time"

const (
	NodeTaskTable = "node_task_%v"
	DayFormat     = "20060102"
	TimeFormat    = "2006-01-02 15:04:05"
)

type NodeSource struct {
	Id          int64     `json:"id" gorm:"column:id"`
	BlockChain  int64     `json:"blockChain" gorm:"column:block_chain"` // 公链code
	TxHash      string    `json:"txHash" gorm:"column:tx_hash"`
	BlockHash   string    `json:"blockHash" gorm:"column:block_hash"`
	BlockNumber string    `json:"blockNumber" gorm:"column:block_number"`
	SourceType  int       `json:"sourceType" gorm:"column:source_type"` // 任务类型 1: 交易 2:区块 3.收据
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time"`
}

type NodeTask struct {
	NodeId      string    `json:"node_id" gorm:"column:node_id"`           // 节点的唯一标识
	BlockNumber string    `json:"block_number" gorm:"column:block_number"` // 区块高度
	BlockHash   string    `json:"block_hash" gorm:"column:block_hash"`     // 区块hash
	TxHash      string    `json:"tx_hash" gorm:"column:tx_hash"`           // 交易hash
	TaskType    int8      `json:"task_type" gorm:"column:task_type"`       //  0:保留 1:同步Tx 2:同步Block 3:同步收据
	BlockChain  int64     `json:"block_chain" gorm:"column:block_chain"`   // 公链code, 默认：100 (etc)
	TaskStatus  int       `json:"task_status" gorm:"column:task_status"`   // 0: 初始 1: 成功. 2: 失败.  3: 执行中 其他：重试次数
	CreateTime  time.Time `json:"create_time" gorm:"column:create_time"`
	LogTime     time.Time `json:"log_time" gorm:"column:log_time"`
	Id          int64     `json:"id" gorm:"column:id"`
}

type Tx struct {
	Id          int64  `json:"id" gorm:"column:id"`
	TxHash      string `json:"hash" gorm:"column:hash"`
	TxTime      string `json:"txTime" gorm:"column:tx_time"`
	TxStatus    string `json:"txStatus" gorm:"column:tx_status"`
	BlockNumber string `json:"blockNumber" gorm:"column:block_number"`
	FromAddr    string `json:"from" gorm:"column:from_addr"`
	ToAddr      string `json:"to" gorm:"column:to_addr"`
	Value       string `json:"value" gorm:"column:value"`
	Fee         string `json:"fee" gorm:"column:fee"`
	GasPrice    string `json:"gasPrice" gorm:"column:gas_price"`
	MaxPrice    string `json:"maxFeePerGas" gorm:"column:max_fee_per_gas"`
	GasLimit    string `json:"gas" gorm:"column:gas"`
	GasUsed     string `json:"gasUsed" gorm:"column:gas_used"`
	BaseFee     string `json:"baseFeePerGas" gorm:"column:base_fee_per_gas"`
	PriorityFee string `json:"maxPriorityFeePerGas" gorm:"column:max_priority_fee_per_gas"`
	InputData   string `json:"input" gorm:"column:input_data"`
	BlockHash   string `json:"blockHash" gorm:"column:block_hash"`
	//"transactionIndex": "0x9a",
	//"type": "0x2",

	TransactionIndex string `json:"transactionIndex" gorm:"column:transaction_index"`
	Type             string `json:"type" gorm:"column:tx_type"`
}
