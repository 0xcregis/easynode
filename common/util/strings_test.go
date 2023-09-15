package util

import "testing"

func TestDiv(t *testing.T) {
	t.Log(Div("1000", 2))
}

func TestNftData(t *testing.T) {
	t.Log(NftData("0xabe68307e498ae6cbe979c23ebd518e8e3e04d26000000000000000000000003000000000000000000000000000000000000000000000000000000000000000d"))
}
