package sqlitestore

import (
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/overlordtm/pmss/pkg/datastore"
)

var (
	blacklistSchema = `
CREATE TABLE IF NOT EXISTS whitelist (
	md5 TEXT,
	sha1 TEXT,
	sha256 TEXT,
	signature TEXT,
	meta TEXT
)`
	BlacklistBatchSize = 200
)

type blacklist struct {
	store *Store

	// prepared statements
	prepareStmtOnce sync.Once
	insertStmt      *sqlx.Stmt
	bulkInsertStmt  *sqlx.Stmt
}

func (w *blacklist) ensureSchema() {
	w.store.db.MustExec(blacklistSchema)
}

func (w *blacklist) FindByMD5(md5 string) (*datastore.BlacklistItem, error) {
	row := &datastore.BlacklistItem{}
	if err := w.store.db.Get(row, "SELECT * FROM whitelist WHERE md5 = ?", md5); err != nil {
		return nil, err
	}

	return row, nil
}

func (b *blacklist) prepareInsertStmt() {
	stmt, err := b.store.db.Preparex("INSERT INTO blacklist (md5, sha1, sha256, signature, meta) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		panic(fmt.Errorf("error while preparing insert statement: %v", err))
	}
	b.insertStmt = stmt
}

func (b *blacklist) prepareBulkInsertStmt(batchSize int) {
	sqlStmt := strings.Builder{}
	sqlStmt.WriteString("INSERT INTO blacklist (md5, sha1, sha256, signature, meta) VALUES ")
	for i := 0; i < batchSize; i++ {
		sqlStmt.WriteString("(?, ?, ?, ?, ?)")
		if i < batchSize-1 {
			sqlStmt.WriteString(", ")
		}
	}
	stmt, err := b.store.db.Preparex(sqlStmt.String())
	if err != nil {
		panic(fmt.Errorf("error while preparing bulk insert statement: %v", err))
	}
	b.bulkInsertStmt = stmt
}

func (b *blacklist) Insert(row datastore.BlacklistItem) error {
	b.prepareStmtOnce.Do(func() {
		b.prepareInsertStmt()
	})

	if _, err := b.insertStmt.Exec(row.MD5, row.SHA1, row.SHA256, row.Signature, row.Meta); err != nil {
		return fmt.Errorf("error while inserting row: %v", err)
	}
	return nil
}

func (b *blacklist) InsertBatch(rows []datastore.BlacklistItem) error {
	b.prepareStmtOnce.Do(func() {
		b.prepareBulkInsertStmt(BlacklistBatchSize)
		b.prepareInsertStmt()
	})

	insertBatch := func(batch []datastore.BlacklistItem, stmt *sqlx.Stmt) {
		args := make([]interface{}, len(batch)*5)

		for i, row := range batch {
			args[i*5] = row.MD5
			args[i*5+1] = row.SHA1
			args[i*5+2] = row.SHA256
			args[i*5+3] = row.Signature
			args[i*5+4] = row.Meta
		}

		stmt.MustExec(args...)
	}

	for i := 0; i < len(rows); i += BlacklistBatchSize {
		end := i + BlacklistBatchSize
		if end > len(rows) {
			end = len(rows)
			// insert the reminder (smaller than single batch) line by line
			for j := i; j < end; j++ {
				b.insertStmt.MustExec(rows[j].MD5, rows[j].SHA1, rows[j].SHA256, rows[j].Signature, rows[j].Meta)
			}
		} else {
			insertBatch(rows[i:end], b.bulkInsertStmt)
		}
	}

	return nil
}
