package sqlitestore

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/overlordtm/pmss/pkg/hashstore"
)

var (
	whitelistSchema = `
CREATE TABLE IF NOT EXISTS whitelist (
	md5 TEXT,
	sha1 TEXT,
	sha256 TEXT,
	path TEXT,
	meta TEXT
)`
	WhitelistBatchSize = 200
)

type whitelist struct {
	store *Store

	// prepared statements
	prepareStmtOnce sync.Once
	insertStmt      *sqlx.Stmt
	bulkInsertStmt  *sqlx.Stmt
}

func encodeWhitelistMeta(meta *hashstore.WhitelistMeta) ([]byte, error) {
	if data, err := json.Marshal(meta); err != nil {
		return nil, fmt.Errorf("error while marshaling whitelist meta: %v", err)
	} else {
		return data, nil
	}
}

func decodeWhitelistMeta(data []byte) (*hashstore.WhitelistMeta, error) {
	meta := &hashstore.WhitelistMeta{}
	if err := json.Unmarshal(data, meta); err != nil {
		return nil, fmt.Errorf("error while unmarshaling whitelist meta: %v", err)
	}
	return meta, nil
}

func (w *whitelist) ensureSchema() {
	w.store.db.MustExec(whitelistSchema)
}

func (w *whitelist) FindByMD5(md5 string) (*hashstore.WhitelistRow, error) {
	row := &hashstore.WhitelistRow{}
	if err := w.store.db.Get(row, "SELECT * FROM whitelist WHERE md5 = ?", md5); err != nil {
		return nil, err
	}
	return row, nil
}

func (s *whitelist) prepareInsertStmt() {
	stmt, err := s.store.db.Preparex("INSERT INTO whitelist (md5, sha1, sha256, path, meta) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		panic(fmt.Errorf("error while preparing insert statement: %v", err))
	}
	s.insertStmt = stmt
}

func (s *whitelist) prepareBulkInsertStmt(batchSize int) {
	sqlStmt := strings.Builder{}
	sqlStmt.WriteString("INSERT INTO whitelist (md5, sha1, sha256, path, meta) VALUES ")
	for i := 0; i < batchSize; i++ {
		sqlStmt.WriteString("(?, ?, ?, ?, ?)")
		if i < batchSize-1 {
			sqlStmt.WriteString(", ")
		}
	}
	stmt, err := s.store.db.Preparex(sqlStmt.String())
	if err != nil {
		panic(fmt.Errorf("error while preparing bulk insert statement: %v", err))
	}
	s.bulkInsertStmt = stmt
}

func (s *whitelist) Insert(row hashstore.WhitelistRow) error {
	s.prepareStmtOnce.Do(func() {
		s.prepareInsertStmt()
	})

	metaBytes, err := encodeWhitelistMeta(&row.Meta)
	if err != nil {
		return err
	}

	_, err = s.insertStmt.Exec("INSERT INTO whitelist (md5, sha1, sha256, path, meta) VALUES (?, ?, ?, ?, ?)", row.MD5, row.SHA1, row.SHA256, row.Path, metaBytes)
	return err
}

func (s *whitelist) InsertBatch(rows []hashstore.WhitelistRow) (err error) {

	s.prepareStmtOnce.Do(func() {
		s.prepareInsertStmt()
		s.prepareBulkInsertStmt(WhitelistBatchSize)
	})

	insertBatch := func(batch []hashstore.WhitelistRow, stmt *sqlx.Stmt) error {
		args := make([]interface{}, len(batch)*5)

		for i, row := range batch {

			metaBytes, err := encodeWhitelistMeta(&row.Meta)
			if err != nil {
				return err
			}

			args[i*5] = row.MD5
			args[i*5+1] = row.SHA1
			args[i*5+2] = row.SHA256
			args[i*5+3] = row.Path
			args[i*5+4] = metaBytes
		}

		_, err := stmt.Exec(args...)
		return err
	}

	for i := 0; i < len(rows); i += WhitelistBatchSize {
		end := i + WhitelistBatchSize
		if end > len(rows) {
			end = len(rows)
			// insert the reminder (smaller than single batch) line by line
			for j := i; j < end; j++ {
				metaBytes, err := encodeWhitelistMeta(&rows[j].Meta)
				if err != nil {
					return err
				}
				s.insertStmt.MustExec(rows[j].MD5, rows[j].SHA1, rows[j].SHA256, rows[j].Path, metaBytes)
			}
		} else {
			if err := insertBatch(rows[i:end], s.bulkInsertStmt); err != nil {
				return err
			}
		}
	}

	return nil
}
