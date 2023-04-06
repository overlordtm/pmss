package debscraper

import (
	"context"
	"testing"
)

func setup() *DebScraper {
	return New(
		WithDistro("buster"),
		WithComponent("main"),
		WithArch("amd64"),
	)
}

func TestListPackages(t *testing.T) {
	s := setup()
	pkgs, err := s.ListPackages(context.Background())
	if err != nil {
		t.Error(err)
	}

	if len(pkgs) == 0 {
		t.Error("no packages found")
	}
}

func TestFetchPackage(t *testing.T) {

	packages := []packageInfo{
		{
			Name:         "libbz2-1.0",
			Version:      "1.0.6-9.2~deb10u1",
			Architecture: "amd64",
			Filename:     "pool/main/b/bzip2/libbz2-1.0_1.0.6-9.2~deb10u1_amd64.deb",
		},
	}

	s := setup()
	for _, pkg := range packages {
		files, err := s.fetchPackage(context.Background(), pkg)
		if err != nil {
			t.Error(err)
		}

		if len(files) == 0 {
			t.Error("no files found")
		}
	}
}

func TestFetchPackage2(t *testing.T) {

	packages := []packageInfo{
		{
			Name:         "libbz2-1.0",
			Version:      "1.0.6-9.2~deb10u1",
			Architecture: "amd64",
			Filename:     "pool/main/b/bzip2/libbz2-1.0_1.0.6-9.2~deb10u1_amd64.deb",
		},
	}

	s := New(
		WithRoundRobinMirrors("http://ftp.si.debian.org/debian", "http://ftp.at.debian.org/debian"),
		WithDistro("buster"),
		WithComponent("main"),
		WithArch("amd64"),
	)
	for _, pkg := range packages {
		files, err := s.retryFetchPackage(context.Background(), 5, pkg)
		if err != nil {
			t.Error(err)
		}

		if len(files) == 0 {
			t.Error("no files found")
		}
	}
}
