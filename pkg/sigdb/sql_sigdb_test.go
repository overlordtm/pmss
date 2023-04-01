package sigdb_test

import (
	"testing"

	"github.com/overlordtm/pmss/pkg/sigdb"
)

func TestInit(t *testing.T) {
	s := sigdb.New(":memory:")
	err := s.Init()
	if err != nil {
		t.Error(err)
	}

	item := sigdb.Item{
		MD5:       "abc",
		SHA1:      "",
		SHA256:    "",
		ImpHash:   "",
		SSDeep:    "",
		TLSH:      "",
		Signature: "Krneki",
		Filename:  "file.txt",
		MimeType:  "application/x-dosexec",
	}

	err = s.SaveItem(item)

	if err != nil {
		t.Error(err)
	}

	foundItem, err := s.FindByMD5("abc")
	if err != nil {
		t.Error(err)
	}

	if foundItem == nil {
		t.Error("item not found")
	}
}
