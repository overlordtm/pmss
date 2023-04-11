package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"errors"

	"github.com/overlordtm/pmss/internal/apiclient"
	"github.com/overlordtm/pmss/pkg/client"
	"github.com/overlordtm/pmss/pkg/multihasher"
)

type Scanner struct {
	client  client.Client
	workers int
}

type scanItem struct {
	path string
	info os.FileInfo
}

func New() *Scanner {
	return &Scanner{
		workers: runtime.NumCPU() * 2,
	}
}

// Scan scans the given paths and returns the results. Path can be either directory or file
func (s *Scanner) Scan(results chan apiclient.File, paths ...string) (err error) {

	wg := sync.WaitGroup{}

	for _, pth2 := range paths {
		pth := pth2
		if info, err2 := os.Stat(pth); err2 != nil {
			err = errors.Join(err, fmt.Errorf("error while getting file info %s: %v", pth, err2))
		} else {
			wg.Add(1)
			if info.IsDir() {
				go func() {
					defer wg.Done()
					s.scanDir(results, pth)
				}()
			} else {
				go func() {
					defer wg.Done()
					r, err := s.scanFile(pth, info)

					if err != nil {
						err = errors.Join(err, err)
						return
					}
					results <- r
				}()
			}
		}
	}

	wg.Wait()
	close(results)

	return err

}

func (s *Scanner) scanDir(results chan apiclient.File, dir string) (err error) {

	queue := make(chan scanItem, 1024)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// panic(err)
				return filepath.SkipDir
				// return err
			}
			if !info.IsDir() && info.Mode().IsRegular() {
				// scan all regulrar files
				queue <- scanItem{path: path, info: info}
			}
			return nil
		})
		close(queue)
	}()

	wg.Add(s.workers)
	for i := 0; i < s.workers; i++ {
		go func() {
			defer wg.Done()
			for item := range queue {
				r, err1 := s.scanFile(item.path, item.info)
				if err1 != nil {
					err = errors.Join(err, err1)
					continue
				}
				results <- r
			}
		}()
	}

	wg.Wait()

	return err
}

func (s *Scanner) scanFile(path string, info os.FileInfo) (r apiclient.File, err error) {

	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return r, fmt.Errorf("error while opening file: %v", err)
	}
	defer f.Close()

	h, err := multihasher.Hash(f)
	if err != nil {
		return r, fmt.Errorf("error while hashing file: %v", err)
	}

	return apiclient.File{
		Path:     path,
		Size:     info.Size(),
		Md5:      &h.MD5,
		Sha1:     &h.SHA1,
		Sha256:   &h.SHA256,
		FileMode: uint32(info.Mode()),
		Mtime:    info.ModTime().Unix(),
	}, nil

}
