package store

type DbStoreInterface interface {
	NewToken(token *NodeToken) error
	UpdateToken(token string, nodeToken *NodeToken) error
	GetNodeTokenByEmail(email string) (*NodeToken, error)
	GetNodeTokenByToken(token string) (*NodeToken, error)
	AddMonitorAddress(blockchain int64, address *MonitorAddress) error
	GetAddressByToken(blockchain int64, token string) ([]*MonitorAddress, error)
	GetAddressByToken2(token string) ([]*MonitorAddress, error)
	GetAddressByToken3(blockchain int64) ([]*MonitorAddress, error)
	DelMonitorAddress(blockchain int64, token string, address string) error
	NewTx(blockchain int64, tx []*Tx) error
	NewBlock(blockchain int64, block []*Block) error
	NewReceipt(blockchain int64, receipt []*Receipt) error
	NewSubTx(blockchain int64, tx []*SubTx) error
	NewBackupTx(blockchain int64, tx []*BackupTx) error

	NewSubFilter(filters []*SubFilter) error
	DelSubFilter(id int64) error
	DelSubFilter2(filter *SubFilter) error
	GetSubFilter(token string, blockChain int64, txCode string) ([]*SubFilter, error)
}
