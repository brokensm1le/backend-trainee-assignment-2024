package http

import (
	"backend-trainee-assignment-2024/internal/banner"
	"github.com/gofiber/fiber/v2"
)

type BannerHandler struct {
	bannerUC banner.Usecase
}

func NewBannerHandler(bannerUC banner.Usecase) *BannerHandler {
	return &BannerHandler{
		bannerUC: bannerUC,
	}
}

func (h *BannerHandler) GetBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) GetFilteredBanners() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) CreateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) UpdateBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func (h *BannerHandler) DeleteBanner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
