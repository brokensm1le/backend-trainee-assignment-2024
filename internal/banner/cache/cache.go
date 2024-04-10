package cache

import "backend-trainee-assignment-2024/internal/banner"

type Cache interface {
	LoadCache() error
	GetBanner(fid int64, tid int64) (banner.GetFilteredBannersResponse, error)
	GetBannersByTID(tid int64) ([]banner.GetFilteredBannersResponse, error)
	GetBannersByFID(fid int64) ([]banner.GetFilteredBannersResponse, error)
}
