package debscraper

import "os"

type HashItem struct {
	Package      string      `csv:"package"`
	Version      string      `csv:"version"`
	Architecture string      `csv:"architecture"`
	Filename     string      `csv:"filename"`
	Size         int64       `csv:"size"`
	Mode         os.FileMode `csv:"mode"`
	Owner        string      `csv:"owner"`
	Group        string      `csv:"group"`
	MD5          string      `csv:"md5"`
	SHA1         string      `csv:"sha1"`
	SHA256       string      `csv:"sha256"`
}
