package sigimport

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Import(dbPath string, filename string) error {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// ensure table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS signatures (
		md5 TEXT,
		sha1 TEXT,
		sha256 TEXT,
		imphash TEXT,
		ssdeep TEXT,
		tlsh TEXT,
		signature TEXT,
		filename TEXT,
		mimetype TEXT)`)
	if err != nil {
		return fmt.Errorf("error while creating table: %v", err)
	}

	f, err := os.OpenFile(filename, os.O_RDONLY, 0644)

	if err != nil {
		return fmt.Errorf("error while opening file: %v", err)
	}
	defer f.Close()

	// build bulk insert sql
	sqlStmt := strings.Builder{}
	sqlStmt.WriteString("INSERT INTO signatures (md5, sha1, sha256, imphash, ssdeep, tlsh, signature, filename, mimetype) VALUES ")

	for i := 0; i < 100; i++ {
		sqlStmt.WriteString("(?, ?, ?, ?, ?, ?, ?, ?, ?)")
		if i < 99 {
			sqlStmt.WriteString(", ")
		}
	}

	bulkInsertStmt, err := db.Prepare(sqlStmt.String())
	if err != nil {
		return fmt.Errorf("error while preparing bulk insert statement: %v", err)
	}
	defer bulkInsertStmt.Close()

	singleInsertStmt, err := db.Prepare("INSERT INTO signatures (md5, sha1, sha256, imphash, ssdeep, tlsh, signature, filename, mimetype) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("error while preparing single insert statement: %v", err)
	}
	defer singleInsertStmt.Close()

	csvReader := csv.NewReader(f)
	csvReader.LazyQuotes = true
	csvReader.Comma = ','
	csvReader.Comment = '#'

	params := make([]interface{}, 0, 900)

	records, err := csvReader.ReadAll()
	if err != nil {
		return fmt.Errorf("error while reading csv: %v", err)
	}

	for i := 0; i < len(records); i = i + 100 {
		for j := 0; j < 100 && i+j < len(records); j++ {
			record := records[i+j]
			params = append(params, record[3], record[2], record[1], record[11], record[12], record[13], record[8], record[5], record[7])
		}

		if len(params) == 900 {
			if _, err := bulkInsertStmt.Exec(params...); err != nil {
				return fmt.Errorf("error while inserting bulk: %v", err)
			}
		} else {
			for k := 0; k < len(params); k = k + 9 {
				if _, err := singleInsertStmt.Exec(params[k : k+9]...); err != nil {
					return fmt.Errorf("error while inserting single item: %v", err)
				}
			}
		}

		params = params[:0]
	}

	res, err := db.Query("SELECT COUNT(*) FROM signatures")
	if err != nil {
		return fmt.Errorf("error while counting signatures: %v", err)
	}

	if !res.Next() {
		return fmt.Errorf("no results for count query")
	}

	var numRows int

	err = res.Scan(&numRows)
	if err != nil {
		return fmt.Errorf("error while scanning count: %v", err)
	}

	fmt.Printf("Loaded %d signatures\n", numRows)

	if len(records) != numRows {
		return fmt.Errorf("number of loaded signatures does not match number of records")
	}

	// item := sigdb.Item{
	// 	MD5:       record[3],
	// 	SHA1:      record[2],
	// 	SHA256:    record[1],
	// 	ImpHash:   record[11],
	// 	SSDeep:    record[12],
	// 	TLSH:      record[13],
	// 	Signature: record[8],
	// 	Filename:  record[5],
	// 	MimeType:  record[7],
	// }

	return nil
}
