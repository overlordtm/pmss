package scanner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"errors"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/overlordtm/pmss/internal/apiclient"
	"github.com/overlordtm/pmss/pkg/client"
	"github.com/overlordtm/pmss/pkg/multihasher"
	"github.com/sirupsen/logrus"
)

type Option func(*Scanner)

func WithExcludePaths(paths []string) Option {
	return func(s *Scanner) {
		s.excludedPaths = paths
	}
}

type Scanner struct {
	maxSize       int64
	client        client.Client
	workers       int
	excludedPaths []string
}

type scanItem struct {
	path string
	info os.FileInfo
}

func New(opts ...Option) *Scanner {

	s := &Scanner{
		excludedPaths: make([]string, 0, 0),
		workers:       runtime.NumCPU() * 2,
	}

	for _, o := range opts {
		o(s)
	}

	return s
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

func (s *Scanner) shouldScan(pth string, info fs.FileInfo) bool {

	if info != nil && !info.Mode().IsRegular() {
		return false
	}

	for _, patt := range s.excludedPaths {

		if strings.HasSuffix(patt, string(filepath.Separator)) && strings.HasPrefix(pth, filepath.FromSlash(patt)) {
			return false
		}

		match, err := doublestar.PathMatch(patt, pth)
		if err != nil {
			logrus.WithError(err).WithField("pattern", patt).Fatal("Bad exclude pattern")
			return false
		}
		if match {
			return false
		}
	}

	return true
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

			if !info.IsDir() && s.shouldScan(path, info) {
				queue <- scanItem{path: path, info: info}
			}

			// if !info.IsDir() && info.Mode().IsRegular() && !(strings.HasPrefix(path, "/proc/") || strings.HasPrefix(path, "/sys/")) {
			// 	// scan all regulrar files
			// 	queue <- scanItem{path: path, info: info}
			// }
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
