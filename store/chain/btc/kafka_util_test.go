package btc

import (
	"testing"
)

func TestGetBlockFromKafka(t *testing.T) {
	str := `
{
  "id": 1695626965410917000,
  "hash": "0000000000000000000053e69a37c5e2b7dcb16cd3c2f34dec814feb47ffecc7",
  "timestamp": "1695623822",
  "blockStatus": "",
  "number": "809250",
  "parentHash": "00000000000000000001edbda0b7b8d5cc8e83fc44a49f15e972aa86534f6c61",
  "blockReward": "",
  "feeRecipient": "",
  "totalDifficulty": "57119871304635.31",
  "size": "1705538",
  "gasLimit": "",
  "gasUsed": "",
  "baseFeePerGas": "",
  "extraData": "",
  "stateRoot": "52282bb24cb9f4494f559c73eef343d69551c9429dbd3093eb08cc30b29fa543",
  "transactions": [
    "d27ade29d9f9f0ea91d3c81a3ef691e10829828bc89870257003d6f610589ad1",
    "6c8fa88f8744642a1a3e563ad0a7ddbd6920e73f0eb9a61c34c2248847d0241e",
    "57d0239212908c599fb3e398b94773a6ed4c86d4231f3b9ca17bc9d22655bf5a",
    "351049dcc2fc9c50b21925221465870ed85ed49b6429073f0e97afc47b24e1ea",
    "8ea5ee165f067285e3312d0ec2f1749a4031a1bc35a6bac1e3dd7d9a8c3d230b",
    "7a10544c45a0414871de33c7eaea9c0ed8f1645c9e937fc61eefaf385720c91d",
    "6ef85392fbcbc03d9eefdddc9930a692180fdd8a03a43f532fa69c2ee10db525",
    "78dacffafd6c4ba4fe5837742d14e7997a27a0eb6ce6e4f0e424d965459c0f3f",
    "6a91aa49c69d111ec4f7fc6bfb03ab78ae4185d639d05a6dd62daf73c93c9247",
    "72fc1ac42887f405aad3cacba938f13b2d71bd29e2a5f348bb3ac8bda688b463",
    "9e253ff7322ca64ec69605d9fb16e65e4e31e37b48e873572090d976c682c686",
    "d27e45e8a257aca05acc1264cbb570fc9955f94c3402aa00a55c5a8a41971fbc",
    "faa78b15e2e7c43c340b62bca10735bf844f1193cce00381c24b2c3b6b02ff58",
    "0f33eb0a7d17811d1050de9089e1227230dff6ec82e4a20a2a71fe936f51e237",
    "38200feed09ce70a6f8997545a3b3e8e110330cfbb940a329f70c22ce956b9ca",
    "7333b6490f094fe7aa449cd98829f44ae0582e12454fb27e6e9d3bdbfadff223",
    "5ae007ed969c2956a9e80df7ab2932878fa18165ad25e2ea9b56c408675c836a",
    "585c9dfde2d37acc15c39a51a09762d836e56b5740c3f3919c404febaf40b792",
    "963e70084cf3228ced7cf2dc17f6f3556c831e6c52ba0866a3c2432ad44ee0ec",
    "9b3c6ff4abe36d64dc14d5f2812185d564fa94cb6b05fe8d957e141268f9eacf",
    "6fd3d4d5d45ab43fe8ec255c95e2603fcb7a6aa3468d7a525b2a4fb192bc9d41",
    "c37d86facc7f5431769577a6a5f8aa86f1cf82d7c05f3b2f3d43737093a4cd90",
    "c404606803bf1cf9f8c38544f6256668e0b2d34a12cba82e78ab7dd96a783fc5",
    "86ef5741cb9a7303f286f5715313d7ae477ab720361cbc1f5083f6520a96b5ed",
    "bd545e2e2d61458bc64d028f016cb4e2309d0e7ab6215a9ff3d98cf46e5be450",
    "3dad8953baa42cacb9f891510f13a63e7cc2b44979c8814db7a25ccc66e26756",
    "efa79c5c7eb3ed2068cfb0cfbb3db2e7874fec8343d0169a7c245613f00ca0a2",
    "e1ddd13b657e414c00506e37af9aff9bf2fa7c86c269717710ab297c41ffb34f",
    "c91c44db9b3679e2e708749e03149a657ef865337b95f0e6f51f33d35c0dcbf1"
  ],
  "transactionsRoot": "",
  "receiptsRoot": "",
  "miner": "",
  "nonce": "1963354372"
}   

`

	b, err := GetBlockFromKafka([]byte(str))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(b)
	}

}

func TestGetTxFromKafka(t *testing.T) {

	str := `
{
  "id": 1695625620675744000,
  "hash": "9f6dc48de9b83db9e4f847802598c78210e39d76361ac6eccd3361972bc1d037",
  "txTime": "1695625259",
  "txStatus": "",
  "blockNumber": "",
  "from": "[{\"txid\":\"8aefcf0f75a9a1e632c4b1c17bb510e3974e9adf545d72df1e505fcfa89952a2\",\"vout\":0,\"scriptSig\":{\"asm\":\"0014d3502375b1d8642d06ad6df64f4ed9e0fee394ab\",\"hex\":\"160014d3502375b1d8642d06ad6df64f4ed9e0fee394ab\"},\"txinwitness\":[\"3045022100e4f1feac7adf3361d9ee01604511f7c47e6283568cb1fe62d32a060c0f6e5b56022035d40b2ddff4b170c23de2897704612cf161d9913fd9260914170daccafc8dad01\",\"038b65670f08b885fbce0994eddeb4ca22b167f5aca18f1846c615e4d0fdd0472a\"],\"prevout\":{\"generated\":false,\"height\":621762,\"value\":0.00278793,\"scriptPubKey\":{\"asm\":\"OP_HASH160 e3b1b6bbbd14dcb39bdcb4067145032e6e9fefa0 OP_EQUAL\",\"desc\":\"addr(3NSxCrsfQhr44VLZAyEsYFnpDT3UFEeM4P)#g6kywuwd\",\"hex\":\"a914e3b1b6bbbd14dcb39bdcb4067145032e6e9fefa087\",\"address\":\"3NSxCrsfQhr44VLZAyEsYFnpDT3UFEeM4P\",\"type\":\"scripthash\"}},\"sequence\":4294967295},{\"txid\":\"a28d115bbdf16c4a5ef1797c30c82444701be9953371e67413aaae67529304b9\",\"vout\":2,\"scriptSig\":{\"asm\":\"\",\"hex\":\"\"},\"txinwitness\":[\"3044022041f91789dbdd823d8d9dcdc22ae054b5a338eb3e80c844bc4ddafda5b95fad7a022072215b7e97f8f644de25f8ea636d5b6685b1673a3c13f81bd0f66e18c74b33c501\",\"02598122f3f1062147c117fde19e67af124351ccd6aadfd50926d7678ecc5526d7\"],\"prevout\":{\"generated\":false,\"height\":785864,\"value\":0.00083080,\"scriptPubKey\":{\"asm\":\"0 a57a9d8e7937518ec2a98ec3f6f58a908747e64a\",\"desc\":\"addr(bc1q54afmrnexagcas4f3mpldav2jzr50ej2e5jdu6)#4n0nutz7\",\"hex\":\"0014a57a9d8e7937518ec2a98ec3f6f58a908747e64a\",\"address\":\"bc1q54afmrnexagcas4f3mpldav2jzr50ej2e5jdu6\",\"type\":\"witness_v0_keyhash\"}},\"sequence\":4294967295}]",
  "to": "[{\"value\":0.00021631,\"n\":0,\"scriptPubKey\":{\"asm\":\"OP_HASH160 50724e9f9905841fec561348b77fbba6c2952cb2 OP_EQUAL\",\"desc\":\"addr(392Nv65LfPNUiK799oYAPhdxXwHtqZ8Cbh)#whx754m9\",\"hex\":\"a91450724e9f9905841fec561348b77fbba6c2952cb287\",\"address\":\"392Nv65LfPNUiK799oYAPhdxXwHtqZ8Cbh\",\"type\":\"scripthash\"}},{\"value\":0.00336747,\"n\":1,\"scriptPubKey\":{\"asm\":\"OP_HASH160 97b3281b726565abb6c437bf9fb8ea6726558fc8 OP_EQUAL\",\"desc\":\"addr(3FX8aXuhgC2f1L4pPvbfSWfdfq4woFakM4)#swf2z4md\",\"hex\":\"a91497b3281b726565abb6c437bf9fb8ea6726558fc887\",\"address\":\"3FX8aXuhgC2f1L4pPvbfSWfdfq4woFakM4\",\"type\":\"scripthash\"}}]",
  "value": "",
  "fee": "0.00003495",
  "gasPrice": "",
  "maxFeePerGas": "",
  "gas": "",
  "gasUsed": "",
  "baseFeePerGas": "",
  "maxPriorityFeePerGas": "",
  "input": "01000000000102a25299a8cf5f501edf725d54df9a4e97e310b57bc1b1c432e6a1a9750fcfef8a0000000017160014d3502375b1d8642d06ad6df64f4ed9e0fee394abffffffffb904935267aeaa1374e6713395e91b704424c8307c79f15e4a6cf1bd5b118da20200000000ffffffff027f5400000000000017a91450724e9f9905841fec561348b77fbba6c2952cb2876b2305000000000017a91497b3281b726565abb6c437bf9fb8ea6726558fc88702483045022100e4f1feac7adf3361d9ee01604511f7c47e6283568cb1fe62d32a060c0f6e5b56022035d40b2ddff4b170c23de2897704612cf161d9913fd9260914170daccafc8dad0121038b65670f08b885fbce0994eddeb4ca22b167f5aca18f1846c615e4d0fdd0472a02473044022041f91789dbdd823d8d9dcdc22ae054b5a338eb3e80c844bc4ddafda5b95fad7a022072215b7e97f8f644de25f8ea636d5b6685b1673a3c13f81bd0f66e18c74b33c5012102598122f3f1062147c117fde19e67af124351ccd6aadfd50926d7678ecc5526d700000000",
  "blockHash": "0000000000000000000143aefa609bb22efbfb8cbcba4e608293f97723373613",
  "transactionIndex": "",
  "type": "1",
  "receipt": ""
}

`

	tx, err := GetTxFromKafka([]byte(str))
	if err != nil {
		t.Error(err)
	} else {
		t.Log(tx)
	}

}

func TestParseTx(t *testing.T) {

	str := `
{
  "id": 1695625620675744000,
  "hash": "9f6dc48de9b83db9e4f847802598c78210e39d76361ac6eccd3361972bc1d037",
  "txTime": "1695625259",
  "txStatus": "",
  "blockNumber": "",
  "from": "[{\"txid\":\"8aefcf0f75a9a1e632c4b1c17bb510e3974e9adf545d72df1e505fcfa89952a2\",\"vout\":0,\"scriptSig\":{\"asm\":\"0014d3502375b1d8642d06ad6df64f4ed9e0fee394ab\",\"hex\":\"160014d3502375b1d8642d06ad6df64f4ed9e0fee394ab\"},\"txinwitness\":[\"3045022100e4f1feac7adf3361d9ee01604511f7c47e6283568cb1fe62d32a060c0f6e5b56022035d40b2ddff4b170c23de2897704612cf161d9913fd9260914170daccafc8dad01\",\"038b65670f08b885fbce0994eddeb4ca22b167f5aca18f1846c615e4d0fdd0472a\"],\"prevout\":{\"generated\":false,\"height\":621762,\"value\":0.00278793,\"scriptPubKey\":{\"asm\":\"OP_HASH160 e3b1b6bbbd14dcb39bdcb4067145032e6e9fefa0 OP_EQUAL\",\"desc\":\"addr(3NSxCrsfQhr44VLZAyEsYFnpDT3UFEeM4P)#g6kywuwd\",\"hex\":\"a914e3b1b6bbbd14dcb39bdcb4067145032e6e9fefa087\",\"address\":\"3NSxCrsfQhr44VLZAyEsYFnpDT3UFEeM4P\",\"type\":\"scripthash\"}},\"sequence\":4294967295},{\"txid\":\"a28d115bbdf16c4a5ef1797c30c82444701be9953371e67413aaae67529304b9\",\"vout\":2,\"scriptSig\":{\"asm\":\"\",\"hex\":\"\"},\"txinwitness\":[\"3044022041f91789dbdd823d8d9dcdc22ae054b5a338eb3e80c844bc4ddafda5b95fad7a022072215b7e97f8f644de25f8ea636d5b6685b1673a3c13f81bd0f66e18c74b33c501\",\"02598122f3f1062147c117fde19e67af124351ccd6aadfd50926d7678ecc5526d7\"],\"prevout\":{\"generated\":false,\"height\":785864,\"value\":0.00083080,\"scriptPubKey\":{\"asm\":\"0 a57a9d8e7937518ec2a98ec3f6f58a908747e64a\",\"desc\":\"addr(bc1q54afmrnexagcas4f3mpldav2jzr50ej2e5jdu6)#4n0nutz7\",\"hex\":\"0014a57a9d8e7937518ec2a98ec3f6f58a908747e64a\",\"address\":\"bc1q54afmrnexagcas4f3mpldav2jzr50ej2e5jdu6\",\"type\":\"witness_v0_keyhash\"}},\"sequence\":4294967295}]",
  "to": "[{\"value\":0.00021631,\"n\":0,\"scriptPubKey\":{\"asm\":\"OP_HASH160 50724e9f9905841fec561348b77fbba6c2952cb2 OP_EQUAL\",\"desc\":\"addr(392Nv65LfPNUiK799oYAPhdxXwHtqZ8Cbh)#whx754m9\",\"hex\":\"a91450724e9f9905841fec561348b77fbba6c2952cb287\",\"address\":\"392Nv65LfPNUiK799oYAPhdxXwHtqZ8Cbh\",\"type\":\"scripthash\"}},{\"value\":0.00336747,\"n\":1,\"scriptPubKey\":{\"asm\":\"OP_HASH160 97b3281b726565abb6c437bf9fb8ea6726558fc8 OP_EQUAL\",\"desc\":\"addr(3FX8aXuhgC2f1L4pPvbfSWfdfq4woFakM4)#swf2z4md\",\"hex\":\"a91497b3281b726565abb6c437bf9fb8ea6726558fc887\",\"address\":\"3FX8aXuhgC2f1L4pPvbfSWfdfq4woFakM4\",\"type\":\"scripthash\"}}]",
  "value": "",
  "fee": "0.00003495",
  "gasPrice": "",
  "maxFeePerGas": "",
  "gas": "",
  "gasUsed": "",
  "baseFeePerGas": "",
  "maxPriorityFeePerGas": "",
  "input": "01000000000102a25299a8cf5f501edf725d54df9a4e97e310b57bc1b1c432e6a1a9750fcfef8a0000000017160014d3502375b1d8642d06ad6df64f4ed9e0fee394abffffffffb904935267aeaa1374e6713395e91b704424c8307c79f15e4a6cf1bd5b118da20200000000ffffffff027f5400000000000017a91450724e9f9905841fec561348b77fbba6c2952cb2876b2305000000000017a91497b3281b726565abb6c437bf9fb8ea6726558fc88702483045022100e4f1feac7adf3361d9ee01604511f7c47e6283568cb1fe62d32a060c0f6e5b56022035d40b2ddff4b170c23de2897704612cf161d9913fd9260914170daccafc8dad0121038b65670f08b885fbce0994eddeb4ca22b167f5aca18f1846c615e4d0fdd0472a02473044022041f91789dbdd823d8d9dcdc22ae054b5a338eb3e80c844bc4ddafda5b95fad7a022072215b7e97f8f644de25f8ea636d5b6685b1673a3c13f81bd0f66e18c74b33c5012102598122f3f1062147c117fde19e67af124351ccd6aadfd50926d7678ecc5526d700000000",
  "blockHash": "0000000000000000000143aefa609bb22efbfb8cbcba4e608293f97723373613",
  "transactionIndex": "",
  "type": "1",
  "receipt": ""
}

`
	tx, err := ParseTx([]byte(str), 300)
	if err != nil {
		t.Error(err)
	} else {
		t.Log(tx)
	}

}
