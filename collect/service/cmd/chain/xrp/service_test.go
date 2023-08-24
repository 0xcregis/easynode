package xrp

import (
	"testing"

	"github.com/0xcregis/easynode/collect"
	"github.com/0xcregis/easynode/collect/config"
	"github.com/0xcregis/easynode/collect/service/db"
	"github.com/sirupsen/logrus"
	"github.com/sunjiangjun/xlog"
)

func Init() (collect.BlockChainInterface, config.Config, *xlog.XLog) {
	cfg := config.LoadConfig("./../../../../../cmd/collect/config_xrp.json")
	x := xlog.NewXLogger()
	store := db.NewTaskCacheService(cfg.Chains[0], x)
	return NewService(cfg.Chains[0], x, store, "9587acc2-04ab-4154-ae11-f6d588c6493f", "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"), cfg, x
}

func TestService_GetBlockByNumber(t *testing.T) {
	s, _, x := Init()
	b, txList := s.GetBlockByNumber("81984330", x.WithFields(logrus.Fields{}), false)
	t.Logf("block.hash:%v, tx.list:%v", b.BlockHash, len(txList))
}

func TestService_GetMultiBlockByNumber(t *testing.T) {
	s, _, x := Init()
	list, _ := s.GetMultiBlockByNumber("81984330", x.WithFields(logrus.Fields{}), false)
	t.Logf("block.list:%v", len(list))
}

func TestService_GetBlockByHash(t *testing.T) {
	s, _, x := Init()
	block, txList := s.GetBlockByHash("B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9", x.WithFields(logrus.Fields{}), true)
	t.Logf("block.hash:%v, tx.list:%v", block.BlockHash, len(txList))
}

func TestService_GetTx(t *testing.T) {
	s, _, x := Init()
	tx := s.GetTx("89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5", x.WithFields(logrus.Fields{}))
	t.Log(tx)
}

func TestService_GetReceipt(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceipt("89CEEBAF42BE602DCFDF0F89B7B9111A4E09CF4D55EC0BFEF5C439246E31C8B5", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}

func TestService_CheckAddress(t *testing.T) {
	s, _, _ := Init()
	str := `
	{"block":"{\"accepted\":\"true\",\"account_hash\":\"836BC07C9FF83ED928979F76A80D8461412A2A69632E6280C7D2CA1AFC965B71\",\"close_flags\":\"0\",\"close_time\":\"745905191\",\"close_time_human\":\"2023-Aug-21 03:53:11.000000000 UTC\",\"close_time_resolution\":\"10\",\"closed\":\"true\",\"hash\":\"B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9\",\"ledger_hash\":\"B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9\",\"ledger_index\":\"81984330\",\"parent_close_time\":\"745905190\",\"parent_hash\":\"E59E77C619E89E8608B93F9AD11ACA6541BB9CB552093917088CCF8102DB975F\",\"seqNum\":\"81984330\",\"totalCoins\":\"99988485440853202\",\"total_coins\":\"99988485440853202\",\"transaction_hash\":\"49C46B8FB727DDF65EFB626013A1BC790ACC326682A3257DDA13306BC333E1F9\"}","blockHash":"B4468DFE533796FDBD54D324626CFC648979EFB08D796840D5F0CCDB9FD655F9","blockNumber":"81984330","hash":"05B881C0A86CBDF930D6F2EB684D365CAACED0D77E95D61F66FE77C1DFBBDE37","tx":"{\"Account\":\"rJumr5e1HwiuV543H7bqixhtFreChWTaHH\",\"Fee\":\"15\",\"Flags\":0,\"LastLedgerSequence\":81984332,\"OfferSequence\":27088391,\"Sequence\":27088396,\"SigningPubKey\":\"02207B680D408E4298BF14EC14CD2751DB20CDE1ADAD6322FBD4239B84F27D3F27\",\"TakerGets\":\"14619558002\",\"TakerPays\":{\"currency\":\"CNY\",\"issuer\":\"rJ1adrpGS3xsnQMb9Cw54tWJVFPuSdZHK\",\"value\":\"96645.21969576311\"},\"TransactionType\":\"OfferCreate\",\"TxnSignature\":\"3045022100DD70EE6BCC5096EC8D41A50E9FD872537786040AC6D1B1EEA1F8851E262051AB02205C981ABDD0E8686A29A0466C67515C7705B88A97DAB0C8562C74DDEDE867180D\",\"hash\":\"05B881C0A86CBDF930D6F2EB684D365CAACED0D77E95D61F66FE77C1DFBBDE37\",\"metaData\":{\"AffectedNodes\":[{\"CreatedNode\":{\"LedgerEntryType\":\"Offer\",\"LedgerIndex\":\"3EC5A7D7C1017FEBEFCB50B958F094F7AF0A12D2F3CED398AA2D05885F0C0674\",\"NewFields\":{\"Account\":\"rJumr5e1HwiuV543H7bqixhtFreChWTaHH\",\"BookDirectory\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F177C60E122ECC4\",\"Sequence\":27088396,\"TakerGets\":\"14619558002\",\"TakerPays\":{\"currency\":\"CNY\",\"issuer\":\"rJ1adrpGS3xsnQMb9Cw54tWJVFPuSdZHK\",\"value\":\"96645.21969576311\"}}}},{\"DeletedNode\":{\"FinalFields\":{\"ExchangeRate\":\"4f144254728891a9\",\"Flags\":0,\"RootIndex\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F144254728891A9\",\"TakerGetsCurrency\":\"0000000000000000000000000000000000000000\",\"TakerGetsIssuer\":\"0000000000000000000000000000000000000000\",\"TakerPaysCurrency\":\"000000000000000000000000434E590000000000\",\"TakerPaysIssuer\":\"0360E3E0751BD9A566CD03FA6CAFC78118B82BA0\"},\"LedgerEntryType\":\"DirectoryNode\",\"LedgerIndex\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F144254728891A9\"}},{\"CreatedNode\":{\"LedgerEntryType\":\"DirectoryNode\",\"LedgerIndex\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F177C60E122ECC4\",\"NewFields\":{\"ExchangeRate\":\"4f177c60e122ecc4\",\"RootIndex\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F177C60E122ECC4\",\"TakerPaysCurrency\":\"000000000000000000000000434E590000000000\",\"TakerPaysIssuer\":\"0360E3E0751BD9A566CD03FA6CAFC78118B82BA0\"}}},{\"ModifiedNode\":{\"FinalFields\":{\"Account\":\"rJumr5e1HwiuV543H7bqixhtFreChWTaHH\",\"Balance\":\"5029821234\",\"Flags\":0,\"OwnerCount\":5,\"Sequence\":27088397},\"LedgerEntryType\":\"AccountRoot\",\"LedgerIndex\":\"649D8F54A04A6503F1B768A1C5A8B7690722F4A96AB1D999261B2661D8223F2D\",\"PreviousFields\":{\"Balance\":\"5029821249\",\"Sequence\":27088396},\"PreviousTxnID\":\"B61B1D8F8E7589C6EB956284050F041DA0541F07728AA4335E79A27E8B0B1E1A\",\"PreviousTxnLgrSeq\":81984330}},{\"ModifiedNode\":{\"FinalFields\":{\"Flags\":0,\"IndexNext\":\"0\",\"IndexPrevious\":\"0\",\"Owner\":\"rJumr5e1HwiuV543H7bqixhtFreChWTaHH\",\"RootIndex\":\"B85E96605A78ABA153E322936C277B4E1524790EE1A340575FA2088E2F6E8F4F\"},\"LedgerEntryType\":\"DirectoryNode\",\"LedgerIndex\":\"B85E96605A78ABA153E322936C277B4E1524790EE1A340575FA2088E2F6E8F4F\"}},{\"DeletedNode\":{\"FinalFields\":{\"Account\":\"rJumr5e1HwiuV543H7bqixhtFreChWTaHH\",\"BookDirectory\":\"623C4C4AD65873DA787AC85A0A1385FE6233B6DE100799474F144254728891A9\",\"BookNode\":\"0\",\"Flags\":0,\"OwnerNode\":\"0\",\"PreviousTxnID\":\"123AE89A7734E975A19DC41E6C7F2E3E7584221113643F4178823A039A3C799C\",\"PreviousTxnLgrSeq\":81984318,\"Sequence\":27088391,\"TakerGets\":\"9592066224\",\"TakerPays\":{\"currency\":\"CNY\",\"issuer\":\"rJ1adrpGS3xsnQMb9Cw54tWJVFPuSdZHK\",\"value\":\"54698.08620206002\"}},\"LedgerEntryType\":\"Offer\",\"LedgerIndex\":\"C11A2C2B49F6A9D855F5399F41E4B04FA44159590EDAE8A4699401963FFA41E9\"}}],\"TransactionIndex\":25,\"TransactionResult\":\"tesSUCCESS\"},\"owner_funds\":\"5009821234\"}"}
    `
	t.Log(s.CheckAddress([]byte(str), nil))
}

func TestService_GetReceiptByBlock(t *testing.T) {
	s, _, x := Init()
	r, err := s.GetReceiptByBlock("", "81984330", x.WithFields(logrus.Fields{}))
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("%+v", r)
	}
}
