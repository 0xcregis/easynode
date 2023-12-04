package chain

import "testing"

func TestGetCoreAddress(t *testing.T) {
	addr := GetCoreAddress(200, "0xb6e268b6675846104feef5582d22f40723164d05")
	t.Log(addr)
}
