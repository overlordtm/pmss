module github.com/overlordtm/pmss

go 1.20

require (
	github.com/golang/protobuf v1.5.3
	github.com/isbm/go-deb v0.0.0-20211119182924-4bb66f353d0a
	github.com/jmoiron/sqlx v1.3.5
	github.com/jszwec/csvutil v1.8.0
	github.com/mattn/go-sqlite3 v1.14.16
	github.com/sirupsen/logrus v1.9.0
	github.com/spf13/cobra v1.6.1
	google.golang.org/protobuf v1.30.0
)

require (
	github.com/andrew-d/lzma v0.0.0-20120628231508-2a7c55cad4a2 // indirect
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb // indirect
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/mitchellh/colorstring v0.0.0-20190213212951-d06e56a500db // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/schollz/progressbar/v3 v3.13.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/xi2/xz v0.0.0-20171230120015-48954b6210f8 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/term v0.6.0 // indirect
)

replace github.com/isbm/go-deb => ./vendor2/go-deb
