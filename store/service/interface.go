package service

type DbMonitorAddressInterface interface {
	AddMonitorAddress(blockchain int64, address *MonitorAddress) error
	GetAddressByToken(blockchain int64, token string) ([]*MonitorAddress, error)
	NewTx(tx []*Tx) error
	NewBlock(block []*Block) error
	NewReceipt(receipt []*Receipt) error
}
