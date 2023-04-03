package hashstore

type WhitelistMeta struct {
	Package string `json:"package"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
	Distro  string `json:"distro"`
	Size    int64  `json:"size"`
	Mode    uint32 `json:"mode"`
	Owner   string `json:"owner"`
	Group   string `json:"group"`
}

type WhitelistRow struct {
	MD5    string
	SHA1   string
	SHA256 string
	Path   string
	Meta   WhitelistMeta
}

type BlacklistMeta struct {
	Platform string `json:"platform"`
}

type BlacklistRow struct {
	MD5       string
	SHA1      string
	SHA256    string
	Signature string
	Meta      BlacklistMeta
}
