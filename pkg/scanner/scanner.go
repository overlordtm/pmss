package scanner

import (
	"fmt"
	"os"
	"path/filepath"

	"errors"

	"github.com/overlordtm/pmss/pkg/client"
	"github.com/overlordtm/pmss/pkg/detector"
	"github.com/overlordtm/pmss/pkg/multihasher"
)

type Scanner struct {
	client client.Client
}

type scanItem struct {
	path string
	info os.FileInfo
}

func New(c client.Client) *Scanner {
	return &Scanner{
		client: c,
	}
}

func (s *Scanner) Scan(dir string) (results []detector.Result, err error) {
	queue := make(chan scanItem, 1024)

	go func() {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				// scan everything that is not a directory
				queue <- scanItem{path: path, info: info}
			}
			return nil
		})
		close(queue)
	}()

	for item := range queue {
		r, err1 := s.scanFile(item.path, item.info)
		if err1 != nil {
			err = errors.Join(err, err1)
			continue
		}
		results = append(results, r)
	}

	return results, err
}

func (s *Scanner) scanFile(path string, info os.FileInfo) (r detector.Result, err error) {

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return r, fmt.Errorf("error while opening file: %v", err)
	}
	defer f.Close()

	h, err := multihasher.Hash(f)

	feat := client.FileFeatures{
		Path:   path,
		Size:   info.Size(),
		MD5:    h.MD5,
		SHA1:   h.SHA1,
		SHA256: h.SHA256,
	}

	return s.client.ScanFeatures(feat)
}
