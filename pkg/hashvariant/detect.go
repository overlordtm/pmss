package hashvariant

import "fmt"

type HashVariant string

const (
	Unknown HashVariant = ""
	MD5     HashVariant = "md5"
	SHA1    HashVariant = "sha1"
	SHA256  HashVariant = "sha256"
	SHA512  HashVariant = "sha512"
	TLSH    HashVariant = "tlsh"
)

var (
	ErrUnknownHashVariant = fmt.Errorf("unknown hash variant")
)

func DetectHashVariant(hash string) HashVariant {
	switch len(hash) {
	case 32:
		return MD5
	case 40:
		return SHA1
	case 64:
		return SHA256
	case 128:
		return SHA512
	case 72:
		return TLSH
	default:
		return Unknown
	}
}
