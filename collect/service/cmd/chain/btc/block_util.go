package btc

import (
	"time"

	"github.com/0xcregis/easynode/collect"
	"github.com/tidwall/gjson"
)

// GetBlockFromJson
/*
*
{
    "result": {
        "hash": "00000000000000000000fbd80dd4a8502f6a06d10c6179c602d5a0ea24f43d38",
        "confirmations": 3782,
        "height": 804901,
        "version": 651051008,
        "versionHex": "26ce4000",
        "merkleroot": "4933fa68c881d9b1896bba468a70d98dbc762760371c2419e7820b23b3305519",
        "time": 1693035700,
        "mediantime": 1693034785,
        "nonce": 1315411571,
        "bits": "17050f7b",
        "difficulty": 55621444139429.57,
        "chainwork": "0000000000000000000000000000000000000000535ffc4a39a980ef017b4cd4",
        "nTx": 4920,
        "previousblockhash": "00000000000000000002e9a67876d80bb07812b24a0e47439ef7d1f9721bd2cd",
        "nextblockhash": "000000000000000000005064676c9b4c843786d8e7d74a350fd09e1f80a94bf4",
        "strippedsize": 733924,
        "size": 1796085,
        "weight": 3997857,
        "tx": [
            "3db988315cebdaba11db5100c21473420a44047f67140f19cd58084554e7724b",
            "af0291ef619d4a50562943b4165b0fd3efac153b378070fe10ae70723ccc65ee",
            "85af5e8f6069d78a426ab5152de5d22babc9e684150b65e5c4242bfbf7f365ee"
        ]
    },
    "error": null,
    "id": null
}
*/
func GetBlockFromJson(json string) (*collect.Block, []*collect.Tx) {
	block := collect.Block{Id: time.Now().UnixNano()}
	r := gjson.Parse(json)
	number := r.Get("height").String()
	block.BlockNumber = number
	block.BlockHash = r.Get("hash").String()
	block.TotalDifficulty = r.Get("difficulty").String()
	block.Nonce = r.Get("nonce").String()
	block.ParentHash = r.Get("previousblockhash").String()
	block.BlockSize = r.Get("size").String()
	block.Root = r.Get("merkleroot").String()
	block.BlockTime = r.Get("time").String()
	txs := r.Get("tx").Array()
	txId := make([]string, 0)
	txList := make([]*collect.Tx, 0)
	for _, tx := range txs {
		if tx.IsObject() {
			x := GetTxFromJson(tx.String())
			x.BlockNumber = number
			txList = append(txList, x)
			txId = append(txId, x.TxHash)
		} else {
			txId = append(txId, tx.String())
		}
	}

	block.Transactions = txId
	return &block, txList
}

// GetTxFromJson
/**
{
    "result": {
        "txid": "10b54fd708ab2e5703979b4ba27ca0339882abc2062e77fbe51e625203a49642",
        "hash": "b0027a73bc8705c3ad1328ca3caf45046924ab036c3aa96348c2cd6a3fb046ce",
        "version": 1,
        "size": 246,
        "vsize": 165,
        "weight": 657,
        "locktime": 0,
        "vin": [
            {
                "txid": "a177ff1980c4b121c5ee2c5971ce92fa3c69c03a6c339dbda21a051dd2037ee1",
                "vout": 1,
                "scriptSig": {
                    "asm": "00146d76e574b5f4825fe740ba6c41aaf1b319dfb80c",
                    "hex": "1600146d76e574b5f4825fe740ba6c41aaf1b319dfb80c"
                },
                "txinwitness": [
                    "304402206701306a4750908fd48dead54331a3c7b4dce04ec10bfc6dd32049e2cff061a5022013c9d66827fabbeaadeb30b41c09aca2daddf4628cd00e3b993b1c86a12ff51901",
                    "034bcb9be1daf6ce1193774d15f863768b621bc95a363f1da5810129e961a23174"
                ],
                "prevout": {
                    "generated": false,
                    "height": 676231,
                    "value": 1.86504802,
                    "scriptPubKey": {
                        "asm": "OP_HASH160 797922d6bb8a1a2e87592871c2f88267a8ad29fe OP_EQUAL",
                        "desc": "addr(3CmJoSNYp993oPHT3UwNYPSEwFfrHaymF4)#w4fe8k9u",
                        "hex": "a914797922d6bb8a1a2e87592871c2f88267a8ad29fe87",
                        "address": "3CmJoSNYp993oPHT3UwNYPSEwFfrHaymF4",
                        "type": "scripthash"
                    }
                },
                "sequence": 4294967295
            }
        ],
        "vout": [
            {
                "value": 0.00105089,
                "n": 0,
                "scriptPubKey": {
                    "asm": "0 422002d927a1cae901eac668444cce8dd0ae60d5",
                    "desc": "addr(bc1qggsq9kf8589wjq02ce5ygnxw3hg2ucx4c9vet7)#zdanl330",
                    "hex": "0014422002d927a1cae901eac668444cce8dd0ae60d5",
                    "address": "bc1qggsq9kf8589wjq02ce5ygnxw3hg2ucx4c9vet7",
                    "type": "witness_v0_keyhash"
                }
            },
            {
                "value": 1.86364713,
                "n": 1,
                "scriptPubKey": {
                    "asm": "OP_HASH160 f5b48d1130dc3d366d1eabf6783a552d1c8e08f4 OP_EQUAL",
                    "desc": "addr(3Q6BmSyUkvqFLdT68RH3ohVY5MWNQP3Qwt)#xg5wpetj",
                    "hex": "a914f5b48d1130dc3d366d1eabf6783a552d1c8e08f487",
                    "address": "3Q6BmSyUkvqFLdT68RH3ohVY5MWNQP3Qwt",
                    "type": "scripthash"
                }
            }
        ],
        "fee": 0.00035,
        "hex": "01000000000101e17e03d21d051aa2bd9d336c3ac0693cfa92ce71592ceec521b1c48019ff77a101000000171600146d76e574b5f4825fe740ba6c41aaf1b319dfb80cffffffff02819a010000000000160014422002d927a1cae901eac668444cce8dd0ae60d529b31b0b0000000017a914f5b48d1130dc3d366d1eabf6783a552d1c8e08f4870247304402206701306a4750908fd48dead54331a3c7b4dce04ec10bfc6dd32049e2cff061a5022013c9d66827fabbeaadeb30b41c09aca2daddf4628cd00e3b993b1c86a12ff5190121034bcb9be1daf6ce1193774d15f863768b621bc95a363f1da5810129e961a2317400000000",
        "blockhash": "0000000000000000000c70f2c7397a67c9e79c8a57b91f715df65f6738d4a5d8",
        "confirmations": 132446,
        "time": 1616673563,
        "blocktime": 1616673563
    },
    "error": null,
    "id": null
}
*/
func GetTxFromJson(json string) *collect.Tx {
	tx := collect.Tx{Id: time.Now().UnixNano()}
	r := gjson.Parse(json)
	tx.BlockHash = r.Get("blockhash").String()
	tx.TxHash = r.Get("txid").String()
	tx.FromAddr = r.Get("vin").String()
	//tx.GasPrice = r.Get("gasPrice").String()
	//tx.GasLimit = r.Get("gas").String()
	tx.Fee = r.Get("fee").String()
	tx.InputData = r.Get("hex").String()
	//tx.MaxPrice = r.Get("maxFeePerGas").String()
	//tx.PriorityFee = r.Get("maxPriorityFeePerGas").String()
	tx.ToAddr = r.Get("vout").String()
	//tx.Value = r.Get("value").String()
	tx.TxTime = r.Get("blocktime").String()
	//tx.TransactionIndex = r.Get("transactionIndex").String()
	tx.Type = r.Get("version").String()
	return &tx
}
