package client

import (
	"os"

	"github.com/overlordtm/pmss/pkg/detector"
)

type FileFeatures struct {
	Path     string
	MD5      string
	SHA1     string
	SHA256   string
	TLSH     string
	FileMode os.FileMode
	Mtime    int64
	Ctime    int64
	Atime    int64
	Size     int64
}

type Client interface {
	ScanFeatures(f FileFeatures) (detector.Result, error)
}
