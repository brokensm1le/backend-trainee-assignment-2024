package usecase

import (
	"backend-trainee-assignment-2024/internal/banner"
	"backend-trainee-assignment-2024/internal/banner/cache"
	"backend-trainee-assignment-2024/internal/cconstant"
)

type BannerUsecase struct {
	repo  banner.Repository
	cache cache.Cache
}

func NewBannerUsecase(repo banner.Repository, cache cache.Cache) banner.Usecase {
	return &BannerUsecase{repo: repo, cache: cache}
}

func (u *BannerUsecase) GetBanner(params *banner.GetBannerParams) (*string, error) {
	if !params.UseLastRevision {
		getBanner, err := u.cache.GetBanner(params.FeatureID, params.TagID)
		if err == nil {
			return &getBanner.Content, nil
		}
	}

	if params.Role == cconstant.RoleAdmin {
		return u.repo.GetContentBannerAdmin(params)
	}
	return u.repo.GetContentBanner(params)
}

func (u *BannerUsecase) GetFilteredBanners(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	if params.FeatureID != -1 && params.TagID != -1 {
		if !params.UseLastRevision {
			getBanner, err := u.cache.GetBanner(params.FeatureID, params.TagID)
			if err == nil {
				return &[]banner.GetFilteredBannersResponse{getBanner}, nil
			}
		}

		if params.Role == cconstant.RoleAdmin {
			return u.repo.GetBannerAdmin(params)
		}
		return u.repo.GetBanner(params)
	} else if params.FeatureID == -1 && params.TagID != -1 {
		if !params.UseLastRevision {
			getBanner, err := u.cache.GetBannersByTID(params.TagID)
			if err == nil {
				return &getBanner, nil
			}
		}

		if params.Role == cconstant.RoleAdmin {
			return u.repo.GetFilteredBannersTIDAdmin(params)
		}
		return u.repo.GetFilteredBannersTID(params)
	}

	if !params.UseLastRevision {
		getBanner, err := u.cache.GetBannersByFID(params.FeatureID)
		if err == nil {
			return &getBanner, nil
		}
	}
	if params.Role == cconstant.RoleAdmin {
		return u.repo.GetFilteredBannersFIDAdmin(params)
	}
	return u.repo.GetFilteredBannersFID(params)
}

func (u *BannerUsecase) CreateBanner(params *banner.CreateBannerParams) (int64, error) {
	return u.repo.CreateBanner(params)
}

func (u *BannerUsecase) DeleteBanner(id int64) error {
	return u.repo.DeleteBanner(id)
}

func (u *BannerUsecase) UpdateUser(params *banner.UpdateBannerParams) error {
	return u.repo.UpdateUser(params)
}
