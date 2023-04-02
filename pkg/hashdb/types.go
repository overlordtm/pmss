package hashdb

type WhitelistItemMeta struct {
	Package string
	Version string
	Arch    string
	Distro  string
	Size    int64
	Mode    int32
	Owner   string
	Group   string
}

type WhitelistRow struct {
	MD5    string
	SHA1   string
	SHA256 string
	Path   string
	Meta   WhitelistItemMeta
}
