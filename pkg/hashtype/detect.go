package hashtype

import (
	"fmt"
)

type Type string

const (
	Unknown Type = ""
	MD5     Type = "md5"
	SHA1    Type = "sha1"
	SHA256  Type = "sha256"
	SHA512  Type = "sha512"
	TLSH    Type = "tlsh"
)

var (
	ErrUnknown = fmt.Errorf("unknown hash variant")
)

func Detect(hash string) Type {
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
