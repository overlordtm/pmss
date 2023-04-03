package debscraper

import (
	"fmt"
	"net/http"
	"testing"
)

func TestMirrorsList(t *testing.T) {
	for _, mirror := range Mirrors {
		res, err := http.Get(mirror)
		if err != nil {
			t.Error(fmt.Errorf("error while fetching mirror %s: %w", mirror, err))
		}
		if res.StatusCode != 200 {
			t.Error(fmt.Errorf("mirror %s returned status code %d", mirror, res.StatusCode))
		}
	}
}
