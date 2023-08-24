package xrp

import (
	"testing"

	"github.com/0xcregis/easynode/store"
)

func TestGetBlockFromKafka(t *testing.T) {
	value := `
{
  "accepted": "true",
  "account_hash": "910FCF5C2CA222286C6D1D3AB40198D4552B639129A0EE516A45E8B4770EA00E",
  "close_flags": "0",
  "close_time": "746087402",
  "close_time_human": "2023-Aug-23 06:30:02.000000000 UTC",
  "close_time_resolution": "10",
  "closed": "true",
  "hash": "224CE7DA0F1D88ED35559E5D1904B20043DD125AFA21CED62FA66780BDC1B491",
  "ledger_hash": "224CE7DA0F1D88ED35559E5D1904B20043DD125AFA21CED62FA66780BDC1B491",
  "ledger_index": "82032099",
  "parent_close_time": "746087401",
  "parent_hash": "9D4985A3CBD7CACD42C0149C440304C3AFBB1BF3157CA4063D94F6BB8B077614",
  "seqNum": "82032099",
  "totalCoins": "99988480317227197",
  "total_coins": "99988480317227197",
  "transaction_hash": "AFA1C7E4CD2956BDCF495F597B7EEE0D96457399FA59DFFA190FF20088EE81DE"
}
`
	block, err := GetBlockFromKafka([]byte(value))
	if err != nil {
		t.Error(err)
	}
	t.Logf("block:%+v", block)
}

func TestGetTxFromKafka(t *testing.T) {

	value := `
{
  "block": "{\"accepted\":\"true\",\"account_hash\":\"2BD14C2BE88F3ED179007E0EFA752BC36127308D824354C55504E31AC9562C21\",\"close_flags\":\"0\",\"close_time\":\"746087000\",\"close_time_human\":\"2023-Aug-23 06:23:20.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_index\":\"82031992\",\"parent_close_time\":\"746086992\",\"parent_hash\":\"8F058D0BA35084C7989834CA1AF3AA195AF47A850FAEB4148C4704C20D887BD0\",\"seqNum\":\"82031992\",\"totalCoins\":\"99988480325556972\",\"total_coins\":\"99988480325556972\",\"transaction_hash\":\"CFD845FDF68FAB6553ADA7F398301C7E347980B00CEC98C5CB011F84F50E34AA\"}",
  "blockHash": "FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6",
  "blockNumber": "82031992",
  "hash": "EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7",
  "tx": "{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Amount\":\"2299401567509\",\"Destination\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Fee\":\"10000\",\"Sequence\":105320,\"SigningPubKey\":\"037DFC657B370962A83DA8035E05B240D2A921AEF5FD59B904D71E3D5B0E3C9D9C\",\"TransactionType\":\"Payment\",\"TxnSignature\":\"3044022075005F0CAFB68AE7FE135AEA8F2496AC4FF1B1EA9DC230D35E3545FDB623EFE3022031B0696D319E3D665E3DF14EC41A2841F60EECE1E9A54817A17E208C11A837B4\",\"hash\":\"EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7\",\"metaData\":{\"AffectedNodes\":[{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Balance\":\"32155814798353\",\"Flags\":0,\"OwnerCount\":1,\"Sequence\":66878203},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"3BFAEEAD933C62413BD9A92BBC11993F2B0FCA28E2CD8FEE69EB8C22EE848E36\",\"PreviousFields\":{\"Balance\":\"29856413230844\"},\"PreviousTxnID\":\"A0BCDF69185C351155D58AC7C4E5C19061637078107EBC33990FE65F12503FD7\",\"PreviousTxnLgrSeq\":82031976}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Balance\":\"1780244485\",\"Flags\":131072,\"OwnerCount\":1,\"Sequence\":105321},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"C19B36F6B6F2EEC9F4E2AF875E533596503F4541DBA570F06B26904FDBBE9C52\",\"PreviousFields\":{\"Balance\":\"2301181821994\",\"Sequence\":105320},\"PreviousTxnID\":\"C72F3FFDE1E59E8DE391A4E850E9031AE8D33156A7A2EBD1E3E355A364E527B7\",\"PreviousTxnLgrSeq\":82031989}}],\"TransactionIndex\":19,\"TransactionResult\":\"tesSUCCESS\",\"delivered_amount\":\"2299401567509\"}}"
}
  `

	tx, err := GetTxFromKafka([]byte(value))
	if err != nil {
		t.Error(err)
	}
	t.Logf("tx:%+v", tx)
}

func TestGetReceiptFromKafka(t *testing.T) {

	value := `

	{
	    "account":"rf4VZ6LYKnPY1uaEiZF8qzD8vrW91vvjbb",
	    "date":745986661,
	    "hash":"89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5",
	    "ledgerIndex":82005629,
	    "status":"success",
	    "transactionIndex":18,
	    "transactionResult":"tecINSUFFICIENT_RESERVE"
	}
`
	r, err := GetReceiptFromKafka([]byte(value))
	if err != nil {
		t.Error(err)
	}
	t.Logf("receipt:%+v", r)
}

func TestParseTx(t *testing.T) {

	value := `
{
  "block": "{\"accepted\":\"true\",\"account_hash\":\"2BD14C2BE88F3ED179007E0EFA752BC36127308D824354C55504E31AC9562C21\",\"close_flags\":\"0\",\"close_time\":\"746087000\",\"close_time_human\":\"2023-Aug-23 06:23:20.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_index\":\"82031992\",\"parent_close_time\":\"746086992\",\"parent_hash\":\"8F058D0BA35084C7989834CA1AF3AA195AF47A850FAEB4148C4704C20D887BD0\",\"seqNum\":\"82031992\",\"totalCoins\":\"99988480325556972\",\"total_coins\":\"99988480325556972\",\"transaction_hash\":\"CFD845FDF68FAB6553ADA7F398301C7E347980B00CEC98C5CB011F84F50E34AA\"}",
  "blockHash": "FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6",
  "blockNumber": "82031992",
  "hash": "EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7",
  "tx": "{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Amount\":\"2299401567509\",\"Destination\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Fee\":\"10000\",\"Sequence\":105320,\"SigningPubKey\":\"037DFC657B370962A83DA8035E05B240D2A921AEF5FD59B904D71E3D5B0E3C9D9C\",\"TransactionType\":\"Payment\",\"TxnSignature\":\"3044022075005F0CAFB68AE7FE135AEA8F2496AC4FF1B1EA9DC230D35E3545FDB623EFE3022031B0696D319E3D665E3DF14EC41A2841F60EECE1E9A54817A17E208C11A837B4\",\"hash\":\"EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7\",\"metaData\":{\"AffectedNodes\":[{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Balance\":\"32155814798353\",\"Flags\":0,\"OwnerCount\":1,\"Sequence\":66878203},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"3BFAEEAD933C62413BD9A92BBC11993F2B0FCA28E2CD8FEE69EB8C22EE848E36\",\"PreviousFields\":{\"Balance\":\"29856413230844\"},\"PreviousTxnID\":\"A0BCDF69185C351155D58AC7C4E5C19061637078107EBC33990FE65F12503FD7\",\"PreviousTxnLgrSeq\":82031976}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Balance\":\"1780244485\",\"Flags\":131072,\"OwnerCount\":1,\"Sequence\":105321},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"C19B36F6B6F2EEC9F4E2AF875E533596503F4541DBA570F06B26904FDBBE9C52\",\"PreviousFields\":{\"Balance\":\"2301181821994\",\"Sequence\":105320},\"PreviousTxnID\":\"C72F3FFDE1E59E8DE391A4E850E9031AE8D33156A7A2EBD1E3E355A364E527B7\",\"PreviousTxnLgrSeq\":82031989}}],\"TransactionIndex\":19,\"TransactionResult\":\"tesSUCCESS\",\"delivered_amount\":\"2299401567509\"}}"
}
  `

	s, err := ParseTx([]byte(value), store.PolygonTopic, 301)
	if err != nil {
		t.Error(err)
	}
	t.Logf("sub.tx:%v", s)

}

func TestGetTxType(t *testing.T) {

	value := `
{
  "block": "{\"accepted\":\"true\",\"account_hash\":\"2BD14C2BE88F3ED179007E0EFA752BC36127308D824354C55504E31AC9562C21\",\"close_flags\":\"0\",\"close_time\":\"746087000\",\"close_time_human\":\"2023-Aug-23 06:23:20.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_index\":\"82031992\",\"parent_close_time\":\"746086992\",\"parent_hash\":\"8F058D0BA35084C7989834CA1AF3AA195AF47A850FAEB4148C4704C20D887BD0\",\"seqNum\":\"82031992\",\"totalCoins\":\"99988480325556972\",\"total_coins\":\"99988480325556972\",\"transaction_hash\":\"CFD845FDF68FAB6553ADA7F398301C7E347980B00CEC98C5CB011F84F50E34AA\"}",
  "blockHash": "FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6",
  "blockNumber": "82031992",
  "hash": "EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7",
  "tx": "{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Amount\":\"2299401567509\",\"Destination\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Fee\":\"10000\",\"Sequence\":105320,\"SigningPubKey\":\"037DFC657B370962A83DA8035E05B240D2A921AEF5FD59B904D71E3D5B0E3C9D9C\",\"TransactionType\":\"Payment\",\"TxnSignature\":\"3044022075005F0CAFB68AE7FE135AEA8F2496AC4FF1B1EA9DC230D35E3545FDB623EFE3022031B0696D319E3D665E3DF14EC41A2841F60EECE1E9A54817A17E208C11A837B4\",\"hash\":\"EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7\",\"metaData\":{\"AffectedNodes\":[{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Balance\":\"32155814798353\",\"Flags\":0,\"OwnerCount\":1,\"Sequence\":66878203},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"3BFAEEAD933C62413BD9A92BBC11993F2B0FCA28E2CD8FEE69EB8C22EE848E36\",\"PreviousFields\":{\"Balance\":\"29856413230844\"},\"PreviousTxnID\":\"A0BCDF69185C351155D58AC7C4E5C19061637078107EBC33990FE65F12503FD7\",\"PreviousTxnLgrSeq\":82031976}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Balance\":\"1780244485\",\"Flags\":131072,\"OwnerCount\":1,\"Sequence\":105321},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"C19B36F6B6F2EEC9F4E2AF875E533596503F4541DBA570F06B26904FDBBE9C52\",\"PreviousFields\":{\"Balance\":\"2301181821994\",\"Sequence\":105320},\"PreviousTxnID\":\"C72F3FFDE1E59E8DE391A4E850E9031AE8D33156A7A2EBD1E3E355A364E527B7\",\"PreviousTxnLgrSeq\":82031989}}],\"TransactionIndex\":19,\"TransactionResult\":\"tesSUCCESS\",\"delivered_amount\":\"2299401567509\"}}"
}
  `

	tp, err := GetTxType([]byte(value))
	if err != nil {
		t.Error(err)
	}
	t.Logf("tx.type:%v", tp)
}

func TestCheckAddress(t *testing.T) {
	value := `
{
  "block": "{\"accepted\":\"true\",\"account_hash\":\"2BD14C2BE88F3ED179007E0EFA752BC36127308D824354C55504E31AC9562C21\",\"close_flags\":\"0\",\"close_time\":\"746087000\",\"close_time_human\":\"2023-Aug-23 06:23:20.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_hash\":\"FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6\",\"ledger_index\":\"82031992\",\"parent_close_time\":\"746086992\",\"parent_hash\":\"8F058D0BA35084C7989834CA1AF3AA195AF47A850FAEB4148C4704C20D887BD0\",\"seqNum\":\"82031992\",\"totalCoins\":\"99988480325556972\",\"total_coins\":\"99988480325556972\",\"transaction_hash\":\"CFD845FDF68FAB6553ADA7F398301C7E347980B00CEC98C5CB011F84F50E34AA\"}",
  "blockHash": "FFF79B99B28E446652B741A682FA40D4C1E485BBE78EA5D3868223C60E6C93A6",
  "blockNumber": "82031992",
  "hash": "EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7",
  "tx": "{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Amount\":\"2299401567509\",\"Destination\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Fee\":\"10000\",\"Sequence\":105320,\"SigningPubKey\":\"037DFC657B370962A83DA8035E05B240D2A921AEF5FD59B904D71E3D5B0E3C9D9C\",\"TransactionType\":\"Payment\",\"TxnSignature\":\"3044022075005F0CAFB68AE7FE135AEA8F2496AC4FF1B1EA9DC230D35E3545FDB623EFE3022031B0696D319E3D665E3DF14EC41A2841F60EECE1E9A54817A17E208C11A837B4\",\"hash\":\"EB9F83623CE7CE8802B6D1D12FB55133437005AAB2E4E09BA61C2B1EAE5751B7\",\"metaData\":{\"AffectedNodes\":[{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rMvCasZ9cohYrSZRNYPTZfoaaSUQMfgQ8G\",\"Balance\":\"32155814798353\",\"Flags\":0,\"OwnerCount\":1,\"Sequence\":66878203},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"3BFAEEAD933C62413BD9A92BBC11993F2B0FCA28E2CD8FEE69EB8C22EE848E36\",\"PreviousFields\":{\"Balance\":\"29856413230844\"},\"PreviousTxnID\":\"A0BCDF69185C351155D58AC7C4E5C19061637078107EBC33990FE65F12503FD7\",\"PreviousTxnLgrSeq\":82031976}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJn2zAPdFA193sixJwuFixRkYDUtx3apQh\",\"Balance\":\"1780244485\",\"Flags\":131072,\"OwnerCount\":1,\"Sequence\":105321},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"C19B36F6B6F2EEC9F4E2AF875E533596503F4541DBA570F06B26904FDBBE9C52\",\"PreviousFields\":{\"Balance\":\"2301181821994\",\"Sequence\":105320},\"PreviousTxnID\":\"C72F3FFDE1E59E8DE391A4E850E9031AE8D33156A7A2EBD1E3E355A364E527B7\",\"PreviousTxnLgrSeq\":82031989}}],\"TransactionIndex\":19,\"TransactionResult\":\"tesSUCCESS\",\"delivered_amount\":\"2299401567509\"}}"
}
  `

	r := CheckAddress([]byte(value), nil, "")

	if !r {
		t.Error("fail")
	} else {
		t.Log("ok")
	}

}
