package taskapi

import "time"

type NodeTask struct {
	Id          int64     `json:"id"  gorm:"column:id"`
	NodeId      string    `json:"nodeId" gorm:"column:node_id"`
	BlockNumber string    `json:"blockNumber" gorm:"column:block_number"`
	BlockHash   string    `json:"blockHash" gorm:"column:block_hash"`
	TxHash      string    `json:"txHash" gorm:"column:tx_hash"`
	TaskType    int       `json:"taskType" gorm:"column:task_type"` // 0:保留 1:同步Tx. 2:同步Block 3:同步Receipt 4:区块Tx 5:区块receipt
	BlockChain  int64     `json:"blockChain" gorm:"column:block_chain"`
	TaskStatus  int       `json:"taskStatus" gorm:"column:task_status"` //0: 初始 1: 成功. 2: 失败.  3: 执行中 其他：重试次数
	CreateTime  time.Time `json:"createTime" gorm:"column:create_time"`
	LogTime     time.Time `json:"logTime" gorm:"column:log_time"`
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
