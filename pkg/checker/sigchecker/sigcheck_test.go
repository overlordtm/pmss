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
	dbPath := fmt.Sprintf("%s/test/test.db", projRoot)
	db := sigdb.New(dbPath)

	defer db.Close()

	err := db.Init()
	if err != nil {
		t.Fatal(err)
	}

	checker := sigchecker.New(db)

	exePath := fmt.Sprintf("%s/test/data/179b98e2cb16a094755f853ae892b47948a8b6a83e7ca050d520e113ff180b2f.exe", projRoot)
	f, err := os.OpenFile(exePath, os.O_RDONLY, 0)
	if err != nil {
		t.Fatal(err)
	}

	defer f.Close()

	h, err := multihasher.Hash(f)
	if err != nil {
		t.Fatal(err)
	}

	r, err := checker.CheckHash(h)
	if err != nil {
		t.Fatal(err)
	}

	if r.Signature != "Formbook" {
		t.Fatal("wrong signature")
	}
}
