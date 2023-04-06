package sqlitestore

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/overlordtm/pmss/pkg/datastore"
)

type Store struct {
	db        *sqlx.DB
	whitelist *whitelist
	blacklist *blacklist
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

	db, err := sqlx.Open("sqlite3", o.dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %v", err)
	}

	store := &Store{
		db: db,
	}

	store.whitelist = &whitelist{
		store: store,
	}
	store.blacklist = &blacklist{
		store: store,
	}

	store.whitelist.ensureSchema()
	store.blacklist.ensureSchema()

	return store, nil
}

func (s *Store) Whitelist() datastore.WhitelistRepository {
	return s.whitelist
}

func (s *Store) Blacklist() datastore.BlacklistRepository {
	return s.blacklist
}

func (s *Store) Close() error {
	return s.db.Close()
}
