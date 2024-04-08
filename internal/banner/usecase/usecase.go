package usecase

import "backend-trainee-assignment-2024/internal/banner"

type BannerUsecase struct {
	repo banner.Repository
}

func NewBannerUsecase(repo banner.Repository) banner.Usecase {
	return &BannerUsecase{repo: repo}
}

func (u *BannerUsecase) GetBanner(params *banner.GetBannerParams) (*string, error) {
	return u.repo.GetBanner(params)
}

func (u *BannerUsecase) CreateBanner(params *banner.CreateBannerParams) (int64, error) {
	return u.repo.CreateBanner(params)
}
