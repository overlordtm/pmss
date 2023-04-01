package sigdb

type Item struct {
	MD5    string
	SHA1   string
	SHA256 string

	ImpHash string
	SSDeep  string
	TLSH    string

	Signature string
	Filename  string
	MimeType  string
}

type SigDb interface {
	SaveItem(item Item) error
	FindByMD5(md5 string) (*Item, error)
	FindBySHA1(sha1 string) (*Item, error)
	FindBySHA256(sha256 string) (*Item, error)
	FindByImpHash(imphash string) (*Item, error)
	FindBySSDeep(ssdeep string) (*Item, error)
	FindByTLSH(tlsh string) (*Item, error)
}
