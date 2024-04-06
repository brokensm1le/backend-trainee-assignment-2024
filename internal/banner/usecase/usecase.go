package usecase

import "backend-trainee-assignment-2024/internal/banner"

type BannerUsecase struct {
	repo banner.Repository
}

func NewBannerUsecase(repo banner.Repository) banner.Usecase {
	return &BannerUsecase{repo: repo}
}
