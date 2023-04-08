package client

import (
	"net/http"
	"os"
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

type Files []FileFeatures

type Client interface {
	// ScanFeatures(f FileFeatures) (detector.Result, error)
	SubmitFiles(files Files)
}

type HTTPClient struct {
	httpClient *http.Client
}

func New() *HTTPClient {
	return &HTTPClient{
		httpClient: &http.Client{},
	}
}
