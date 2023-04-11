package debscraper

import (
	"bufio"
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/isbm/go-deb"
	"github.com/overlordtm/pmss/internal/utils"
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/sirupsen/logrus"
)

type dummyLogger struct{}

func (d *dummyLogger) Printf(format string, args ...interface{}) {}
func (d *dummyLogger) Println(args ...interface{})               {}
func (d *dummyLogger) Print(args ...interface{})                 {}

func init() {
	deb.SetLogger(&dummyLogger{})
}

type DebScraper struct {
	mirrorUrl mirrorFunc
	distro    string
	component string
	arch      string
	osType    datastore.OsType

	httpClient *http.Client
	logger     *logrus.Logger
}

func New(opts ...Option) *DebScraper {

	scraper := &DebScraper{
		mirrorUrl: RandomMirror(Mirrors...),
		distro:    "buster",
		component: "main",
		arch:      "amd64",
		osType:    datastore.OsTypeDebian,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		logger: logrus.New(),
	}

	for _, opt := range opts {
		opt(scraper)
	}

	return scraper
}

func (s *DebScraper) ListPackages(ctx context.Context) ([]datastore.Package, error) {

	mirrorUrl := s.mirrorUrl()

	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages.gz", mirrorUrl, s.distro, s.component, s.arch), nil)

	if err != nil {
		return nil, fmt.Errorf("error while creating request: %w", err)
	}

	s.logger.WithField("url", req.URL.String()).Debugln("Fetching Packages.gz")
	res, err := s.httpClient.Do(req)
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

	pkgs := make([]datastore.Package, 0)

	pkg := datastore.Package{}

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

			md5 := strings.TrimSpace(strings.TrimPrefix(line, "MD5Sum: "))

			if md5 != "" {
				pkg.MD5Sum = utils.StringPtr(md5)
			}

		} else if strings.HasPrefix(line, "SHA256: ") {

			sha256 := strings.TrimSpace(strings.TrimPrefix(line, "SHA256: "))
			if sha256 != "" {
				pkg.SHA256 = utils.StringPtr(sha256)
			}

		} else if strings.HasPrefix(line, "Size: ") {
			pkg.Size, _ = strconv.ParseUint(strings.TrimPrefix(line, "Size: "), 10, 64)
		} else if line == "" {

			if pkg.Name == "" || pkg.Version == "" || pkg.Architecture == "" || pkg.Filename == "" || (pkg.MD5Sum == nil && pkg.SHA256 == nil) {
				return nil, errors.New("invalid package info")
			}

			pkg.Distro = s.distro
			pkg.Component = s.component
			pkg.OsType = s.osType
			pkgs = append(pkgs, pkg)
			pkg = datastore.Package{}
		}
	}

	return pkgs, r.Err()
}

func (s *DebScraper) Scrape(ctx context.Context, concurrency int, hashItemCh chan datastore.KnownFile, progress ProgressDelegate) (err error) {
	pkgs, err := s.ListPackages(ctx)
	if err != nil {
		return fmt.Errorf("error while listing packages: %w", err)
	}

	progress.Start(int64(len(pkgs)))

	workCh := make(chan datastore.Package, 0)

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

			for {
				select {
				case <-ctx.Done():
					err = errors.Join(err, ctx.Err())
					return
				case pkg, ok := <-workCh:
					if !ok {
						err = errors.Join(err, fmt.Errorf("work channel closed"))
						return
					}
					hashItems, err := s.FetchPackage(ctx, 5, pkg)

					if err != nil {
						s.logger.WithError(err).WithField("package", pkg.Name).Error("error while fetching package")
					} else {
						for _, item := range hashItems {
							hashItemCh <- item
						}
					}
					progress.Done(1)
				}
			}
		}()
	}

	wg.Wait()
	return err
}

func (s *DebScraper) FetchPackage(ctx context.Context, attempts int, pkg datastore.Package) (items []datastore.KnownFile, err error) {
	for i := 0; i < attempts; i++ {
		items, err := func() ([]datastore.KnownFile, error) {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("panic: %v", r)
				}
			}()
			return s.fetchPackage(ctx, pkg)
		}()

		if err == nil {
			return items, nil
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
	}
	return nil, fmt.Errorf("failed after %d attempts: %w", attempts, err)
}

func (s *DebScraper) fetchPackage(ctx context.Context, pkgInfo datastore.Package) ([]datastore.KnownFile, error) {
	mirrorUrl := s.mirrorUrl()

	pkgUri := fmt.Sprintf("%s/%s", mirrorUrl, pkgInfo.Filename)

	s.logger.WithFields(logrus.Fields{
		"package": pkgInfo.Name,
		"version": pkgInfo.Version,
		"arch":    pkgInfo.Architecture,
		"pkgUri":  pkgUri,
		"size":    humanize.Bytes(pkgInfo.Size),
	}).Info("fetching package")

	options := &deb.PackageOptions{
		Hash:                 deb.HASH_SHA1,
		RecalculateChecksums: true,
		MetaOnly:             false,
	}

	pkg, err := deb.OpenPackageFile(pkgUri, options)

	if err != nil {
		return nil, fmt.Errorf("error while opening package file: %w", err)
	}

	res := make([]datastore.KnownFile, 0)

	for _, f := range pkg.Files() {
		md5 := pkg.GetFileMd5Sums(f.Name())
		sha1 := pkg.GetFileChecksum(f.Name())

		if !f.IsDir() && (md5 != "" || sha1 != "") {
			res = append(res, datastore.KnownFile{
				Path:      utils.StringPtr(f.Name()),
				MD5:       &md5,
				SHA1:      &sha1,
				Size:      utils.Int64Ptr(f.Size()),
				PackageID: &pkgInfo.ID,
				Status:    datastore.FileStatusSafe,

				// Package:      pkgInfo.Name,
				// Version:      pkgInfo.Version,
				// Architecture: pkgInfo.Architecture,
				// Filename:     f.Name(),
				// Mode:         f.Mode(),
				// Size:         f.Size(),
				// Owner:        f.Owner(),
				// Group:        f.Group(),
				// MD5:          md5,
				// SHA1:         sha1,
			})
		}
	}

	return res, nil
}
