package cconstant

import "time"

const (
	AuthDB    string = "taskdb.public.auth"
	BannerDB  string = "taskdb.public.banner"
	VersionDB string = "taskdb.public.vbanner"
)

const (
	Salt            = "xjifcmefdx2oxe3x"
	SignedKey       = "efcj34s3dr4cwdxxjuu34"
	AccessTokenTTL  = 2 * time.Hour
	RefreshTokenTTl = 30 * 24 * time.Hour
)

const (
	RoleUser  int = 0
	RoleAdmin int = 1
)
