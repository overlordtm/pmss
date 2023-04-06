package datastore

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Store is a datastore
type Store struct {
	opts options

	packageRepository *packageRepository
}

type options struct {
	dbUrl string
}

type Option func(*options)

func WithDbUrl(dbUrl string) Option {
	return func(o *options) {
		o.dbUrl = dbUrl
	}
}

func New(opts ...Option) (*Store, error) {

	o := options{
		dbUrl: ":memory:",
	}

	for _, option := range opts {
		option(&o)
	}

	db, err := gorm.Open(sqlite.Open(o.dbUrl), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %v", err)
	}

	db.AutoMigrate(&Package{})
	db.AutoMigrate(&WhitelistItem{})
	db.AutoMigrate(&BlacklistItem{})

	return &Store{
		opts:              o,
		packageRepository: &packageRepository{db},
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
