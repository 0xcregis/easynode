package service

type MonitorAddress struct {
	Id          int64  `json:"id" gorm:"column:id"`
	Token       string `json:"token" gorm:"column:token"`
	Address     string `json:"address" gorm:"column:address"`
	BlockChain  int64  `json:"blockChain" gorm:"column:block_chain"`
	TxType      string `json:"txType" gorm:"column:tx_type"`
	AddressType string `json:"addressType" gorm:"column:address_type"` //1:外部账户 2:合约账户
	Decimals    string `json:"decimals" gorm:"decimals"`
	Symbol      string `json:"symbol" gorm:"symbol"`
	TokenName   string `json:"tokenName" gorm:"token_name"`
}

type NodeToken struct {
	Id    int64  `json:"id" gorm:"column:id"`
	Token string `json:"token" gorm:"column:token"`
	Email string `json:"email" gorm:"column:email"`
}

/**
CREATE TABLE IF NOT EXISTS ether.tx
(
id UInt64,--时间戳
tx_hash String,
tx_time String,
tx_status String,
block_number UInt64,
from_addr String,
to_addr String,
value String,
fee String,
gas_price String,
max_price String,
gas_limit String,
gas_used String,
base_fee String,
priority_fee String,
input_date String, --extra data
block_hash String
)

{
        "blockHash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
        "blockNumber": "0xf3f088",
        "hash": "0x5917da4788cdc1383215541744beb93fd804c1902e221d2c5555ce99d9bfff42",
        "accessList": [],
        "chainId": "0x1",
        "from": "0xf4e07370db628044ee8556d1dedb0417bd518970",
        "gas": "0x186a0",
        "gasPrice": "0x2ea75f237",
        "input": "0x095ea7b3000000000000000000000000a152f8bb749c55e9943a3a0a3111d18ee2b3f94effffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
        "maxFeePerGas": "0x45ecedb30",
        "maxPriorityFeePerGas": "0x3b9aca00",
        "nonce": "0x0",
        "r": "0x5db86fbd5caf3b3a0896762e90cf99502671b856d8912c059facd2cf9fb1504b",
        "s": "0x43f6c37ea28057e3768fdc1504dee4e362237c26a0582aad57fe94e504372248",
        "to": "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
        "transactionIndex": "0x9a",
        "type": "0x2",
        "v": "0x0",
        "value": "0x0"
    }

*/

type Tx struct {
	Id          uint64 `json:"id" gorm:"column:id"`
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

/**
{
        "transactionHash": "0x5917da4788cdc1383215541744beb93fd804c1902e221d2c5555ce99d9bfff42",
        "blockHash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
        "blockNumber": "0xf3f088",
        "logs": [
            {
                "transactionHash": "0x5917da4788cdc1383215541744beb93fd804c1902e221d2c5555ce99d9bfff42",
                "address": "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
                "blockHash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
                "blockNumber": "0xf3f088",
                "data": "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
                "logIndex": "0x147",
                "removed": false,
                "topics": [
                    "0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925",
                    "0x000000000000000000000000f4e07370db628044ee8556d1dedb0417bd518970",
                    "0x000000000000000000000000a152f8bb749c55e9943a3a0a3111d18ee2b3f94e"
                ],
                "transactionIndex": "0x9a"
            }
        ],
        "contractAddress": null,
        "effectiveGasPrice": "0x2ea75f237",
        "cumulativeGasUsed": "0xc2ec5d",
        "from": "0xf4e07370db628044ee8556d1dedb0417bd518970",
        "gasUsed": "0xb5d7",
        "logsBloom": "0x00000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000002000000000100000400000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000002090000000000000200000000000000000000000000000000000000000000000",
        "status": "0x1",
        "to": "0x95ad61b0a150d79219dcf64e1e6cc01f0b64c4ce",
        "transactionIndex": "0x9a",
        "type": "0x2"
    }
*/

type Receipt struct {
	Id                uint64 `json:"id"`
	BlockHash         string `json:"blockHash" gorm:"column:block_hash"`
	LogsBloom         string `json:"logsBloom" gorm:"column:logs_bloom"`
	ContractAddress   string `json:"contractAddress" gorm:"column:contract_address"`
	TransactionIndex  string `json:"transactionIndex" gorm:"column:transaction_index"`
	Type              string `json:"type" gorm:"column:tx_type"`
	TransactionHash   string `json:"transactionHash" gorm:"column:transaction_hash"`
	GasUsed           string `json:"gasUsed" gorm:"column:gas_used"`
	BlockNumber       string `json:"blockNumber" gorm:"column:block_number"`
	CumulativeGasUsed string `json:"cumulativeGasUsed" gorm:"column:cumulative_gas_used"`
	From              string `json:"from" gorm:"column:from_addr"`
	To                string `json:"to" gorm:"column:to_addr"`
	EffectiveGasPrice string `json:"effectiveGasPrice" gorm:"column:effective_gas_price"`
	Logs              *Logs  `json:"logs" gorm:"-"`
	LogList           string `json:"-" gorm:"column:logs"`
	CreateTime        string `json:"createTime" gorm:"column:create_time"` // 2006-01-02
	Status            string `json:"status" gorm:"column:status"`
}

type Logs []struct {
	BlockHash        string   `json:"blockHash" gorm:"column:block_hash"`
	Address          string   `json:"address" gorm:"column:address"`
	LogIndex         string   `json:"logIndex" gorm:"column:log_index"`
	Data             string   `json:"data" gorm:"column:data"`
	Removed          bool     `json:"removed" gorm:"column:removed"`
	Topics           []string `json:"topics" gorm:"column:topics"`
	BlockNumber      string   `json:"blockNumber" gorm:"column:block_number"`
	TransactionIndex string   `json:"transactionIndex" gorm:"column:transaction_index"`
	TransactionHash  string   `json:"transactionHash" gorm:"column:transaction_hash"`
}

/**
CREATE TABLE IF NOT EXISTS ether.block
(id UInt64,--时间戳
block_hash String,
block_time String,
block_status String,
block_number UInt64,
parent_hash String,
block_reward String,
fee_recipient String,
total_difficulty String,
block_size UInt32,
gas_limit String,
gas_used String,
base_fee String,
input_date String,
root String,
tx_hash String,
receipt_hash String,
boinbase String,
nonce String
)

{
        "number": "0xf3f088",
        "hash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
        "transactions": [
            "0x66db3e5d338a7c4e62a5e93c91cdadfa3137f192a492397abdc986eb3b42bdb9",
            "0x16682ad6180ad2c65e1135c96462833249f41eba1866e61193d8a3880f129948"
        ],
        "difficulty": "0x0",
        "extraData": "0x627920406275696c64657230783639",
        "gasLimit": "0x1c9c380",
        "gasUsed": "0xc3fd9d",
        "logsBloom": "0xcff73416c8861815909c01c9b99bf32c5952375919140225d08114d4d912895576d64b2985ce445860d3b270028d07ced68209624e2829a1a2f312aa1a7ff40ef886a0dd69888f6ce913a43eeb5632228e80092798604742004b3cc1cd443d01929491a50efa2604ca76ff880ae4c9624a22d960415a6d5490dad592a08e40000a06931052271016954e0e2045121cce9450d091c3cc04582ea924dd2018099d8fd94577a9096ef23a2341d0dc9496100154860229c4afeb795d94ec0291acf5031ef40b26092958559850042436c28b5c3b364a5d1e4a1d7280b923547c360030b6ab9d40a03c68a0848f6088050e0a3456604c00c80e5ddd9eb8d0920bda86",
        "miner": "0x690b9a9e9aa1c9db991c7721a92d351db4fac990",
        "mixHash": "0x3d7e8678c2c523c20e27e2a024a277a4b4bdcd6315f2d51a38bfa5cf6a1dc654",
        "nonce": "0x0000000000000000",
        "parentHash": "0x73b66f8d4a99efe315b3a2bb82787c93693ef3e7f23a232c245f402f69576bcf",
        "receiptsRoot": "0x59adae8e952f90e9653e71a866a90b051b9307636698f6aa4eeb07ba820709d0",
        "sha3Uncles": "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
        "size": "0xe610",
        "stateRoot": "0x46d87b856a87ec0f9a4dd5dd1fbb4ca00d7210b2f58c63a5187d2ebdb3d0f478",
        "timestamp": "0x63759f6b",
        "totalDifficulty": "0xc70d815d562d3cfa955",
        "transactionsRoot": "0xf25865fe021fa09e3d5dbbcc215430e3556a72f4f577c982e968570171f26e4a",
        "uncles": [],
        "baseFeePerGas": "0x2aedb2837"
    }
*/

type Block struct {
	Id              uint64   `json:"id" gorm:"column:id"`
	BlockHash       string   `json:"hash" gorm:"column:hash"`
	BlockTime       string   `json:"timestamp" gorm:"column:block_time"`
	BlockStatus     string   `json:"blockStatus" gorm:"column:block_status"`
	BlockNumber     string   `json:"number" gorm:"column:block_number"`
	ParentHash      string   `json:"parentHash" gorm:"column:parent_hash"`
	BlockReward     string   `json:"blockReward" gorm:"column:block_reward"`
	FeeRecipient    string   `json:"feeRecipient" gorm:"column:fee_recipient"`
	TotalDifficulty string   `json:"totalDifficulty" gorm:"column:total_difficulty"`
	BlockSize       string   `json:"size" gorm:"column:block_size"`
	GasLimit        string   `json:"gasLimit" gorm:"column:gas_limit"`
	GasUsed         string   `json:"gasUsed" gorm:"column:gas_used"`
	BaseFee         string   `json:"baseFeePerGas" gorm:"column:base_fee_per_gas"`
	ExtraData       string   `json:"extraData" gorm:"column:extra_data"`
	Root            string   `json:"stateRoot" gorm:"column:state_root"`
	Transactions    []string `json:"transactions" gorm:"-"`
	TransactionList string   `json:"-" gorm:"column:transactions"`
	TxRoot          string   `json:"transactionsRoot" gorm:"column:transactions_root"`
	ReceiptRoot     string   `json:"receiptsRoot" gorm:"column:receipts_root"`
	Coinbase        string   `json:"miner" gorm:"column:miner"`
	Nonce           string   `json:"nonce" gorm:"column:nonce"`
}

type WsReqMessage struct {
	Id         int64
	Code       int64   //1:订阅资产转移交易，2:取消订阅资产转移交易
	BlockChain []int64 `json:"blockChain"`
	Params     map[string]string
}

type WsRespMessage struct {
	Id         int64
	Code       int64
	BlockChain []int64 `json:"blockChain"`
	Status     int     //0:成功 1：失败
	Err        string
	Params     map[string]string
	Resp       interface{}
}

type WsPushMessage struct {
	Code       int64 //1:tx,2:block,3:receipt
	BlockChain int64 `json:"blockChain"`
	Data       interface{}
}
