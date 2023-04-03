package hashstore

type HashDb interface {
	Whitelist() Whitelist
	Blacklist() Blacklist
}

type Whitelist interface {
	FindByMD5(string) (*WhitelistRow, error)
	InsertBatch([]WhitelistRow) error
	Insert(WhitelistRow) error
}

type Blacklist interface {
	FindByMD5(string) (*BlacklistRow, error)
	InsertBatch([]BlacklistRow) error
	Insert(BlacklistRow) error
}
