package pkgscraper

import (
	"context"
	"fmt"

	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/overlordtm/pmss/pkg/debscraper"
	"gorm.io/gorm"
)

func ScrapeDebianMirror(ctx context.Context, db *gorm.DB, distro, arch, component string) error {
	scraper := debscraper.New(debscraper.WithDistro(distro), debscraper.WithArch(arch), debscraper.WithComponent(component), debscraper.WithOsType(datastore.OsTypeDebian))
	if packages, err := scraper.ListPackages(ctx); err != nil {
		return err
	} else {

		for _, pkg := range packages {
			if err := datastore.Packages().Save(pkg)(db); err != nil {
				continue
			}
		}

	}
	return nil
}

func ScrapeDebianPackage(ctx context.Context, tx *gorm.DB, pkg datastore.Package) error {
	scraper := debscraper.New(debscraper.WithDistro(pkg.Distro), debscraper.WithArch(pkg.Architecture), debscraper.WithComponent(pkg.Component))

	knownFiles, err := scraper.FetchPackage(ctx, 3, pkg)
	if err != nil {
		return fmt.Errorf("error while fetching package: %w", err)
	}

	for _, knownFile := range knownFiles {
		if err := datastore.KnownFiles().Create(knownFile)(tx); err != nil {
			return fmt.Errorf("error while creating known file: %w", err)
		}
	}
	return nil
}

func ScrapeUbuntu(distro, arch, component string) ([]datastore.Package, error) {
	return nil, nil
}
