package usecase

import "backend-trainee-assignment-2024/internal/banner"

type BannerUsecase struct {
	repo banner.Repository
}

func NewBannerUsecase(repo banner.Repository) banner.Usecase {
	return &BannerUsecase{repo: repo}
}

func (u *BannerUsecase) GetBanner(params *banner.GetBannerParams) (*string, error) {
	return u.repo.GetContentBanner(params)
}

func (u *BannerUsecase) GetFilteredBanners(params *banner.GetFilteredBannersParams) (*[]banner.GetFilteredBannersResponse, error) {
	if params.FeatureID != -1 && params.TagID != -1 {
		return u.repo.GetBanner(params)
	} else if params.FeatureID == -1 && params.TagID != -1 {
		return u.repo.GetFilteredBannersTID(params)
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
