package hashdb

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

const (
	bulkInsertCount = 300
)

var (
	schema = `
CREATE TABLE IF NOT EXISTS whitelist (
	md5 TEXT,
	sha1 TEXT,
	sha256 TEXT,
	path TEXT,
	meta TEXT
)

CREATE TABLE IF NOT EXISTS whitelist (
	md5 TEXT,
	sha1 TEXT,
	sha256 TEXT,
	signature TEXT,
	meta TEXT
)`
)

type SqlHashDb struct {
	db                 *sqlx.DB
	insertStmt         *sqlx.Stmt
	bulkInsertStmt     *sqlx.Stmt
	bulkInsertStmtOnce sync.Once
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

func New(opts ...Option) (*SqlHashDb, error) {

	o := options{
		dbUrl: ":memory:",
	}

	db, err := sqlx.Open("sqlite3", o.dbUrl)
	if err != nil {
		return nil, fmt.Errorf("error while opening database: %v", err)
	}

	db.MustExec(schema)

	return &SqlHashDb{
		db: db,
	}, nil
}

func (s *SqlHashDb) Close() error {
	return s.db.Close()
}

func (s *SqlHashDb) FindByMD5(md5 string) (*WhitelistRow, error) {
	row := &WhitelistRow{}
	if err := s.db.Get(row, "SELECT * FROM whitelist WHERE md5 = ?", md5); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *SqlHashDb) BulkInsert(rows []WhitelistRow) (err error) {

	s.bulkInsertStmtOnce.Do(func() {
		stmt, err1 := s.db.Preparex("INSERT INTO whitelist (md5, sha1, sha256, path, meta) VALUES (?, ?, ?)")

		if err1 != nil {
			err = fmt.Errorf("error while preparing insert statement: %v", err1)
			return
		}
		s.insertStmt = stmt

		sqlStmt := strings.Builder{}
		sqlStmt.WriteString("INSERT INTO whitelist (md5, sha1, sha256, path, meta) VALUES ")
		for i := 0; i < bulkInsertCount; i++ {
			sqlStmt.WriteString("(?, ?, ?, ?, ?)")
			if i < bulkInsertCount-1 {
				sqlStmt.WriteString(", ")
			}
		}
		stmt, err1 = s.db.Preparex(sqlStmt.String())
		if err1 != nil {
			err = fmt.Errorf("error while preparing bulk insert statement: %v", err1)
			return
		}
		s.bulkInsertStmt = stmt
	})

	if err != nil {
		return err
	}

	// tx := s.db.MustBegin()

	insertBatch := func(batch []WhitelistRow, stmt *sqlx.Stmt) {
		args := make([]interface{}, len(batch)*5)

		for i, row := range batch {
			args[i*5] = row.MD5
			args[i*5+1] = row.SHA1
			args[i*5+2] = row.SHA256
			args[i*5+3] = row.Path
			args[i*5+4] = row.Meta
		}

		stmt.MustExec(args...)
	}

	for i := 0; i < len(rows); i += bulkInsertCount {
		end := i + bulkInsertCount
		if end > len(rows) {
			end = len(rows)

			for j := i; j < end; j++ {
				s.insertStmt.MustExec(rows[j].MD5, rows[j].SHA1, rows[j].SHA256, rows[j].Path, rows[j].Meta)
			}
		}

		insertBatch(rows[i:end], s.bulkInsertStmt)
	}

	return nil
}
