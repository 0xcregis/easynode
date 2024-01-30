package chain

import (
	"testing"

	"github.com/segmentio/kafka-go"
)

func TestGetCoreAddress(t *testing.T) {
	addr := GetCoreAddress(200, "0xb6e268b6675846104feef5582d22f40723164d05")
	t.Log(addr)
}

func TestParseTx(t *testing.T) {
	req := []byte(`
{
  "id": 1706598246238798305,
  "hash": "0xabcca0ed01b966327f2f1f239f8a3201b5e81bf1523083b6872c12137b0b900a",
  "txTime": "1706598233",
  "txStatus": "",
  "blockNumber": "115499728",
  "from": "0x545f731e3ce6ab51c7a30ca08bf0f1a953e30826",
  "to": "0x65cb3ff66cdd12ade34cab66a95108a3af034543",
  "value": "0x38d7ea4c68000",
  "fee": "",
  "gasPrice": "0x61cde8a",
  "maxFeePerGas": "",
  "gas": "0x5208",
  "gasUsed": "",
  "baseFeePerGas": "0x5f872b1",
  "maxPriorityFeePerGas": "",
  "input": "0x",
  "blockHash": "0xc7467877690f58fd7403eee2912c4d44324fffb183b245afffae7d16f2eb839a",
  "transactionIndex": "0x3",
  "type": "0x0",
  "receipt": "{\"id\":1706598246425046657,\"blockHash\":\"0xc7467877690f58fd7403eee2912c4d44324fffb183b245afffae7d16f2eb839a\",\"logsBloom\":\"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000\",\"contractAddress\":\"\",\"transactionIndex\":\"0x3\",\"type\":\"0x0\",\"transactionHash\":\"0xabcca0ed01b966327f2f1f239f8a3201b5e81bf1523083b6872c12137b0b900a\",\"gasUsed\":\"0x5208\",\"l1Fee\":\"0x2b7056aa4da2\",\"blockNumber\":\"115499728\",\"cumulativeGasUsed\":\"0x37a16\",\"from\":\"0x545f731e3ce6ab51c7a30ca08bf0f1a953e30826\",\"to\":\"0x65cb3ff66cdd12ade34cab66a95108a3af034543\",\"effectiveGasPrice\":\"0x61cde8a\",\"logs\":[],\"createTime\":\"2024-01-30\",\"status\":\"0x1\"}"
}
`)
	msg := &kafka.Message{}
	msg.Value = req
	sub, err := ParseTx(64, msg)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sub)
	}

}
