package debscraper

import (
	"fmt"
	"testing"
)

func setup() *DebScraper {
	return New()
}

func TestXxx(t *testing.T) {
	s := setup()
	pkgs, err := s.listPackages()
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) == 0 {
		t.Error("no packages found")
	}

	files, err := s.fetchPackage(pkgs[0])
	if err != nil {
		t.Error(err)
	}

	if len(files) == 0 {
		t.Error("no files found")
	}

	fmt.Println(files[0])
}

func TestScrape(t *testing.T) {
	s := setup()
	hashItemCh := make(chan HashItem)
	go s.Scrape(10, hashItemCh)

	for item := range hashItemCh {
		fmt.Println(item)
	}
}
