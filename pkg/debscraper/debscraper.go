package debscraper

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/isbm/go-deb"
	"github.com/sirupsen/logrus"
)

type DebScraper struct {
	opts       options
	httpClient *http.Client
}

type PackageInfo struct {
	Name         string
	Version      string
	Architecture string
	Filename     string
	MD5Sum       string
	SHA256       string
}

type HashItem struct {
	Package      string
	Version      string
	Architecture string
	Filename     string
	Size         int64
	Mode         os.FileMode
	Owner        string
	Group        string
	MD5          string
	SHA1         string
	SHA256       string
}

type options struct {
	// Mirror URL
	MirrorUrl string
	// Distribution name
	Distro string
	// Component name
	Component string
	// Architecture
	Arch string
}

type option func(*options)

func WithMirrorUrl(mirrorUrl string) func(*options) {
	return func(o *options) {
		o.MirrorUrl = mirrorUrl
	}
}

func WithDistro(distro string) func(*options) {
	return func(o *options) {
		o.Distro = distro
	}
}

func WithComponent(component string) func(*options) {
	return func(o *options) {
		o.Component = component
	}
}

func WithArch(arch string) func(*options) {
	return func(o *options) {
		o.Arch = arch
	}
}

func New(opts ...option) *DebScraper {

	options := options{
		MirrorUrl: "http://ftp.debian.org/debian",
		Distro:    "buster",
		Component: "main",
		Arch:      "amd64",
	}

	for _, opt := range opts {
		opt(&options)
	}

	return &DebScraper{
		opts:       options,
		httpClient: &http.Client{},
	}
}

func (s *DebScraper) listPackages() ([]PackageInfo, error) {
	res, err := s.httpClient.Get(fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages.gz", s.opts.MirrorUrl, s.opts.Distro, s.opts.Component, s.opts.Arch))
	if err != nil {
		return nil, fmt.Errorf("error while fetching Packages.gz: %w", err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status while fetching Packages.gz: %s", res.Status)
	}

	gzipReader, err := gzip.NewReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error while creating gzip reader: %w", err)
	}

	defer gzipReader.Close()

	r := bufio.NewScanner(gzipReader)

	pkgs := make([]PackageInfo, 0)

	pkg := PackageInfo{}

	for r.Scan() {
		line := r.Text()
		if strings.HasPrefix(line, "Package: ") {
			pkg.Name = strings.TrimPrefix(line, "Package: ")
		} else if strings.HasPrefix(line, "Version: ") {
			pkg.Version = strings.TrimPrefix(line, "Version: ")
		} else if strings.HasPrefix(line, "Architecture: ") {
			pkg.Architecture = strings.TrimPrefix(line, "Architecture: ")
		} else if strings.HasPrefix(line, "Filename: ") {
			pkg.Filename = strings.TrimPrefix(line, "Filename: ")
		} else if strings.HasPrefix(line, "MD5Sum: ") {
			pkg.MD5Sum = strings.TrimPrefix(line, "MD5Sum: ")
		} else if strings.HasPrefix(line, "SHA256: ") {
			pkg.SHA256 = strings.TrimPrefix(line, "SHA256: ")
		} else if line == "" {
			pkgs = append(pkgs, pkg)
			pkg = PackageInfo{}
		}
	}

	return pkgs, r.Err()
}

func (s *DebScraper) fetchPackage(pkgInfo PackageInfo) ([]HashItem, error) {

	options := &deb.PackageOptions{
		Hash:                 deb.HASH_SHA1,
		RecalculateChecksums: true,
		MetaOnly:             false,
	}

	pkgUri := fmt.Sprintf("%s/%s", s.opts.MirrorUrl, pkgInfo.Filename)

	pkg, err := deb.OpenPackageFile(pkgUri, options)

	if err != nil {
		return nil, fmt.Errorf("error while opening package file: %w", err)
	}

	res := make([]HashItem, 0)

	for _, f := range pkg.Files() {
		md5 := pkg.GetFileMd5Sums(f.Name())
		sha1 := pkg.GetFileChecksum(f.Name())

		if !f.IsDir() && (md5 != "" || sha1 != "") {
			res = append(res, HashItem{
				Package:      pkgInfo.Name,
				Version:      pkgInfo.Version,
				Architecture: pkgInfo.Architecture,
				Filename:     f.Name(),
				Mode:         f.Mode(),
				Size:         f.Size(),
				Owner:        f.Owner(),
				Group:        f.Group(),
				MD5:          md5,
				SHA1:         sha1,
			})
		}
	}

	return res, nil
}

func (s *DebScraper) Scrape(concurrency int, hashItemCh chan HashItem) error {
	pkgs, err := s.listPackages()
	if err != nil {
		return fmt.Errorf("error while listing packages: %w", err)
	}

	workCh := make(chan PackageInfo, len(pkgs))

	go func() {
		for _, pkg := range pkgs {
			workCh <- pkg
		}
		close(workCh)
	}()

	wg := sync.WaitGroup{}
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for pkg := range workCh {
				hashItems, err := s.fetchPackage(pkg)
				if err != nil {
					logrus.WithError(err).WithField("package", pkg.Name).Error("error while fetching package")
				} else {
					for _, item := range hashItems {
						hashItemCh <- item
					}
				}
			}
		}()
	}

	wg.Wait()
	return nil
}
