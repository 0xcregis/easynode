package btc

import (
	"testing"
)

func TestGetBalance(t *testing.T) {

	list, err := GetBalance("https://www.oklink.com/api/v5/explorer/address/utxo", "a6938ee9-4678-4cd1-90ca-ba13ee472ede", "1PL6qjNjEMRhTLAnHEFJWwnvjjKGWAwFws")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(list)
	}
}

func TestGetBalanceEx(t *testing.T) {

	list, err := GetBalanceEx([]string{"https://www.oklink.com/api/v5/explorer/address/utxo#a6938ee9-4678-4cd1-90ca-ba13ee472ede"}, "1PL6qjNjEMRhTLAnHEFJWwnvjjKGWAwFws")
	if err != nil {
		t.Error(err)
	} else {
		t.Log(list)
	}
}
