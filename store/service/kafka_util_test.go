package service

import (
	"log"
	"testing"
)

func TestParseTxForTron(t *testing.T) {

	str := `
{
  "blockId": "00000000022949509ea3e17d4e946d27ceb4f4a1b73c2e4d9886f98a3455304a",
  "receipt": "{\"id\":\"60bce1ce829e0ca55806324356601e52f28c562a974f69c4e884863d87ec9418\",\"fee\":1100000,\"blockNumber\":36260176,\"blockTimeStamp\":1692064845000,\"from\":\"\",\"to\":\"\",\"contractResult\":[\"\"],\"contract_address\":\"\",\"receipt\":{\"energy_fee\":0,\"energy_usage_total\":0,\"net_usage\":0,\"result\":\"\",\"energy_penalty_total\":0},\"log\":null}",
  "tx": "{\"ret\":[{\"contractRet\":\"SUCCESS\"}],\"signature\":[\"ffaefd848f031ded00e6ab4429ce94564acec3918e47217ac6f81467433b090f1ef35f57eb0d945fe4c0c357172b379a3b3406666660dc7797767630d512507600\"],\"txID\":\"60bce1ce829e0ca55806324356601e52f28c562a974f69c4e884863d87ec9418\",\"raw_data\":{\"contract\":[{\"parameter\":{\"value\":{\"amount\":1000000000,\"owner_address\":\"41ee7115dbfdde5a70178e85b6439010c4c48bafe8\",\"to_address\":\"41caaa4e9f2dbe9b1c0d288ed745db482bbd304c1e\"},\"type_url\":\"type.googleapis.com/protocol.TransferContract\"},\"type\":\"TransferContract\"}],\"ref_block_bytes\":\"494d\",\"ref_block_hash\":\"23bdfbb5d6c57e38\",\"expiration\":1692065138836,\"timestamp\":1692064838836},\"raw_data_hex\":\"0a02494d220823bdfbb5d6c57e38409491c4b79f315a69080112650a2d747970652e676f6f676c65617069732e636f6d2f70726f746f636f6c2e5472616e73666572436f6e747261637412340a1541ee7115dbfdde5a70178e85b6439010c4c48bafe8121541caaa4e9f2dbe9b1c0d288ed745db482bbd304c1e188094ebdc0370b4e9b1b79f31\"}"
}
`
	log.Println(ParseTxForTron([]byte(str)))
}

func TestDiv(t *testing.T) {
	t.Log(Div("1000", 2))
}
