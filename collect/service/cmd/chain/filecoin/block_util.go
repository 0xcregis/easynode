package filecoin

import (
	"encoding/json"
	"time"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/common/util"
	"github.com/tidwall/gjson"
)

//GetReceiptFromJson
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
func GetReceiptFromJson(js string) *collect.Receipt {
	var receipt collect.Receipt
	r := gjson.Parse(js)

	logs := r.Get("logs").String()
	var temp []*collect.Logs
	_ = json.Unmarshal([]byte(logs), &temp)
	receipt.Logs = temp
	receipt.TransactionHash = r.Get("transactionHash").String()
	receipt.BlockHash = r.Get("blockHash").String()
	//receipt.BlockNumber=r.Get("blockNumber").String()

	number := r.Get("blockNumber").String()
	receipt.BlockNumber = number
	number, err := util.HexToInt(number)
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
