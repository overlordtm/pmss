package checker

import (
	"os"

	"github.com/overlordtm/pmss/pkg/multihasher"
)

type HashChecker interface {
	CheckHash(multihasher.Hashes) (*Result, error)
}

type FileInfoChecker interface {
	CheckInfo(os.FileInfo) (*Result, error)
}
