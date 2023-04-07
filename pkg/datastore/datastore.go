package datastore

import (
	"fmt"

	"gorm.io/gorm"
)

// Store is a datastore
type Store struct {
	opts options

	packageRepository   *packageRepository
	whitelistRepository *whitelistRepository
	blacklistRepository *blacklistRepository
	fileRepository      *fileRepository
}

type options struct {
	dialector gorm.Dialector
}

type Option func(*options)

func WithDb(dialector gorm.Dialector) Option {
	return func(o *options) {
		o.dialector = dialector
	}
}

func New(opts ...Option) (*Store, error) {

	o := options{}

	for _, option := range opts {
		option(&o)
	}

	db, err := gorm.Open(o.dialector, &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %v", err)
	}

	db.AutoMigrate(&Package{})
	db.AutoMigrate(&WhitelistItem{})
	db.AutoMigrate(&BlacklistItem{})
	db.AutoMigrate(&File{})

	return &Store{
		opts:                o,
		packageRepository:   &packageRepository{db},
		whitelistRepository: &whitelistRepository{db},
		blacklistRepository: &blacklistRepository{db},
		fileRepository:      &fileRepository{db},
	}, nil
}

func (ds *Store) Packages() PackageRepository {
	return ds.packageRepository
}

func (ds *Store) Whitelist() WhitelistRepository {
	return nil
}

func (ds *Store) Blacklist() BlacklistRepository {
	return nil
}
