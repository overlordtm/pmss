package multihasher

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
)

type Hashes struct {
	MD5    string
	SHA1   string
	SHA256 string
}

func Hash(r io.Reader) (Hashes, error) {
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	// tlsHash := tlsh.New()
	// ssdeepHash := ssdeep.New()

	reader := bufio.NewReaderSize(r, os.Getpagesize())
	multiWriter := io.MultiWriter(md5, sha1, sha256)

	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return Hashes{}, err
	}

	return Hashes{
		MD5:    hex.EncodeToString(md5.Sum(nil)),
		SHA1:   hex.EncodeToString(sha1.Sum(nil)),
		SHA256: hex.EncodeToString(sha256.Sum(nil)),
		// TLSH:   tlsHash.String(),
		// SSDeep: hex.EncodeToString(ssdeepHash.Sum(nil)),
	}, nil
}
