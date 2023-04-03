package debscraper

import "github.com/sirupsen/logrus"

type options struct {
	mirrorUrl mirrorFunc
	distro    string
	component string
	arch      string
	logger    *logrus.Logger
}

type Option func(*options)

func WithMirrorUrl(mirrorUrl string) func(*options) {
	return func(o *options) {
		o.mirrorUrl = func() string {
			return mirrorUrl
		}
	}
}

func WithRoundRobinMirrors(url ...string) func(*options) {
	return func(o *options) {
		o.mirrorUrl = RoundRobinMirror(url...)
	}
}

func WithRandomMirrors(url ...string) func(*options) {
	return func(o *options) {
		o.mirrorUrl = RandomMirror(url...)
	}
}

func WithDistro(distro string) func(*options) {
	return func(o *options) {
		o.distro = distro
	}
}

func WithComponent(component string) func(*options) {
	return func(o *options) {
		o.component = component
	}
}

func WithArch(arch string) func(*options) {
	return func(o *options) {
		o.arch = arch
	}
}

func WithLogger(logger *logrus.Logger) func(*options) {
	return func(o *options) {
		o.logger = logger
	}
}
