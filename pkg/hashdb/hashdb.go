package hashdb

type HashDb interface {
	FindByMD5(string) error
	BulkInsert([]WhitelistRow) error
}
