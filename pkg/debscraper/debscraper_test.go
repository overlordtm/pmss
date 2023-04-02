package debscraper

import (
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
}
