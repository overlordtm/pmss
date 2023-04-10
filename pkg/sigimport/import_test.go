package sigimport

import (
	"testing"
)

func TestImport(t *testing.T) {
	t.SkipNow()
	err := Import(":memory:", "testdata/malware-bazaar-recent-2023-04-01.csv")
	if err != nil {
		t.Error(err)
	}
}
