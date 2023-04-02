package sigdb

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type SqlSigDb struct {
	db     *sql.DB
	dbPath string
}

func New(dbPath string) *SqlSigDb {
	return &SqlSigDb{
		db:     nil,
		dbPath: dbPath,
	}
}

func (s *SqlSigDb) Init() error {
	db, err := sql.Open("sqlite3", s.dbPath)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS signatures (
		md5 TEXT,
		sha1 TEXT,
		sha256 TEXT,
		imphash TEXT,
		ssdeep TEXT,
		tlsh TEXT,
		signature TEXT,
		filename TEXT,
		mimetype TEXT
	)
		`)
	if err != nil {
		return fmt.Errorf("error while creating table: %v", err)
	}

	s.db = db

	return nil
}

func (s *SqlSigDb) SaveItem(item Item) error {
	_, err := s.db.Exec(`INSERT INTO signatures (md5, sha1, sha256, imphash, ssdeep, tlsh, signature, filename, mimetype)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		item.MD5, item.SHA1, item.SHA256, item.ImpHash, item.SSDeep, item.TLSH, item.Signature, item.Filename, item.MimeType)
	if err != nil {
		return fmt.Errorf("error while inserting item: %v", err)
	}

	return nil
}

func (s *SqlSigDb) FindBy(field string, value string) (*Item, error) {
	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM signatures WHERE %s = ?", field), value)
	if err != nil {
		return nil, fmt.Errorf("error while querying database: %v", err)
	}

	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	item := &Item{}

	err = rows.Scan(&item.MD5, &item.SHA1, &item.SHA256, &item.ImpHash, &item.SSDeep, &item.TLSH, &item.Signature, &item.Filename, &item.MimeType)
	if err != nil {
		return nil, fmt.Errorf("error while scanning row: %v", err)
	}

	return item, nil
}

func (s *SqlSigDb) FindByMD5(md5 string) (*Item, error) {
	return s.FindBy("md5", md5)
}

func (s *SqlSigDb) FindBySHA1(sha1 string) (*Item, error) {
	return s.FindBy("sha1", sha1)
}

func (s *SqlSigDb) FindBySHA256(sha256 string) (*Item, error) {
	return s.FindBy("sha256", sha256)
}

func (s *SqlSigDb) FindByImpHash(imphash string) (*Item, error) {
	return s.FindBy("imphash", imphash)
}

func (s *SqlSigDb) FindBySSDeep(ssdeep string) (*Item, error) {
	return s.FindBy("ssdeep", ssdeep)
}

func (s *SqlSigDb) FindByTLSH(tlsh string) (*Item, error) {
	return s.FindBy("tlsh", tlsh)
}

func (s *SqlSigDb) Close() error {
	return s.db.Close()
}
