package service

type DbMonitorAddressInterface interface {
	AddMonitorAddress(blockchain int64, address *MonitorAddress) error
	GetAddressByToken(blockchain int64, token string) ([]*MonitorAddress, error)
	NewTx(blockchain int64,tx []*Tx) error
	NewBlock(blockchain int64,block []*Block) error
	NewReceipt(blockchain int64,receipt []*Receipt) error
}
