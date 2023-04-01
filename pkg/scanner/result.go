package scanner

import "github.com/overlordtm/pmss/pkg/checker"

type Result struct {
	Path         string
	Err          error
	CheckResults []checker.Result
}
