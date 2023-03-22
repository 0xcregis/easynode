package service

import "time"

const (
	NodeInfoTable   = "node_info"
	NodeSourceTable = "node_source"
	NodeTaskTable   = "node_task_%v"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
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
