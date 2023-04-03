package debscraper

import "github.com/schollz/progressbar/v3"

type ProgressDelegate interface {
	Start(int64)
	Done(int)
}

type NoopProgressBar struct{}

func (d *NoopProgressBar) Start(total int64) {}
func (d *NoopProgressBar) Done(i int)        {}

type CliProgressBar struct {
	bar *progressbar.ProgressBar
}

func (c *CliProgressBar) Start(total int64) {
	c.bar = progressbar.Default(
		total,
		"downloading",
	)
}

func (c *CliProgressBar) Done(i int) {
	c.bar.Add(i)
}
