package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"strconv"
	"strings"
	"time"
)

/**
{
        "number": "0xf3f088",
        "hash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
"transactions":[
   {
                "blockHash": "0xb49d607f5b80890531e3e1d57798a7573cf8e18048ec0df34e3c81d48115078f",
                "blockNumber": "0xf3f088",
                "hash": "0x66db3e5d338a7c4e62a5e93c91cdadfa3137f192a492397abdc986eb3b42bdb9",
                "accessList": [
                    {
                        "address": "0xe37009f6a756e8997fac8da44da7b5b87731fc22",
                        "storageKeys": [
                            "0xc4ba43fad33cbfd83623af6a18d276f51e6983dfa3b7cbdfca86c9d20686f31a",
                            "0xda08f9c72fdaacc1a51cab773779619b140e9591e39b9b47cb79ad0c3fcbdbc1"
                        ]
                    },
                    {
                        "address": "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2",
                        "storageKeys": [
                            "0x74ca828535d35a5b1358536ee1cd8c04d06a7c9414e4288c637fc0440236a8c4",
                            "0xe82d4ace14a34d99367f5c69eb6b54767615f9cc44834d5cf79494c9672644a9"
                        ]
                    },
                    {
                        "address": "0x6fa870fc1f1d206f689621990bacf31fde420367",
                        "storageKeys": [
                            "0x0000000000000000000000000000000000000000000000000000000000000006",
                            "0x0000000000000000000000000000000000000000000000000000000000000007",
                            "0x0000000000000000000000000000000000000000000000000000000000000009",
                            "0x000000000000000000000000000000000000000000000000000000000000000a",
                            "0x000000000000000000000000000000000000000000000000000000000000000c",
                            "0x0000000000000000000000000000000000000000000000000000000000000008"
                        ]
                    }
                ],
                "chainId": "0x1",
                "from": "0x4321d28bb600e85400e5658e0530c457e2f03250",
                "gas": "0x2d8b2",
                "gasPrice": "0x2aedb2837",
                "input": "0x066fa870fc1f1d206f689621990bacf31fde420367000000840968310ee5",
                "maxFeePerGas": "0x2aedb2837",
                "maxPriorityFeePerGas": "0x0",
                "nonce": "0x5d1",
                "r": "0xf0589386446e1e4e08987c981483b4ea987e3cca166d61fe900f9f24fcdcf7fe",
                "s": "0x38b1f53807a6da002f1ee62e41403e00f08bc6d25576b7becdea7c99d4ed57aa",
                "to": "0xe9e83e5bf3d200e48059085b780031002200c100",
                "transactionIndex": "0x0",
                "type": "0x2",
                "v": "0x0",
                "value": "0x150f"
            }
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
func GetBlockFromJson(json string) (*Block, []*Tx) {
	block := Block{Id: time.Now().UnixNano()}
	r := gjson.Parse(json)
	number := r.Get("number").String()
	block.BlockNumber = number
	number, err := HexToInt(number)
	if err == nil {
		block.BlockNumber = number
	}
	block.BlockHash = r.Get("hash").String()
	block.TotalDifficulty = r.Get("totalDifficulty").String()
	block.ExtraData = r.Get("extraData").String()
	block.GasLimit = r.Get("gasLimit").String()
	block.GasUsed = r.Get("gasUsed").String()
	block.Coinbase = r.Get("miner").String()
	block.Nonce = r.Get("nonce").String()
	block.ParentHash = r.Get("parentHash").String()
	block.ReceiptRoot = r.Get("receiptsRoot").String()
	block.BlockSize = r.Get("size").String()
	block.Root = r.Get("stateRoot").String()

	blockTime := r.Get("timestamp").String()
	block.BlockTime = blockTime
	blockTime, err = HexToInt(blockTime)
	if err == nil {
		block.BlockTime = blockTime
	}

	block.TxRoot = r.Get("transactionsRoot").String()
	block.BaseFee = r.Get("baseFeePerGas").String()
	txs := r.Get("transactions").Array()
	txId := make([]string, 0)
	txList := make([]*Tx, 0)
	for _, tx := range txs {
		x := GetTxFromJson(tx.String())
		txList = append(txList, x)
		txId = append(txId, x.TxHash)
	}

	block.Transactions = txId
	return &block, txList
}

/**
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
func GetTxFromJson(json string) *Tx {
	tx := Tx{Id: time.Now().UnixNano()}
	r := gjson.Parse(json)
	tx.BlockHash = r.Get("blockHash").String()
	//tx.BlockNumber = r.Get("blockNumber").String()
	number := r.Get("blockNumber").String()
	tx.BlockNumber = number
	number, err := HexToInt(number)
	if err == nil {
		tx.BlockNumber = number
	}
	tx.TxHash = r.Get("hash").String()
	tx.FromAddr = r.Get("from").String()
	tx.GasPrice = r.Get("gasPrice").String()
	tx.GasLimit = r.Get("gas").String()
	tx.InputData = r.Get("input").String()
	tx.MaxPrice = r.Get("maxFeePerGas").String()
	tx.PriorityFee = r.Get("maxPriorityFeePerGas").String()
	tx.ToAddr = r.Get("to").String()
	tx.Value = r.Get("value").String()
	tx.TransactionIndex = r.Get("transactionIndex").String()
	tx.Type = r.Get("type").String()
	return &tx
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

func GetReceiptFromJson(js string) *Receipt {
	var receipt Receipt
	r := gjson.Parse(js)

	logs := r.Get("logs").String()
	var temp Logs
	_ = json.Unmarshal([]byte(logs), &temp)
	receipt.Logs = &temp
	receipt.TransactionHash = r.Get("transactionHash").String()
	receipt.BlockHash = r.Get("blockHash").String()
	//receipt.BlockNumber=r.Get("blockNumber").String()

	number := r.Get("blockNumber").String()
	receipt.BlockNumber = number
	number, err := HexToInt(number)
	if err == nil {
		receipt.BlockNumber = number
	}
	receipt.ContractAddress = r.Get("contractAddress").String()
	receipt.EffectiveGasPrice = r.Get("effectiveGasPrice").String()
	receipt.CumulativeGasUsed = r.Get("cumulativeGasUsed").String()
	receipt.From = r.Get("from").String()
	receipt.GasUsed = r.Get("gasUsed").String()
	receipt.LogsBloom = r.Get("logsBloom").String()
	receipt.Status = r.Get("status").String()
	receipt.To = r.Get("to").String()
	receipt.TransactionIndex = r.Get("transactionIndex").String()
	receipt.Type = r.Get("type").String()
	receipt.Id = time.Now().UnixNano()
	receipt.CreateTime = time.Now().Format("2006-01-02")
	return &receipt
}

func GetReceiptListFromJson(js string) []*Receipt {
	array := gjson.Parse(js).Array()
	list := make([]*Receipt, 0, len(array))
	for _, r := range array {
		var receipt Receipt
		logs := r.Get("logs").String()
		var temp Logs
		_ = json.Unmarshal([]byte(logs), &temp)
		receipt.Logs = &temp
		receipt.TransactionHash = r.Get("transactionHash").String()
		receipt.BlockHash = r.Get("blockHash").String()
		//receipt.BlockNumber=r.Get("blockNumber").String()

		number := r.Get("blockNumber").String()
		receipt.BlockNumber = number
		number, err := HexToInt(number)
		if err == nil {
			receipt.BlockNumber = number
		}
		receipt.ContractAddress = r.Get("contractAddress").String()
		receipt.EffectiveGasPrice = r.Get("effectiveGasPrice").String()
		receipt.CumulativeGasUsed = r.Get("cumulativeGasUsed").String()
		receipt.From = r.Get("from").String()
		receipt.GasUsed = r.Get("gasUsed").String()
		receipt.LogsBloom = r.Get("logsBloom").String()
		receipt.Status = r.Get("status").String()
		receipt.To = r.Get("to").String()
		receipt.TransactionIndex = r.Get("transactionIndex").String()
		receipt.Type = r.Get("type").String()
		receipt.Id = time.Now().UnixNano()
		receipt.CreateTime = time.Now().Format("2006-01-02")
		list = append(list, &receipt)
	}

	return list
}

func HexToInt(hex string) (string, error) {
	if len(hex) < 1 {
		return hex, errors.New("params is null")
	}
	if !strings.HasPrefix(hex, "0x") {
		return hex, errors.New("input string must be hex string")
	}
	i, err := strconv.ParseInt(hex, 0, 64)
	if err != nil {
		return hex, err
	}
	return fmt.Sprintf("%v", i), nil
}
