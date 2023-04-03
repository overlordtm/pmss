package debscraper

import (
	"math/rand"
	"sync/atomic"
)

type mirrorFunc func() string

var Mirrors = []string{
	"http://ftp.at.debian.org/debian",
	"http://ftp.de.debian.org/debian",
	"http://ftp.si.debian.org/debian",
	"http://ftp.hr.debian.org/debian",
	"http://ftp.ch.debian.org/debian",
	"http://ftp.cz.debian.org/debian",
	"http://ftp.hu.debian.org/debian",
	"http://ftp.it.debian.org/debian",
	"http://ftp.nl.debian.org/debian",
	"http://ftp.pl.debian.org/debian",
	"http://ftp.ro.debian.org/debian",
	"http://ftp.sk.debian.org/debian",
	"http://ftp.es.debian.org/debian",
	"http://ftp.pt.debian.org/debian",
	"http://ftp.fr.debian.org/debian",
	"http://ftp.uk.debian.org/debian",
	"http://ftp.us.debian.org/debian",
}

func RandomMirror(url ...string) func() string {
	return func() string {
		return url[rand.Intn(len(url))]
	}
}

func RoundRobinMirror(url ...string) func() string {
	var idx uint32
	return func() string {
		i := atomic.AddUint32(&idx, 1)
		return url[i%uint32(len(url))]
	}
}
