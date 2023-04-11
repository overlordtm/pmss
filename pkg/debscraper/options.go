package debscraper

import (
	"github.com/overlordtm/pmss/pkg/datastore"
	"github.com/sirupsen/logrus"
)

type Option func(*DebScraper)

func WithMirrorUrl(mirrorUrl string) Option {
	return func(o *DebScraper) {
		o.mirrorUrl = func() string {
			return mirrorUrl
		}
	}
}

func WithRoundRobinMirrors(url ...string) Option {
	return func(o *DebScraper) {
		o.mirrorUrl = RoundRobinMirror(url...)
	}
}

func WithRandomMirrors(url ...string) Option {
	return func(o *DebScraper) {
		o.mirrorUrl = RandomMirror(url...)
	}
}

func WithDistro(distro string) Option {
	return func(o *DebScraper) {
		o.distro = distro
	}
}

func WithComponent(component string) Option {
	return func(o *DebScraper) {
		o.component = component
	}
}

func WithArch(arch string) Option {
	return func(o *DebScraper) {
		o.arch = arch
	}
}

func WithLogger(logger *logrus.Logger) Option {
	return func(o *DebScraper) {
		o.logger = logger
	}
}

func WithOsType(osType datastore.OsType) Option {
	return func(o *DebScraper) {
		o.osType = osType
	}
}
