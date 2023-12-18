package tron

import (
	"testing"

	"github.com/0xcregis/easynode/store"
)

func TestParseTx(t *testing.T) {
	req := []byte(`
{
  "blockId": "00000000025fc02f22f068a6e04913196756836fb46c73bd80b9340a7ca37ced",
  "receipt": "{\"id\":\"e7757a0b0d2df629bf363f39062934c977b3a204f82f8a882f019cadeddf6298\",\"fee\":267000,\"blockNumber\":39829551,\"blockTimeStamp\":1702881432000,\"from\":\"\",\"to\":\"\",\"contractResult\":[\"\"],\"contract_address\":\"\",\"receipt\":{\"energy_fee\":0,\"energy_usage_total\":0,\"net_usage\":0,\"result\":\"\",\"energy_penalty_total\":0},\"log\":null,\"internal_transactions\":null}",
  "tx": "{\"ret\":[{\"contractRet\":\"SUCCESS\"}],\"signature\":[\"65c15050fb80bde6e911e9be9322a99cf5d3fd65a6eb6e7fa56b2139df670bd17d2b8f5bb9f082b6873606e9181e9daaeb82d0fc605268dfeee3c435dc6ab69e1c\"],\"txID\":\"e7757a0b0d2df629bf363f39062934c977b3a204f82f8a882f019cadeddf6298\",\"raw_data\":{\"contract\":[{\"parameter\":{\"value\":{\"amount\":1000000,\"owner_address\":\"418e5ad0c377005da36e592424e84ea0fa317321c2\",\"to_address\":\"412f83df6e9705f82e001e32e7d3f45551fda36d03\"},\"type_url\":\"type.googleapis.com/protocol.TransferContract\"},\"type\":\"TransferContract\"}],\"ref_block_bytes\":\"c02d\",\"ref_block_hash\":\"e732b171fc746141\",\"expiration\":1702881486000,\"timestamp\":1702881426000},\"raw_data_hex\":\"0a02c02d2208e732b171fc74614140b0c995ddc7315a67080112630a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412320a15418e5ad0c377005da36e592424e84ea0fa317321c21215412f83df6e9705f82e001e32e7d3f45551fda36d0318c0843d70d0f491ddc731\"}"
}
`)
	sub, err := ParseTx(req, store.TronTopic, 198)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(sub)
	}
}
