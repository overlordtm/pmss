package sigchecker_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/overlordtm/pmss/pkg/checker/sigchecker"
	"github.com/overlordtm/pmss/pkg/multihasher"
	"github.com/overlordtm/pmss/pkg/sigdb"
)

func TestSigChecker_CheckHash(t *testing.T) {
	projRoot := os.Getenv("PMSS_PROJ_ROOT")
	dbPath := fmt.Sprintf("%s/test/0ad.db", projRoot)
	db := sigdb.New(dbPath)

	err := db.Init()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	checker := sigchecker.New(db)

	testFile := "0ad.appdata.xml"
	testFilePath := fmt.Sprintf("%s/test/data/%s", projRoot, testFile)
	f, err := os.OpenFile(testFilePath, os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	h, err := multihasher.Hash(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("hash: %s\n", h.MD5)
	r, err := checker.CheckHash(h)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("r: %+v\n", r)

	if r.Signature != "Formbook" {
		t.Fatal("wrong signature")
	}
}
